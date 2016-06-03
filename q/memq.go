package q

import (
	"fmt"
	"log"
	"strings"
	"time"
)

type MsgIndex struct {
	Id       string
	GroupId  string
	Priority int
	GetTime  int64 //锁定时间
	Expire   int64 //有效时间
	Props    map[string]string
	Retry    int //消费重试次数

	Data string

	mark_del bool
}

func NewMsgIndex(msgid string) *MsgIndex {
	return &MsgIndex{Id: msgid, mark_del: false, Props: make(map[string]string)}
}

type PrioQueue struct {
	Dst_cli_id     string
	Dst_topic_id   string
	Next_client_id string

	totalCount int64

	min_locking_time int64 //当前队列中在途消息最小的锁定时间， -1 为无锁定

	max_onroad      int                  //最大在途消息
	max_tries       int                  //最大消费重试次数
	grouplocks      map[string]int64     //已锁定的组
	messages        map[string]*List     //以groupid为key 保存消息链表
	messageids      map[string]*MsgIndex //以messageid 为key， 保存的消息索引
	lockingMessages map[string]*MsgIndex //所有在途消息
	plist           *SortList            //优先级链表， 其值为groupid
}

func (p *PrioQueue) setminlocktime(tm_now, tm_lock int64) {
	if p.min_locking_time < tm_now {
		p.min_locking_time = tm_lock
	} else if tm_lock < p.min_locking_time {
		p.min_locking_time = tm_lock
	}
}

func NewQueue(topic_id, cli_id string) *PrioQueue {
	r := &PrioQueue{Dst_cli_id: cli_id, Dst_topic_id: topic_id,
		max_onroad:       20,
		max_tries:        10,
		min_locking_time: -1,
		grouplocks:       make(map[string]int64),
		messages:         make(map[string]*List),
		messageids:       make(map[string]*MsgIndex),
		lockingMessages:  make(map[string]*MsgIndex),
		plist:            NewSortList()}
	return r
}

func (p *PrioQueue) Push(mi *MsgIndex) bool {
	if _, ok := p.messageids[mi.Id]; ok {
		log.Printf("push: exists %s\n", mi.Id)
		return false // already exists
	}
	if mi.GroupId == "" {
		mi.GroupId = fmt.Sprintf("[groupid]%d", mi.Priority)
	} else if mi.GetTime > 0 {
		p.grouplocks[mi.Id] = mi.GetTime
	}

	t1 := time.Now().Unix()
	if mi.GetTime > t1 {
		p.lockingMessages[mi.Id] = mi
		p.setminlocktime(t1, mi.GetTime)
	}
	p.plist.Insert(mi.Priority, mi.GroupId) //groupid 添加到优先级链表中

	lm, ok := p.messages[mi.GroupId]
	if !ok {
		lm = newList()
		p.messages[mi.GroupId] = lm
	}
	lm.pushBack(mi)
	p.messageids[mi.Id] = mi

	return true
}

func (p *PrioQueue) remove(mi *MsgIndex) {
	//删除一个消息
	delete(p.grouplocks, mi.GroupId)
	delete(p.lockingMessages, mi.Id)
	delete(p.messageids, mi.Id)
}

func (p *PrioQueue) Pull(locktime int64) *MsgIndex {
	var mi *MsgIndex
	mi = nil
	time_now := time.Now().Unix()
	if len(p.lockingMessages) >= p.max_onroad && p.min_locking_time > time_now {
		log.Printf("pull: max onroad reached, %d\n", p.max_onroad)
		return nil
	}
	p.plist.Iter(func(key int, groupid string) int {
		//		log.Printf("check group: %s\n", groupid)
		// 检查组锁定
		if strings.Index(groupid, "[groupid]") != 0 {
			tm, ok := p.grouplocks[groupid]
			if ok {
				if tm <= time_now {
					// 组锁定时间过期, 解锁
					log.Printf("pull: remove expired locking group %s\n", groupid)
					delete(p.grouplocks, groupid)
				} else {
					log.Printf("pull: group locked: %s\n", groupid)
					return 1 // continue next group
				}
			}
		}
		lm, ok := p.messages[groupid]
		if !ok || lm.Size() == 0 {
			//			log.Printf("pull: remove an empty group: %s\n", groupid)
			return -1
		}
		for x := lm.head; x != nil; {
			mx := x.value.(*MsgIndex)
			//			log.Printf("pull: check id %s\n", mx.Id)
			if mx.mark_del {
				// 标记删除的， 删除
				// p.remove(mx), 这个是在 commit 的时候删除的
				//				log.Printf("pull: msg removed %s\n", mx.Id)
				x = lm.remove(x)
				continue
			}
			if mx.Expire > 0 && mx.Expire <= time_now {
				// 超过有效期， 删除
				p.remove(mx)
				x = lm.remove(x)
				log.Printf("pull: msg expired %s\n", mx.Id)
				continue
			}
			if mx.GetTime > 0 {
				if mx.GetTime > time_now {
					// 消息锁定时间未超时
					x = x.next
					log.Printf("pull: msg locked %s\n", mx.Id)
					continue
				}
			}
			if mx.Retry >= p.max_tries {
				// 超过重试次数， 删除
				p.remove(mx)
				x = lm.remove(x)
				log.Printf("pull: msg too retries %s, %d\n", mx.Id, p.max_tries)
				continue
			}
			mi = mx
			return 0 //found
		}
		return 1
	})
	if mi != nil {
		mi.GetTime = time_now + locktime
		mi.Retry++
		p.lockingMessages[mi.Id] = mi
		if strings.Index(mi.GroupId, "[groupid]") != 0 {
			p.grouplocks[mi.GroupId] = mi.GetTime
		}
		p.setminlocktime(time_now, mi.GetTime)
	}
	return mi
}

func (p *PrioQueue) Commit(mid, note string) bool {
	mi, ok := p.lockingMessages[mid]
	if !ok {
		return false
	}
	mi.mark_del = true
	p.remove(mi)
	return true
}

func (p *PrioQueue) Rollback(mid, note string) bool {
	mi, ok := p.lockingMessages[mid]
	if !ok {
		return false
	}
	mi.GetTime = 0
	delete(p.grouplocks, mi.GroupId)
	delete(p.lockingMessages, mi.Id)

	return true
}
func (p *PrioQueue) Failure(mid, note string) bool {
	mi, ok := p.lockingMessages[mid]
	if !ok {
		return false
	}
	mi.Props["note"] = note
	mi.mark_del = true
	p.remove(mi)
	return true
}

func (p *PrioQueue) Delete(mid, note string) bool {
	mi, ok := p.messageids[mid]
	if !ok {
		return false
	}
	mi.Props["note"] = note
	mi.mark_del = true
	p.remove(mi)
	return true
}

func (p *PrioQueue) Size() int {
	return len(p.messageids)
}

func (p *PrioQueue) Clear() {
	// TODO 删除全部消息
}
