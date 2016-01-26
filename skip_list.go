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
	value interface{}
	// forwards []*node
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
	// n.forwards = make([]*node, level)
	n.forwards = make([]unsafe.Pointer, level)
	for i := 0; i < level; i++ {
		n.forwards[i] = nil
	}
	return n
}

func (n *node) SetNext(level int, n1 *node) {
	atomic.StorePointer(&(n.forwards[level]), unsafe.Pointer(n1))
	// n.forwards[level] = n1
}

func (n *node) Next(level int) *node {
	// return n.forwards[level]
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
	//s.RLock()
	//defer s.RUnlock()

	//x := s.head
	//for i := s.level - 1; i >= 0; i-- {
	//	for {
	//		if x.forwards[i] != nil && s.comparer.Less(x.forwards[i].value, key) {
	//			x = x.forwards[i]
	//		} else {
	//			break
	//		}
	//	}
	//}

	n := s.findGreaterOrEqual(key, nil)
	//if s.comparer.Equal(key, x.forwards[0].value) {
	//	return x.forwards[0]
	//}
	if s.comparer.Equal(key, n.value) {
		return n
	}
	return nil
}

func (s *skipList) Insert(key interface{}) {
	s.Lock()
	defer s.Unlock()

	//update := make([]*node, maxLevel)
	//x := s.head
	//for i := maxLevel - 1; i >= 0; i-- {
	//	for {
	//		if x.forwards[i] != nil && s.comparer.Less(x.forwards[i].value, key) {
	//			x = x.forwards[i]
	//		} else {
	//			break
	//		}
	//	}
	//	update[i] = x
	//}
	prev := make([]*node, maxLevel)
	n := s.findGreaterOrEqual(key, prev)

	//if x.forwards[0] != nil && s.comparer.Equal(key, x.forwards[0].value) {
	//	x.forwards[0].value = key
	//	return
	//}
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

	//update := make([]*node, maxLevel)
	//x := s.head
	//for i := s.level - 1; i >= 0; i-- {
	//	for {
	//		if x.forwards[i] != nil && s.comparer.Less(x.forwards[i].value, key) {
	//			x = x.forwards[i]
	//		} else {
	//			break
	//		}
	//	}
	//	update[i] = x
	//}
	//if x.forwards[0] != nil && !s.comparer.Equal(x.forwards[0].value, key) {
	//	return
	//}
	//n := x.forwards[0]
	prev := make([]*node, maxLevel)
	n := s.findGreaterOrEqual(key, prev)
	if !s.comparer.Equal(n.value, key) {
		return
	}

	level := s.currentLevel()
	for i := 0; i < level; i++ {
		if prev[i] != nil && prev[i].Next(i) != nil && s.comparer.Equal(key, prev[i].Next(i).value) {
			prev[i].SetNext(i, n.Next(i))
		}
		//if update[i].forwards[i] != nil && s.comparer.Equal(key, update[i].forwards[i].value) {
		//	update[i].forwards[i] = n.forwards[i]
		//}
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
