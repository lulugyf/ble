package q

import (
	"fmt"
	//	"log"
)

// 一个简单的双向链表
type listItem struct {
	value interface{}
	prev  *listItem
	next  *listItem
}
type List struct {
	head *listItem
	tail *listItem
	size int
}

func newList() *List {
	return &List{head: nil, tail: nil, size: 0}
}
func (l *List) Size() int {
	return l.size
}
func (l *List) insertBefore(lx *listItem, value interface{}) {
	px := &listItem{value: value, next: lx, prev: lx.prev}
	if lx.prev == nil {
		l.head = px
	} else {
		lx.prev.next = px
	}
	lx.prev = px
	l.size++
}

func (l *List) pushBack(value interface{}) {
	px := &listItem{value: value, prev: nil, next: nil}
	if l.tail == nil {
		// empty
		l.head = px
		l.tail = px
	} else {
		p1 := l.tail
		if p1.prev == nil {
			// only on element
			l.tail = px
			l.head.next = px
			px.prev = l.head
		} else {
			// more than one
			p1.next = px
			px.prev = p1
			l.tail = px
		}
	}
	l.size++
}
func (l *List) remove(p1 *listItem) *listItem {
	if p1 == nil {
		return nil
	}
	x := p1.next
	if p1.prev == nil {
		l.head = p1.next
	} else {
		p1.prev.next = p1.next
	}
	if p1.next == nil {
		l.tail = p1.prev
	} else {
		p1.next.prev = p1.prev
	}
	l.size--
	return x
}

////////////////////////////////
// 一个带排序的双向链表
type PrioItem struct {
	p    int
	data *List
	prev *PrioItem
	next *PrioItem
}
type SortList struct {
	head  *PrioItem
	tail  *PrioItem
	sizep int //key 节点数
	sizec int //元素个数
}

func NewSortList() *SortList {
	return &SortList{head: nil, tail: nil, sizep: 0, sizec: 0}
}
func (p *SortList) Insert(key int, groupid string) bool {
	p1 := p.head
	for ; p1 != nil; p1 = p1.next {
		if key == p1.p {
			// found match key
			if groupid != p1.data.tail.value.(string) {
				p1.data.pushBack(groupid)
				p.sizec++
			}
			break
		} else if key > p1.p {
			continue
		} else {
			// insert before p1
			//			log.Printf("insert %d before %d\n", key, p1.p)
			px := &PrioItem{p: key, data: newList(), next: p1, prev: p1.prev}
			px.data.pushBack(groupid)
			if p1.prev == nil {
				p.head = px
			} else {
				p1.prev.next = px
			}
			p1.prev = px
			p.sizec++
			p.sizep++
			break
		}
	}
	if p1 == nil {
		// insert tail
		//		log.Printf("==insert tail: %d\n", key)
		px := &PrioItem{p: key, data: newList(), prev: nil, next: nil}
		px.data.pushBack(groupid)
		if p.tail == nil {
			// empty
			p.head = px
			p.tail = px
		} else {
			p1 = p.tail
			if p1.prev == nil {
				// only on element
				p.tail = px
				p.head.next = px
				px.prev = p.head
			} else {
				// more than one
				p1.next = px
				px.prev = p1
				p.tail = px
			}
		}
		p.sizec++
		p.sizep++

	}
	return true
}

func (p *SortList) Remove(key int) interface{} {
	for p1 := p.head; p1 != nil; p1 = p1.next {
		if p1.p == key {
			data := p1.data

			if p1.prev == nil {
				p.head = p1.next
			} else {
				p1.prev.next = p1.next
			}
			if p1.next == nil {
				p.tail = p1.prev
			} else {
				p1.next.prev = p1.prev
			}
			p.sizep--
			p.sizec -= data.size

			return data
		} else if p1.p > key {
			return nil
		}
	}
	return nil
}

func (p *SortList) Find(key int) interface{} {
	for p1 := p.head; p1 != nil; p1 = p1.next {
		if p1.p == key {
			return p1.data
		} else if p1.p > key {
			return nil
		}
	}
	return nil
}

//删除节点， 并返回前一个元素
func (l *SortList) remove(p1 *PrioItem) *PrioItem {
	if p1 == nil {
		return nil
	}
	x := p1.prev

	if p1.prev == nil {
		l.head = p1.next
	} else {
		p1.prev.next = p1.next
	}
	if p1.next == nil {
		l.tail = p1.prev
	} else {
		p1.next.prev = p1.prev
	}
	l.sizep--
	l.sizec -= p1.data.size

	return x
}

/* F return value: 1-continue 0-stop -1-remove & continue */
func (p *SortList) Iter(F func(int, string) int) {
LOOP:
	for p1 := p.tail; p1 != nil; {
		if p1.data.Size() == 0 {
			// 该key下无数据了， 删除节点
			p1 = p.remove(p1)
			continue
		}
		for x := p1.data.head; x != nil; {
			r := F(p1.p, x.value.(string))
			if r == 0 {
				break LOOP
			} else if r == -1 {
				x = p1.data.remove(x)
				p.sizec--
			} else if r == 1 {
				x = x.next
			} else {
				break LOOP
			}
		}
		p1 = p1.prev
	}
}

func Prio_test() {
	p := NewSortList()
	p.Insert(5, "five")
	p.Insert(6, "six")
	p.Insert(7, "seven")
	p.Insert(2, "two")
	p.Insert(3, "three")
	p.Insert(3, "three 1")

	p.Iter(func(key int, groupid string) int {
		//fmt.Printf("%d==%s\n", key, groupid)
		if groupid == "six" {
			return -1
		} else {
			return 1
		}
	})

	p.Iter(func(key int, groupid string) int {
		fmt.Printf("%d==%s\n", key, groupid)
		return 1
	})

	//	fmt.Printf("==%s\n", p.Find(7).(string))

	//	p.Remove(3)
	//	//	p.Iter()
	//	p.Remove(2)
	//	p.Remove(7)
	//	//	p.Iter()

	//	p.Remove(5)
	////	p.Iter()
	//	p.Remove(6)
	//	p.Iter()

	//fmt.Printf("==%s\n", p.Find(2).(string))
}
