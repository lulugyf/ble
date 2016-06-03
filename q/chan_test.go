package q

import (
	"testing"
)

func TestChanBase(t *testing.T) {
	p := NewQueue("t1", "c1")
	QMap.Put(p)

	m := NewMessage(SEND_COMMIT)
	m.SetProp(TARGET_TOPIC, "t1")
	m.SetProp(GROUP, "group1")
	m.SetProp(PRIORITY, 4)
	m.SetProp(MESSAGE_ID, "11111")
	m.SetProp(CLIENT_ID, "producer11")
	t.Logf("made a request message.\n")

	mr := NewMessage(ANSWER)

	err := QMap.REQ_SendCommit(m, mr, "broker_id")
	if err != nil {
		t.Errorf("send failed, err: %v\n", err)
	}
	if mr.GetPropStr(RESULT_CODE) != OK {
		t.Errorf("get result code not ok: %s\n", mr.GetPropStr(RESULT_CODE))
	}

	if p.Size() != 1 {
		t.Errorf("send commit failed\n")
	}

	///////////////////////////
	m1 := NewMessage(PULL)
	m1.SetProp(TARGET_TOPIC, "t1")
	m1.SetProp(CLIENT_ID, "c1")
	m1.SetProp(PROCESSING_TIME, 20)
	m1r := NewMessage(ANSWER)
	err = QMap.REQ_Pull(m1, m1r, "broker_id")
	if err != nil {
		t.Errorf("pull failed")
	}
	if m1r.GetPropStr(RESULT_CODE) != OK {
		t.Errorf("get result code not ok: %s\n", m1r.GetPropStr(RESULT_CODE))
	}
	if m1r.GetPropStr(MESSAGE_ID) != "11111" {
		t.Errorf("message id mismatch")
	}
	mid := m1r.GetPropStr(MESSAGE_ID)

	///////////////////////////
	m2 := NewMessage(PULL)
	m2.SetProp(TARGET_TOPIC, "t1")
	m2.SetProp(CLIENT_ID, "c1")
	m2.SetProp(PULL_CODE, PULL_COMMIT)
	m2.SetProp(MESSAGE_ID, mid)
	m2r := NewMessage(ANSWER)

	err = QMap.REQ_Pull(m2, m2r, "broker_id")
	if err != nil {
		t.Errorf("pull failed")
	}
	if m2r.GetPropStr(RESULT_CODE) != OK {
		t.Errorf("get result code not ok: %s\n", m2r.GetPropStr(RESULT_CODE))
	}
}

func TestChanParameterMissing(t *testing.T) {
	p := NewQueue("t1", "c1")
	QMap.Put(p)

	m := NewMessage(SEND_COMMIT)
	//	m.SetProp(TARGET_TOPIC, "t1")
	m.SetProp(GROUP, "group1")
	m.SetProp(PRIORITY, 4)
	m.SetProp(MESSAGE_ID, "11111")
	m.SetProp(CLIENT_ID, "producer11")
	t.Logf("made a request message.\n")

	mr := NewMessage(ANSWER)

	err := QMap.REQ_SendCommit(m, mr, "broker_id")
	if err != nil {
		t.Errorf("REQ_SendCommit failed")
	}
	if mr.GetPropStr(RESULT_CODE) != REQUIRED_PARAMETER_MISSING {
		t.Errorf("get result code not REQUIRED_PARAMETER_MISSING: %s\n",
			mr.GetPropStr(RESULT_CODE))
	}

	///////////////////////////
	m1 := NewMessage(PULL)
	//	m1.SetProp(TARGET_TOPIC, "t1")
	m1.SetProp(CLIENT_ID, "c1")
	m1.SetProp(PROCESSING_TIME, 20)
	mr = NewMessage(ANSWER)
	err = QMap.REQ_Pull(m1, mr, "broker_id")
	if err != nil {
		t.Errorf("REQ_Pull failed")
	}
	if mr.GetPropStr(RESULT_CODE) != REQUIRED_PARAMETER_MISSING {
		t.Errorf("get result code not REQUIRED_PARAMETER_MISSING: %s\n",
			mr.GetPropStr(RESULT_CODE))
	}
}

func TestChanNoMore(t *testing.T) {
	p := NewQueue("t1", "c1")
	QMap.Put(p)

	m := NewMessage(SEND_COMMIT)
	m.SetProp(TARGET_TOPIC, "t1")
	m.SetProp(GROUP, "group1")
	m.SetProp(PRIORITY, 4)
	m.SetProp(MESSAGE_ID, "11111")
	m.SetProp(CLIENT_ID, "producer11")

	mr := NewMessage(ANSWER)

	err := QMap.REQ_SendCommit(m, mr, "broker_id")
	if err != nil {
		t.Errorf("send failed, err: %v\n", err)
	}
	if mr.GetPropStr(RESULT_CODE) != OK {
		t.Errorf("get result code not ok: %s\n", mr.GetPropStr(RESULT_CODE))
	}

	//pull first
	m1 := NewMessage(PULL)
	m1.SetProp(TARGET_TOPIC, "t1")
	m1.SetProp(CLIENT_ID, "c1")
	m1.SetProp(PROCESSING_TIME, 20)
	mr = NewMessage(ANSWER)
	err = QMap.REQ_Pull(m1, mr, "broker_id")
	if err != nil {
		t.Errorf("REQ_Pull failed")
	}
	if mr.GetPropStr(RESULT_CODE) != OK {
		t.Errorf("get result code not OK: %s\n",
			mr.GetPropStr(RESULT_CODE))
	}
	if mr.GetPropStr(MESSAGE_ID) != "11111" {
		t.Errorf("message-id mismatched")
	}

	// commit-and-next
	m1 = NewMessage(PULL)
	m1.SetProp(TARGET_TOPIC, "t1")
	m1.SetProp(CLIENT_ID, "c1")
	m1.SetProp(PROCESSING_TIME, 20)

	m1.SetProp(PULL_CODE, PULL_COMMIT_AND_NEXT)
	m1.SetProp(MESSAGE_ID, mr.GetPropStr(MESSAGE_ID))
	mr = NewMessage(ANSWER)
	err = QMap.REQ_Pull(m1, mr, "broker_id")
	if err != nil {
		t.Errorf("REQ_Pull failed")
	}
	if mr.GetPropStr(RESULT_CODE) != NO_MORE_MESSAGE {
		t.Errorf("get result code not NO_MORE_MESSAGE: %s\n",
			mr.GetPropStr(RESULT_CODE))
	}
}
