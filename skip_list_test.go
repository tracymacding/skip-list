package skip_list

import (
	"fmt"
	"math/rand"
	"sync"
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
	fmt.Printf("Thread %d insert records [%d - %d] cost %f s, tps: %f\n", threadID, keyFrom, keyTo, cost.Seconds(), float64((keyTo-keyFrom))/cost.Seconds())
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
	fmt.Printf("Thread %d get records [%d - %d] cost %f s, tps: %f\n", threadID, keyFrom, keyTo, cost.Seconds(), float64((keyTo-keyFrom))/cost.Seconds())
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
	for i := 0; i < 10; i++ {
		wg.Wrap(SingleThreadInsert, skipList, i*1000000, (i+1)*1000000, i)
	}
	for i := 0; i < 0; i++ {
		wg.Wrap(SingleThreadGet, skipList, (i)*100000, (i+1)*100000, i)
	}
	wg.Wait()
}

type record struct {
	key     int
	deleted bool
}

type records struct {
	data map[int]*record
	sync.Mutex
	rand *rand.Rand
}

func newRecords() *records {
	return &records{
		data: make(map[int]*record),
		rand: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (r *records) pickUpOne() *record {
	r.Lock()
	defer r.Unlock()

	i := 0
	for k, record := range r.data {
		if i == rand.Intn(len(r.data)) {
			delete(r.data, k)
			return record
		}
	}
	return nil
}

func (r *records) insert(record *record) {
	r.Lock()
	defer r.Unlock()

	r.data[record.key] = record
}

func DeleteRecord(args []interface{}) {
	r := args[0].(*records)
	s := args[1].(*skipList)
	for {
		record := r.pickUpOne()
		if record != nil {
			s.Delete(record.key)
			record.deleted = true
			r.insert(record)
		}
	}
}

func GetAndCheck(args []interface{}) {
	r := args[0].(*records)
	s := args[1].(*skipList)
	for {
		record := r.pickUpOne()
		if record != nil {
			n := s.Find(record.key)
			if record.deleted && n != nil {
				panic("delete key find ok")
			}
			if !record.deleted && n == nil {
				panic("exist key not found")
			}
			r.insert(record)
		}
	}
}

func PutRecord(args []interface{}) {
	r := args[0].(*records)
	s := args[1].(*skipList)
	for {
		key := time.Now().Nanosecond()
		s.Insert(key)
		record := &record{
			key:     key,
			deleted: false,
		}
		r.insert(record)
	}
}

func TestCorrectness(t *testing.T) {
	println("correctness test begin")
	defer println("correctness test end")
	skipList := NewSkipList(&IntComparer{})
	SingleThreadInsert([]interface{}{skipList, 0, 1000000, 0})

	wg := new(WaitGroupWrapper)
	r := newRecords()
	for i := 0; i < 10; i++ {
		wg.Wrap(GetAndCheck, r, skipList)
	}
	for i := 0; i < 10; i++ {
		wg.Wrap(DeleteRecord, r, skipList)
	}
	for i := 0; i < 10; i++ {
		wg.Wrap(PutRecord, r, skipList)
	}
	// wg.Wait()
	time.Sleep(time.Second * 100)
}
