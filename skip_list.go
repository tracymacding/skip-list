package skip_list

import (
	"fmt"
	"math/rand"
	"time"
)

const (
	maxLevel = 12
	p        = 0.5
)

type Comparer interface {
	Less(v1, v2 interface{}) bool
	Equal(v1, v2 interface{}) bool
	Min() interface{}
}

type node struct {
	value    interface{}
	level    int
	forwards []*node
	drawed   bool
}

type skipList struct {
	head     *node
	level    int
	comparer Comparer
}

func NewNode(v interface{}, level int) *node {
	n := new(node)
	n.level = level
	n.value = v
	n.forwards = make([]*node, level)
	for i := 0; i < level; i++ {
		n.forwards[i] = nil
	}
	n.drawed = false
	return n
}

func (n *node) DrawMysel() {
	if !n.drawed {
		for i := 0; i < n.level; i++ {
			println("|")
			println("|")
			print("---")
			println("|")
			println("|")
		}
	}
}

func (n *node) LinkOther(other *node, level int, distance int) {

}

func NewSkipList(c Comparer) *skipList {
	s := new(skipList)
	// s.level = maxLevel
	s.level = 0
	s.head = NewNode(c.Min(), maxLevel)
	s.comparer = c
	return s
}

// TODO:
//func (s *skipList) Dump() {
//	for i := s.level - 1; i >= 0; i-- {
//		x := s.head
//		for x := s.head; x != nil; x = x.forwards[i] {
//			nodeDistance := s.distanceBetweenNode(x, x.forwards[i])
//			x.DrawMysel()
//			for i := 0; i < nodeDistance-1; i++ {
//				print("-")
//				print(" ")
//			}
//			print(">")
//			x.forwards[i].DrawMysel()
//			// 连接相邻两个节点
//			x.LinkOther(x.forwards[i], i, nodeDistance)
//		}
//	}
//}

func (s *skipList) Find(key interface{}) *node {
	x := s.head
	for i := s.level - 1; i >= 0; i-- {
		for {
			if x.forwards[i] != nil && s.comparer.Less(x.forwards[i].value, key) {
				x = x.forwards[i]
			} else {
				break
			}
		}
	}
	if s.comparer.Equal(key, x.forwards[0].value) {
		return x.forwards[0]
	}
	return nil
}

func (s *skipList) Insert(key interface{}) {

	update := make([]*node, maxLevel)
	x := s.head
	for i := maxLevel - 1; i >= 0; i-- {
		for {
			if x.forwards[i] != nil && s.comparer.Less(x.forwards[i].value, key) {
				x = x.forwards[i]
			} else {
				break
			}
		}
		update[i] = x
		// fmt.Printf("update[%d]: %v\n", i, x)
	}

	if x.forwards[0] != nil && s.comparer.Equal(key, x.forwards[0].value) {
		x.forwards[0].value = key
		return
	}

	level := genRandomLevel()
	fmt.Printf("new level: %d\n", level)
	n := NewNode(key, level)
	if level > s.level {
		//for i := s.level; i < level; i++ {
		//	s.head.forwards[i] = n
		//}
		s.level = level
	}
	for i := 0; i < level; i++ {
		n.forwards[i] = update[i].forwards[i]
		// fmt.Printf("node [%d], level [%d], forward:%v\n", n.value, i, update[i].forwards[i])
		update[i].forwards[i] = n
		// fmt.Printf("node [%d], level [%d], forward:%v\n", update[i].value, i, n.value)
	}
}

func (s *skipList) Delete(key interface{}) {
	update := make([]*node, maxLevel)
	x := s.head
	for i := s.level - 1; i >= 0; i-- {
		for {
			if x.forwards[i] != nil && s.comparer.Less(x.forwards[i].value, key) {
				x = x.forwards[i]
				// } else if x.forwards[i] != nil && s.comparer.Equal(key, x.forwards[i].value) {
				// update[i] = x
				// break
			} else {
				break
			}
		}
		update[i] = x
		// fmt.Printf("update[%d]: %v\n", i, x)
	}
	if x.forwards[0] != nil && !s.comparer.Equal(x.forwards[0].value, key) {
		return
	}
	n := x.forwards[0]

	for i := 0; i < s.level; i++ {
		if update[i].forwards[i] != nil && s.comparer.Equal(key, update[i].forwards[i].value) {
			update[i].forwards[i] = n.forwards[i]
		}
	}
	for {
		if s.level > 0 && s.head.forwards[s.level-1] == nil {
			s.level -= 1
		} else {
			break
		}
	}
	fmt.Printf("skip list level now is: %d\n", s.level)
}

func genRandomLevel() int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	i := 1
	for ; i <= maxLevel; i++ {
		if r.Float32() >= float32(p) {
			break
		}
	}
	return i
}
