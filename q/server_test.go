package q

import (
	"bytes"
	//	"io"
	"testing"
)

func TestMsgProperties(t *testing.T) {
	m := NewMessage(QUERY)
	m.SetProp(CLIENT_ID, "broker_id11")

	if m.GetProp(CLIENT_ID).(string) != "broker_id11" {
		t.Errorf("set and get string failed\n")
	}
	m.SetProp(PRIORITY, 1000)
	if m.GetProp(PRIORITY).(int) != 1000 {
		t.Error("set and get Int failed")
	}

	if m.GetProp(PRIORITYNAME) != nil {
		t.Error("GetProp with none-exists property failed")
	}

	if m.GetPropStr(PRIORITYNAME) != "" {
		t.Error("GetPropStr with none-exists property failed")
	}
}

func TestMsgIO(t *testing.T) {
	m := NewMessage(SEND)
	m.SetProp(CLIENT_ID, "pub11")
	m.SetProp(GROUP, "hello.group")
	m.Content = []byte("there have some data in your brain.")

	buf := bytes.NewBuffer(make([]byte, 0))
	t.Logf("buf.len=%d\n", buf.Len())
	err := SendMessage(buf, m)
	if err != nil {
		t.Errorf("write message failed:%v\n", err)
	}
	t.Logf("buf.len=%d, str=[%s]\n", buf.Len(), string(buf.Bytes()))
	//	b := make([]byte, 16)
	//	io.ReadFull(buf, b)
	//	t.Logf("==head[%s]\n", string(b))

	m1, err := RecvMessage(buf)
	if err != nil {
		t.Errorf("recvMessage failed: %v\n", err)
	} else {
		if m1.GetProp(GROUP).(string) != "hello.group" {
			t.Errorf("property from recv message failed")
		}
	}
}
