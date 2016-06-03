package q

/*
报文结构说明:
MESSAGE{
  HEADER {
	ALL_SIZE        (8bytes)  报文总长度，10进制字符串，左补'0'
	PROPERTIES_SIZE (4bytes)  属性字段长度，10进制字符串，左补'0'
	MESSAGE_TYPE    (2bytes)  消息类型
	RESERVE         (2bytes)  保留
	}
  PROPERTIES ($PROPERTIES_SIZE)  属性，json格式字符串
  CONTENT  ($ALL_SIZE - 16 - $PROPERTIES_SIZE bytes)
}

*/

import (
	//	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
)

type Message struct {
	Type    int
	Props   map[string]interface{}
	Content []byte
}

func NewMessage(tp int) *Message {
	return &Message{Type: tp, Props: make(map[string]interface{})}
}
func (m *Message) GetProp(key string) interface{} {
	if v, ok := m.Props[key]; ok {
		return v
	}
	return nil
}
func (m *Message) GetPropInt(key string) int {
	v := m.GetProp(key)
	if v == nil {
		return -1
	}
	if i, ok := v.(int); ok {
		return i
	} else if s, ok := v.(string); ok {
		i := -1
		fmt.Sscanf("%d", s, &i)
		return i
	}
	return -1
}
func (m *Message) GetPropStr(key string) string {
	v := m.GetProp(key)
	if v == nil {
		return ""
	}
	if s, ok := v.(string); ok {
		return s
	}
	return ""
}
func (m *Message) SetProp(key string, value interface{}) {
	m.Props[key] = value
}

const (
	HEAD_LEN = 16
)

func RecvMessage(r io.Reader) (*Message, error) {
	head := make([]byte, HEAD_LEN)
	n, err := io.ReadFull(r, head)
	if err != nil || n != len(head) {
		log.Printf("readFull header failed %v or bytes not enough: %d\n", err, n)
		return nil, err
	}
	all_len, _ := strconv.Atoi(string(head[0:8]))
	props_len, _ := strconv.Atoi(string(head[8:12]))
	mtype, _ := strconv.Atoi(string(head[12:14]))

	if mtype == BREAKHEART {
		log.Printf("it's a heartbeat\n")
		return RecvMessage(r)
	}

	props_data := make([]byte, props_len)
	n, err = io.ReadFull(r, props_data)
	if err != nil || n != props_len {
		log.Printf("readFull propertes failed %v or bytes not enough: %d\n", err, n)
		return nil, errors.New("read properties failed")
	}

	content_len := all_len - HEAD_LEN - props_len
	content := make([]byte, content_len)
	if content_len > 0 {
		n, err = io.ReadFull(r, content)
		if err != nil || n != len(content) {
			log.Printf("readFull content failed %v or bytes not enough: %d\n", err, n)
			return nil, errors.New("read content failed")
		}
	}
	var props map[string]interface{}
	err = json.Unmarshal(props_data, &props)
	if err != nil {
		log.Printf("parse properties failed: %v\n", err)
		return nil, errors.New("parse property json failed")
	}

	return &Message{Type: mtype, Props: props, Content: content}, nil
}

func SendMessage(w io.Writer, msg *Message) error {
	props_data, err := json.Marshal(msg.Props)
	if err != nil {
		log.Printf("marshal properties failed: %v\n", err)
		return err
	}
	total_len := HEAD_LEN + len(props_data) + len(msg.Content)
	head := fmt.Sprintf("%08d%04d%02d00", total_len, len(props_data), msg.Type)
	n, err := w.Write([]byte(head))
	if err != nil || n != HEAD_LEN {
		log.Printf("write header failed: err=%v n=%d\n", err, n)
		return errors.New("write header failed")
	}
	n, err = w.Write(props_data)
	if err != nil || n != len(props_data) {
		log.Printf("write props failed: err=%v n=%d\n", err, n)
		return errors.New("write props failed")
	}
	if msg.Content != nil && len(msg.Content) > 0 {
		n, err = w.Write(msg.Content)
		if err != nil || n != len(msg.Content) {
			log.Printf("write content failed: err=%v n=%d\n", err, n)
			return errors.New("write content failed")
		}
	}
	return nil
}

func HandleConnection(conn net.Conn, qm *qmap) {
	defer conn.Close()

	broker_id := ""
	remote_addr := conn.RemoteAddr()

	log.Printf("connection from %s\n", remote_addr)
	for {
		msg, err := RecvMessage(conn)
		msgr := NewMessage(ANSWER)
		if err != nil {
			log.Printf("recvMessage from %s failed: %v\n", remote_addr, err)
			break
		}

		switch msg.Type {
		case QUERY:
			broker_id = msg.GetPropStr(CLIENT_ID)
			log.Printf("query, client-id:[%s]\n", broker_id)
			msgr.SetProp(RESULT_CODE, OK)
		case SEND_COMMIT:
			log.Printf("send-commit\n")
			err = qm.REQ_SendCommit(msg, msgr, broker_id)
		case PULL:
			log.Printf("pull\n")
			err = qm.REQ_Pull(msg, msgr, broker_id)
		case DELETE:
			err = qm.REQ_Delete(msg, msgr, broker_id)
		case UNLOCK:
			err = qm.REQ_Unlock(msg, msgr, broker_id)
		default:
			msgr.SetProp(RESULT_CODE, UNSUPPORTED_MESSAGE_TYPE)
		}

		if err != nil {
			msgr.SetProp(RESULT_CODE, INTERNAL_SERVER_ERROR)
		}
		if err = SendMessage(conn, msgr); err != nil {
			log.Printf("send message failed: %v\n", err)
			break
		}
	}
	fmt.Printf("%s disconnected", remote_addr)
}

func main1() {
	fmt.Println("start")
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		// handle error
	}
	for {
		conn, err := ln.Accept() // this blocks until connection or error
		if err != nil {
			// handle error
			continue
		}
		go func(conn net.Conn) {
			HandleConnection(conn, QMap)
		}(conn) // a goroutine handles conn so that the loop can accept other connections
	}
}
