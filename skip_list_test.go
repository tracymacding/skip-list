package skip_list

import (
	"testing"
)

type IntComparer struct {
}

func (ic IntComparer) Less(key1, key2 interface{}) bool {
	return key1.(int) < key2.(int)
}

func (ic IntComparer) Equal(key1, key2 interface{}) bool {
	return key1.(int) == key2.(int)
}

func TestSkipList(t *testing.T) {
	skipList := NewSkipList(&IntComparer{})
	skipList.Insert(1)
	skipList.Insert(2)
	skipList.Insert(3)
	n := skipList.Find(1)
	// fmt.Printf("%+v\n", *n)
	skipList.Delete(1)
	n = skipList.Find(1)
	// fmt.Printf("%v\n", n)
	if n != nil {
		t.Fatal()
	}
	n = skipList.Find(2)
	// fmt.Printf("%v\n", n)
	if n == nil {
		t.Fatal()
	}
}
