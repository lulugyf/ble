package q

import (
	"fmt"
	"testing"
)

func TestList1(t *testing.T) {
	l := newList()
	l.pushBack("one")
	if l.Size() != 1 {
		t.Errorf("size error: %d\n", l.Size())
	}
	l.remove(l.head)
	if l.Size() != 0 {
		t.Errorf("size error after remove only one\n")
	}
}

func TestList2(t *testing.T) {
	l := newList()
	l.pushBack("one")
	l.pushBack("two")
	if l.Size() != 2 {
		t.Errorf("size error: %d\n", l.Size())
	}
	x := l.head
	l.remove(x)
	if l.Size() != 1 || l.head.value.(string) != "two" {
		t.Errorf("remove head failed\n")
	}
	l.remove(l.tail)
	if l.Size() != 0 {
		t.Errorf("remove tail failed\n")
	}
}

func TestList3(t *testing.T) {
	l := newList()
	l.pushBack("one")
	l.pushBack("two")
	l.pushBack("tree")
	if l.Size() != 3 {
		t.Errorf("size error: %d\n", l.Size())
	}
	l.remove(l.head.next)
	if l.Size() != 2 || l.head.value.(string) != "one" || l.tail.value.(string) != "tree" {
		t.Errorf("remove middle failed\n")
	}
}

func TestSortList1(t *testing.T) {
	p := NewSortList()
	p.Insert(5, "five")
	p.Insert(6, "six")
	p.Insert(7, "seven")
	p.Insert(2, "two")
	p.Insert(3, "three")
	p.Insert(3, "three 1")

	if p.sizep != 5 || p.sizec != 6 {
		t.Errorf("size error after insert\n")
	}

	p.Iter(func(key int, groupid string) int {
		//fmt.Printf("%d==%s\n", key, groupid)
		if groupid == "six" {
			return -1
		} else {
			return 1
		}
	})
	if p.sizep != 5 || p.sizec != 5 {
		t.Errorf("size error after insert, sizep:%d sizec:%d\n", p.sizep, p.sizec)
	}

	c := 0
	p.Iter(func(key int, groupid string) int {
		fmt.Printf("%d==%s\n", key, groupid)
		c++
		return 1
	})
	if c != 5 || p.sizep != 4 {
		t.Errorf("Iter count failed, sizep:%d, c:%d\n", p.sizep, c)
	}
}
