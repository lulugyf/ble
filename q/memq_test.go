package q

import (
	"fmt"
	"log"
	"testing"
	//	"time"
)

func sendone(p *PrioQueue, id, group, data string, priority int) bool {
	x := NewMsgIndex(id)
	x.Priority = priority
	x.GroupId = group
	x.Data = data
	return p.Push(x)
}

// for group and priority sort consume
func TestBase(t *testing.T) {
	p := NewQueue("t1", "c1")
	for i := 0; i < 10; i++ {
		sendone(p, fmt.Sprintf("id11.%d", i), fmt.Sprintf("groupid-%d", i%3),
			fmt.Sprintf("xxx--%d", i), 100+i)
	}

	if p.Size() != 10 {
		t.Errorf("Push size failed\n")
	}

	first_id := p.Pull(10).Id
	p.Commit(first_id, "")
	last_id := ""
	for {
		mi := p.Pull(60)
		if mi != nil {
			last_id = mi.Id
			log.Printf("====PULL id:%s data:%s\n", mi.Id, mi.Data)
			p.Commit(mi.Id, "ok")
			//			log.Printf("commit: %s return %t\n", mi.Id, r)
		} else {
			log.Printf("not found\n")
			break
		}
	}
	if p.Size() != 0 {
		t.Errorf("last size error: %d\n", p.Size())
	}
	if first_id != "id11.0" || last_id != "id11.7" {
		t.Errorf("priority and group sort failed\n")
	}
}

/*
func TestTimeout(t *testing.T) {
	p := NewQueue("t1", "c1")
	if !sendone(p, "id1111", "group11", "data11", 100) {
		t.Errorf("send failed")
	}

	mi := p.Pull(1)
	id1 := mi.Id
	t.Printf("id:%s retry: %d\n", mi.Id, mi.Retry)
	time.Sleep(time.Second * 2)
	mi = p.Pull(1)
	if mi.Id != id1 {
		t.Errorf("pull after timeout, id mismatched\n")
	}

	log.Printf("id:%s retry: %d\n", mi.Id, mi.Retry)
	if mi.Retry != 2 {
		t.Errorf("retry number error\n")
	}

	p.Commit(mi.Id)
	if p.Size() != 0 {
		t.Errorf("commit failed\n")
	}
} */

func TestToomanyRetries(t *testing.T) {
	p := NewQueue("t1", "c1")
	if !sendone(p, "id1111", "group11", "data11", 100) {
		t.Errorf("send failed")
	}
	c := 0
	c1 := 0
	for {
		mi := p.Pull(10)
		if mi == nil {
			break
		}
		c++
		c1 = mi.Retry
		p.Rollback(mi.Id, "rollback")
	}
	t.Logf("last retry=%d\n", c1)
	if c != p.max_tries {
		t.Errorf("test max retries failed, want %d but %d\n", p.max_tries, c)
	}
}
