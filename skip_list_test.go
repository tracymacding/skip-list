package skip_list

import (
	"fmt"
	"testing"
	"time"
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
	skipList.Delete(1)
	n = skipList.Find(1)
	if n != nil {
		t.Fatal()
	}
	n = skipList.Find(2)
	if n == nil {
		t.Fatal()
	}
}

func SingleThreadInsert(args []interface{}) {
    s := args[0].(*skipList)
    keyFrom := args[1].(int)
    keyTo := args[2].(int)
    threadID := args[3].(int)

    start := time.Now()
    for i := keyFrom; i < keyTo; i++ {
        s.Insert(i)
    }
    cost := time.Since(start)
    fmt.Printf("Thread %d insert records [%d - %d] cost %f s\n", threadID, keyFrom, keyTo, cost.Seconds())
}

func SingleThreadGet(args []interface{}) {
    s := args[0].(*skipList)
    keyFrom := args[1].(int)
    keyTo := args[2].(int)
    threadID := args[3].(int)

    start := time.Now()
    for i := keyFrom; i < keyTo; i++ {
        s.Find(i)
    }
    cost := time.Since(start)
    fmt.Printf("Thread %d get records [%d - %d] cost %f s, tps: %f\n", threadID, keyFrom, keyTo, cost.Seconds(), float64((keyTo - keyFrom))/cost.Seconds())
}

func TestSkipListPerformance(t *testing.T) {
    skipList := NewSkipList(&IntComparer{})
    wg := new(WaitGroupWrapper)
    for i := 0; i < 1; i++ {
        wg.Wrap(SingleThreadInsert, skipList, i*1000000, (i+1)*1000000, i)
    }
    wg.Wait()
}

func TestSkipListRWPerformance(t *testing.T) {
    skipList := NewSkipList(&IntComparer{})
	SingleThreadInsert([]interface{}{skipList, 0, 1000000, 0})
	
    wg := new(WaitGroupWrapper)
    for i := 0; i < 1; i++ {
        wg.Wrap(SingleThreadInsert, skipList, (i+1)*1000000, (i+2)*1000000, i)
    }
    for i := 0; i < 2; i++ {
        wg.Wrap(SingleThreadGet, skipList, (i)*500000, (i+1)*500000, i)
    }
	wg.Wait()
}
