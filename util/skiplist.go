package util

import (
	"math/rand"
)

/*


level 3   -1       5         1
level 2   -1       5   7     1
level 1   -1   3   5   7     1
level 0   -1 2 3 4 5 6 7 8 9 1

一个跳表，应该具有以下特征：
    	一个跳表应该有几个层（level）组成；
    	跳表的第一层包含所有的元素；
    	每一层都是一个有序的链表；
    	如果元素x出现在第i层，则所有比i小的层都包含x；
    	第i层的元素通过一个down指针指向下一层拥有相同值的元素；
    	在每一层中，-1和1两个元素都出现(分别表示INT_MIN和INT_MAX)；
    	Top指针指向最高层的第一个元素。

*/

const (
	MIN_KEY = -1
	MAX_KEY = 1
)

type SKNode struct {
	key   int
	value interface{}
	next  *SKNode
	down  *SKNode
}

type Skiplist struct {
	max_level int
	header    *SKNode
}

func NewSkiplist(max_level int) *Skiplist {
	l := &Skiplist{max_level: max_level, header: nil}

	var n *SKNode = nil
	for i := 0; i < max_level; i++ {
		nt := &SKNode{key: MAX_KEY, value: nil, next: nil, down: nil}
		nh := &SKNode{key: MIN_KEY, value: nil, next: nt, down: nil}
		if n != nil {
			nh.down = n
			nt.down = n.next
		}
		n = nh
	}

	l.header = n
	return l
}

/*从那一层开始插入是取随机数的*/
func (sk *Skiplist) Insert(key int, value interface{}) {
	target_level := rand.Intn(sk.max_level)

	// 按查找方法找到指定层的插入位置
	l1 := sk.header
	var l2 *SKNode
	for level := sk.max_level - 1; level >= target_level; level-- {
		l2 = l1.next
		if l2.key == MAX_KEY || l2.key > key { //到达本level边界, 或者遇到更大的key
			l1 = l1.down   // 向下遍历
			if l1 == nil { //到底了
				break
			}
			continue
		}
		if l2.key == key { //找到的话， 就覆盖
			for l2 != nil {
				l2.value = value
				l2 = l2.down
			}
			return
		}
		l1 = l2
	}

	//在 l1 后面添加
	var l3 *SKNode = nil
	for {
		l2 = &SKNode{key: key, value: value, next: l1.next, down: nil}
		if l3 != nil {
			l3.down = l2
		}
		l1.next = l2
		l3 = l2
		l1 = l1.down
		if l1 == nil {
			break
		}
		// 在下一层找到插入位置 （注：有可能在下一层有这个key )
		for {
			l2 = l1.next
			if l2.key == key {
				for l2 != nil {
					l2.value = value
					l2 = l2.down
				}
				return
			} else if l2.key > key { //found
				break
			} else {
				l1 = l2
			}
		}
	}
}

/*查找顺序是从高level向低level找， 直到 0 level */
func (sk *Skiplist) Find(key int) *SKNode {
	l1 := sk.header
	var l2 *SKNode
	for {
		l2 = l1.next
		if l2.key == MAX_KEY || l2.key > key { //到达本level边界, 或者遇到更大的key
			l1 = l1.down   // 向下遍历
			if l1 == nil { //到底了
				break
			}
			continue
		}
		if l2.key == key {
			return l2 //found
		}
		l1 = l2
	}
	return nil //not found
}

func (sk *Skiplist) Delete(key int) {
	l1 := sk.header
	var l2 *SKNode
	for {
		l2 = l1.next
		if l2.key == MAX_KEY || l2.key > key { //到达本level边界, 或者遇到更大的key
			l1 = l1.down   // 向下遍历
			if l1 == nil { //到底了
				break
			}
			continue
		}
		if l2.key == key { //found
			break
		}
		l1 = l2
	}

	if l1 != nil && l1.next.key == key {
		for {
			l2 = l1.next
			l1.next = l2.next
			l2.down = nil
			l1 = l1.down
			if l1 == nil {
				break
			}
			// 在下一层找到删除位置
			for {
				l2 = l1.next
				if l2.key == key {
					break
				}
				l1 = l2
			}
		}
	}
}
