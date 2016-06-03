package q

/*
1. 根据配置建立内存队列组， 为每个内存队列分配一个 chan 接收网络请求
2. 请求到了后， 根据target-topic client-id 找到
*/

import (
	"errors"
	"log"
	"sync"
	"time"
)

var QMap = NewQMap() //全局的一个保存所有队列的对象

const REQ_TIMEOUT = 20 //请求应答的超时时间， 单位秒

const (
	Q_COMMIT   = 100
	Q_ROLLBACK = 200
	Q_PULL     = 300
	Q_FAIL     = 400
)

type Q1 struct {
	qtype int    //请求类型 commit: 100; rollback: 200; pull: 300;
	msgid string //消息id
	tm    int    //时间， 对于pull为process_time

	m   *MsgIndex //应答的数据， 成功则返回对应消息
	ret bool      //应答是否成功
}

// 一个请求对象
type Q struct {
	data interface{}      //数据
	ch   chan interface{} //应答 chan
}

// 管理socket上并发请求到内存队列的操作
type qchan struct {
	req chan *Q //请求chan
	ans chan *Q // 应答 chan, 需要预先存于一组 chan interface{}
	q   *PrioQueue
}

func newQC(q *PrioQueue) *qchan {
	MAX_REQ := 30
	ans := make(chan *Q, MAX_REQ)
	for i := 0; i < MAX_REQ; i++ {
		ans <- &Q{ch: make(chan interface{})}
	}
	return &qchan{q: q,
		req: make(chan *Q), //请求chan中不设缓存
		ans: ans}           //允许30个并发请求等待
}

/* 通过单 goroutine 来操作一个 memq, 确保数据一致性
根据请求的数据类型， 来区分请求类型
	SEND_COMMIT: *MsgIndex,  返回 bool
	PULL:
		"":           只pull一个消息， 返回 *MsgIndex
		PULL_COMMIT:
		PULL_COMMIT_AND_NEXT:
		PULL_ROLLBACK:
		PULL_ROLLBACK_AND_NEXT:
		PULL_ROLLBACK_BUT_RETRY:
	UNLOCK:
	DELETE:
*/
func (qc *qchan) deal_request() {
	for r := range qc.req {
		if v, ok := r.data.(*MsgIndex); ok {
			r.ch <- qc.q.Push(v)
		} else if v, ok := r.data.(*Q1); ok {
			switch v.qtype {
			case Q_PULL:
				log.Printf("before pull: size: %d\n", qc.q.Size())
				x := qc.q.Pull(int64(v.tm))
				log.Printf("after pull: size=%d, x=[%v]", qc.q.Size(), x)
				r.ch <- x
			case Q_COMMIT:
				x := qc.q.Commit(v.msgid, "commit")
				log.Printf("commit: %s, return: %v, size=%d", v.msgid, x, qc.q.Size())
				r.ch <- x
			case Q_ROLLBACK:
				r.ch <- qc.q.Rollback(v.msgid, "rollback")
			case Q_FAIL:
				r.ch <- qc.q.Failure(v.msgid, "fail")
			}
		}
	}
}
func (qc *qchan) request(r interface{}) interface{} {
	// 先从 ans中取得一个Q， 然后把数据放进去， 然后在 Q.ch 上等待结果
	select {
	case q := <-qc.ans:
		defer func() { qc.ans <- q }() //Q 需要放回到chan中
		q.data = r
		select { //发送请求
		case qc.req <- q:
			select { //等待应答
			case s := <-q.ch:
				return s
			case <-time.After(time.Second * REQ_TIMEOUT):
				log.Printf("wait answer timeout\n")
			}
		case <-time.After(time.Second * REQ_TIMEOUT):
			log.Printf("send request timeout\n")
		}
	case <-time.After(time.Second * REQ_TIMEOUT):
		log.Printf("get Q timeout\n")

	}
	return nil
}
func (qc *qchan) close() {
	close(qc.req)
}

type qmap struct {
	sync.RWMutex
	m map[string]map[string]*qchan
}

func NewQMap() *qmap {
	return &qmap{m: make(map[string]map[string]*qchan)}
}

func (m *qmap) Put(q *PrioQueue) error {
	m.Lock()
	defer m.Unlock()
	var topic, client string
	topic = q.Dst_topic_id
	client = q.Dst_cli_id

	if x1, ok := m.m[topic]; ok {
		if _, ok := x1[client]; ok {
			return errors.New("already exists")
		}
		q := newQC(q)
		x1[client] = q
		go q.deal_request() //启动队列处理goroutine
	} else {
		log.Printf("===not found topic[%s], add new\n", topic)
		x1 := make(map[string]*qchan)
		q := newQC(q)
		x1[client] = q
		m.m[topic] = x1
		//		log.Printf("==== go deal_request for [%s] [%s]\n", topic, client)
		go q.deal_request() //启动队列处理goroutine
	}

	return nil
}

