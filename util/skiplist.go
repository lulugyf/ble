package util

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

type NodeLevel struct {
	level  int
	header *SKNode
}

type Skiplist struct {
	max_level int
	levels    []NodeLevel
}

func NewSkiplist(max_level int) *Skiplist {
	l := &Skiplist{max_level: max_level, levels: make([]NodeLevel, max_level)}

	for i := 0; i < max_level; i++ {
		nl := &l.levels[i]
		nl.level = max_level - i - 1
		n := &SKNode{key: MAX_KEY, value: nil, next: nil, down: nil}
		nl.header = &SKNode{key: MIN_KEY, value: nil, next: n, down: nil}
	}

	for i := 1; i < max_level; i++ {
		l.levels[i].header.down = l.levels[i-1].header
	}
	return l
}

func (sk *Skiplist) Insert(key int, value interface{}) {

}
