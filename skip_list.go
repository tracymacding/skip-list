package skip_list

import (
	//"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"
)

const (
	maxLevel = 12
	p        = 0.25
)

type Comparer interface {
	Less(v1, v2 interface{}) bool
	Equal(v1, v2 interface{}) bool
}

type node struct {
	value    interface{}
	forwards []unsafe.Pointer
}

type skipList struct {
	sync.RWMutex
	head     *node
	level    int
	comparer Comparer
	rand     *rand.Rand
}

func NewNode(v interface{}, level int) *node {
	n := new(node)
	n.value = v
	n.forwards = make([]unsafe.Pointer, level)
	for i := 0; i < level; i++ {
		n.forwards[i] = nil
	}
	return n
}

func (n *node) SetNext(level int, n1 *node) {
	atomic.StorePointer(&(n.forwards[level]), unsafe.Pointer(n1))
}

func (n *node) Next(level int) *node {
	return (*node)(atomic.LoadPointer(&n.forwards[level]))
}

func NewSkipList(c Comparer) *skipList {
	s := new(skipList)
	s.level = 0
	s.head = NewNode(0, maxLevel)
	s.comparer = c
	source := rand.NewSource(time.Now().UnixNano())
	s.rand = rand.New(source)
	return s
}

func (s *skipList) keyIsAfterNode(key interface{}, node *node) bool {
	if node != nil && s.comparer.Less(node.value, key) {
		return true
	}
	return false
}

func (s *skipList) currentLevel() int {
	return s.level
}

func (s *skipList) findGreaterOrEqual(key interface{}, prev []*node) *node {
	x := s.head
	level := s.currentLevel()

	for i := level - 1; i >= 0; i-- {
		for {
			next := x.Next(i)
			if s.keyIsAfterNode(key, next) {
				x = next
			} else {
				if prev != nil {
					prev[i] = x
				}
				if i == 0 {
					return next
				} else {
					break
				}
			}
		}
	}
	return nil
}

func (s *skipList) Find(key interface{}) *node {

	n := s.findGreaterOrEqual(key, nil)
	if n != nil && s.comparer.Equal(key, n.value) {
		return n
	}
	return nil
}

func (s *skipList) Insert(key interface{}) {
	s.Lock()
	defer s.Unlock()

	prev := make([]*node, maxLevel)
	n := s.findGreaterOrEqual(key, prev)

	if n != nil && s.comparer.Equal(key, n.value) {
		n.value = key
		return
	}

	level := s.genRandomLevel()
	n = NewNode(key, level)
	if level > s.level {
		for i := s.currentLevel(); i < level; i++ {
			prev[i] = s.head
		}
		s.level = level
	}
	for i := 0; i < level; i++ {
		n.SetNext(i, prev[i].Next(i))
		// fmt.Printf("node [%d], level [%d], forward:%v\n", n.value, i, prev[i].Next(i))
		prev[i].SetNext(i, n)
		//fmt.Printf("node [%d], level [%d], forward:%+v\n", prev[i].value, i, n.value)
	}
}

func (s *skipList) Delete(key interface{}) {
	s.Lock()
	defer s.Unlock()

	prev := make([]*node, maxLevel)
	n := s.findGreaterOrEqual(key, prev)
	if n != nil && !s.comparer.Equal(n.value, key) {
		return
	}

	level := s.currentLevel()
	for i := 0; i < level; i++ {
		if prev[i] != nil && prev[i].Next(i) != nil && s.comparer.Equal(key, prev[i].Next(i).value) {
			prev[i].SetNext(i, n.Next(i))
		}
	}
	for {
		if s.level > 0 && s.head.Next(s.level-1) == nil {
			s.level -= 1
		} else {
			break
		}
	}
}

func (s *skipList) genRandomLevel() int {
	i := 1
	for ; i < maxLevel; i++ {
		if s.rand.Float32() >= float32(p) {
			break
		}
	}
	return i
}