func (sm *qmap) Get(topic, client string) *qchan {
	sm.RLock()
	defer sm.RUnlock()

	if c, ok := sm.m[topic]; ok {
		if x, ok := c[client]; ok {
			return x
		}
	}
	return nil
}
func (m *qmap) getClients(topic string) []string {
	m.RLock()
	defer m.RUnlock()

	if c, ok := m.m[topic]; ok {
		keys := make([]string, 0, len(c))
		for k := range c {
			keys = append(keys, k)
		}
		return keys
	}
	return nil
}

func (m *qmap) Close() {
	for _, c := range m.m {
		for _, q := range c {
			q.close()
			q.q.Clear()
		}
	}
}

/*生产者二阶段提交 send-commit*/
func (m *qmap) REQ_SendCommit(msg *Message, msgr *Message, broker_id string) error {
	topic := msg.GetPropStr(TARGET_TOPIC)
	if topic == "" {
		log.Printf("missing parameter TARGET_TOPIC[%s]\n", topic)
		msgr.SetProp(RESULT_CODE, REQUIRED_PARAMETER_MISSING)
		return nil
	}
	clients := m.getClients(topic)
	if clients == nil {
		log.Printf("can not found clients for topic[%s]\n", topic)
		msgr.SetProp(RESULT_CODE, INTERNAL_SERVICE_UNAVAILABLE)
		return nil
	}

	groupid := msg.GetPropStr(GROUP)
	priority := msg.GetPropInt(PRIORITY)
	produc_cli := msg.GetPropStr(CLIENT_ID) // producer client
	src_topic := msg.GetPropStr(TOPIC)
	effective_time := msg.GetPropInt(EFFECTIVE_TIME)
	expire_time := msg.GetPropInt(EXPIRE_TIME)
	priority_name := msg.GetPropStr(PRIORITYNAME)

	msgids := make([]string, 0)
	if mid := msg.GetPropStr(MESSAGE_ID); mid != "" {
		msgids = append(msgids, mid)
	} else if v, ok := msg.Props[BATCH_MESSAGE_ID]; ok {
		log.Printf("send-commit: batch-message-id")
		if v1, ok := v.([]string); ok {
			for i := range v1 {
				msgids = append(msgids, v1[i])
			}
		}
	}
	if len(msgids) == 0 {
		msgr.SetProp(RESULT_CODE, BAD_REQUEST)
		return nil
	}
	all_succ := true
	for _, cli_id := range clients {
		q := m.Get(topic, cli_id)
		if q == nil {
			log.Printf("can not found qchan %s %s\n", topic, cli_id)
			continue
		}
		for _, mid := range msgids {
			v := NewMsgIndex(mid)
			v.GroupId = groupid
			if priority > 0 {
				v.Priority = priority
			}
			if expire_time > 0 {
				v.Expire = int64(expire_time)
			}
			if effective_time > 0 {
				v.GetTime = int64(effective_time)
			}
			if priority_name != "" {
				v.Props["priority_name"] = priority_name
			}
			if produc_cli != "" {
				v.Props["produce_cli_id"] = produc_cli
			}
			if src_topic != "" {
				v.Props["produce_cli_id"] = src_topic
			}
			if broker_id != "" {
				v.Props["broker_id"] = broker_id
			}

			r := q.request(v) // 通过chan发送请求, 并等待应答
			if r == nil {
				log.Printf("q.request can not get answer")
				all_succ = false //这个消息发送失败
				break
			} else if r1, ok := r.(bool); !ok || !r1 {
				log.Printf("send-commit: return false[%v] or none-bool[%v]", r1, ok)
				all_succ = false
				break
			}
			log.Printf("send-commit succ %s", mid)
		}
	}

	if all_succ {
		msgr.SetProp(RESULT_CODE, OK)
	} else {
		msgr.SetProp(RESULT_CODE, INTERNAL_DATA_ACCESS_EXCEPTION)
	}

	return nil
}
func (m *qmap) REQ_Pull(msg *Message, msgr *Message, broker_id string) error {
	topic := msg.GetPropStr(TARGET_TOPIC)
	client := msg.GetPropStr(CLIENT_ID)
	if topic == "" || client == "" {
		log.Printf("REQUIRED_PARAMETER_MISSING, TARGET_TOPIC=[%s],CLIENT_ID=[%s]",
			topic, client)
		msgr.SetProp(RESULT_CODE, REQUIRED_PARAMETER_MISSING)
		return nil
	}
	q := m.Get(topic, client)
	if q == nil {
		log.Printf("pull: not found topic[%s] client[%s]", topic, client)
		msgr.SetProp(RESULT_CODE, INTERNAL_SERVICE_UNAVAILABLE)
		return nil
	}
	getnext := true

	if pc := msg.GetPropStr(PULL_CODE); pc != "" {

		msgid := msg.GetPropStr(MESSAGE_ID)
		switch pc {
		case PULL_COMMIT:
			getnext = false
			log.Printf("pull-commit: %s", msgid)

			r := q.request(&Q1{qtype: Q_COMMIT, msgid: msgid})
			if b, ok := r.(bool); ok {
				if !b { //commit failed
					log.Printf("commit msgid=%s failed\n", msgid)
					msgr.SetProp(RESULT_CODE, INTERNAL_SERVER_ERROR)
					return nil
				}
			} else {
				log.Printf("invalid commit return type, want bool")
				msgr.SetProp(RESULT_CODE, INTERNAL_SERVER_ERROR)
				return nil
			}
		case PULL_COMMIT_AND_NEXT:
			log.Printf("pull-commit-&-next: %s", msgid)
			r := q.request(&Q1{qtype: Q_COMMIT, msgid: msgid})
			if b, ok := r.(bool); ok {
				if !b { //commit failed
					log.Printf("commit msgid=%s failed\n", msgid)
					msgr.SetProp(RESULT_CODE, INTERNAL_SERVER_ERROR)
					return nil
				}
			} else {
				log.Printf("invalid commit return type, want bool")
				msgr.SetProp(RESULT_CODE, INTERNAL_SERVER_ERROR)
				return nil
			}
		case PULL_ROLLBACK:
			getnext = false
			log.Printf("pull-rollback: %s", msgid)
			r := q.request(&Q1{qtype: Q_ROLLBACK, msgid: msgid})
			if b, ok := r.(bool); ok {
				if !b { //commit failed
					log.Printf("rollback msgid=%s failed\n", msgid)
					msgr.SetProp(RESULT_CODE, INTERNAL_SERVER_ERROR)
					return nil
				}
			} else {
				log.Printf("invalid rollback return type, want bool")
				msgr.SetProp(RESULT_CODE, INTERNAL_SERVER_ERROR)
				return nil
			}
		case PULL_ROLLBACK_AND_NEXT:
			log.Printf("pull-rollback-&-next: %s", msgid)
			r := q.request(&Q1{qtype: Q_FAIL, msgid: msgid})
			if b, ok := r.(bool); ok {
				if !b { //commit failed
					log.Printf("fail msgid=%s failed\n", msgid)
					msgr.SetProp(RESULT_CODE, INTERNAL_SERVER_ERROR)
					return nil
				}
			} else {
				log.Printf("invalid fail return type, want bool")
				msgr.SetProp(RESULT_CODE, INTERNAL_SERVER_ERROR)
				return nil
			}
		}
	}

	if getnext { //获取下一条消息
		ptime := msg.GetPropInt(PROCESSING_TIME)
		r := &Q1{qtype: Q_PULL, tm: ptime}
		r1 := q.request(r)
		log.Printf("pull return [%v], is nill=%t", r1, r1 == nil)
		if v, ok := r1.(*MsgIndex); ok {
			if v == nil {
				log.Printf("pull return nil\n")
				msgr.SetProp(RESULT_CODE, NO_MORE_MESSAGE)
				return nil
			}
			//成功获得消息
			msgr.SetProp(RESULT_CODE, OK)
			msgr.SetProp(MESSAGE_ID, v.Id)
			msgr.SetProp(CONSUMER_RETRY, v.Retry)
			log.Printf("pull: got one: %s retry: %d", v.Id, v.Retry)
		} else {
			log.Printf("pull return invalid type\n")
			msgr.SetProp(RESULT_CODE, INTERNAL_SERVER_ERROR)
		}

	} else {
		msgr.SetProp(RESULT_CODE, OK)
	}
	return nil
}

func (m *qmap) REQ_Delete(msg *Message, msgr *Message, broker_id string) error {
	topic := msg.GetPropStr(TARGET_TOPIC)
	client := msg.GetPropStr(CLIENT_ID)
	q := m.Get(topic, client)
	if q == nil {
		msgr.SetProp(RESULT_CODE, INTERNAL_SERVICE_UNAVAILABLE)
		return nil
	}
	return nil
}
func (m *qmap) REQ_Unlock(msg *Message, msgr *Message, broker_id string) error {
	topic := msg.GetPropStr(TARGET_TOPIC)
	client := msg.GetPropStr(CLIENT_ID)
	q := m.Get(topic, client)
	if q == nil {
		msgr.SetProp(RESULT_CODE, INTERNAL_SERVICE_UNAVAILABLE)
		return nil
	}
	return nil
}
