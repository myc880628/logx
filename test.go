package main

import (
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

const (
	cacheLineSize = 64
)

type (
	Queue struct {
		size      int32
		_         [cacheLineSize - 4]byte
		head      int32
		capacity0 int32
		_         [cacheLineSize - 8]byte
		tail      int32
		capacity1 int32
		_         [cacheLineSize - 8]byte
		state     []int32
		elements  []interface{}
	}
)

func NewQueue(capacity int32) *Queue {
	if capacity < 0 {
		panic("NewQueue error: negatice capacity")
	}
	return &Queue{
		capacity0: capacity,
		capacity1: capacity,
		state:     make([]int32, capacity),
		elements:  make([]interface{}, capacity),
	}
}

var (
	a, b, c int
)

func (q *Queue) Push(element interface{}) bool {
	for q.size < q.capacity1 {
		tail := q.tail
		if !atomic.CompareAndSwapInt32(&q.state[tail], 0, 1) {
			runtime.Gosched()
			continue
		}
		if !atomic.CompareAndSwapInt32(&q.tail, tail, (tail+1)%q.capacity1) {
			atomic.StoreInt32(&q.state[tail], 0)
			continue
		}
		q.elements[tail] = element
		atomic.StoreInt32(&q.state[tail], 2)
		atomic.AddInt32(&q.size, 1)
		return true
	}
	return false
}

var (
	x, y, z int
)

var (
	count int32
)

func (q *Queue) Pop() (interface{}, bool) {
	for q.size > 0 {
		head := q.head
		if !atomic.CompareAndSwapInt32(&q.state[head], 2, 3) {
			runtime.Gosched()
			x++
			continue
		}
		if q.head != head {
			y++
			atomic.StoreInt32(&q.state[head], 2)
			continue
		}
		element := q.elements[head]
		atomic.StoreInt32(&q.head, (head+1)%q.capacity0)
		atomic.StoreInt32(&q.state[head], 0)
		atomic.AddInt32(&q.size, -1)
		return element, true
	}
	return nil, false
}

const (
	NumGroutines = 1000
	NumOpertions = 10000
)

func main() {
	q := NewQueue(NumGroutines * NumOpertions)

	// q := &list.List{}
	// q = q.Init()
	// var set map[int32]bool = make(map[int32]bool, 10000)
	// var mu sync.Mutex
	// ***********

	starTime := time.Now().UnixMilli()
	var wg sync.WaitGroup
	for i := 0; i < NumGroutines; i++ {
		wg.Add(1)
		go func() {
			defer func() {
				wg.Done()
			}()
			for j := 0; j < NumOpertions; j++ {
				v := atomic.AddInt32(&count, 1)
				if ok := q.Push(v); !ok {
					println("XXX")
				}
				// mu.Lock()
				// q.PushBack(j)
				// mu.Unlock()
			}
		}()
	}

	for i := 0; i < NumGroutines; i++ {
		wg.Add(1)
		go func() {
			defer func() {
				wg.Done()
			}()
			for j := 0; j < NumOpertions; j++ {
				for {
					if _, ok := q.Pop(); ok {
						// mu.Lock()
						// set[v.(int32)] = true
						// mu.Unlock()
						break
					}
				}

				// for {
				// 	var v interface{}
				// 	mu.Lock()
				// 	if q.Len() != 0 {
				// 		v = q.Remove(q.Front())
				// 	}
				// 	mu.Unlock()
				// 	if v != nil {
				// 		break
				// 	}
				// }
			}
		}()
	}
	wg.Wait()
	endTime := time.Now().UnixMilli()
	fmt.Println("push: ", float64(endTime-starTime))
	println(a, b)

	// starTime = time.Now().UnixMilli()
	// for i := 0; i < NumGroutines; i++ {
	// 	wg.Add(1)
	// 	go func() {
	// 		defer func() {
	// 			wg.Done()
	// 		}()
	// 		for j := 0; j < NumOpertions; j++ {
	// 			if _, ok := q.Pop(); !ok {
	// 				println("XXXXX")
	// 			}

	// 			// mu.Lock()
	// 			// q.Remove(q.Front())
	// 			// mu.Unlock()
	// 		}
	// 	}()
	// }
	// wg.Wait()
	// endTime = time.Now().UnixMilli()
	// fmt.Println("pop: ", float64(endTime-starTime))

	fmt.Println(x, y, z)
	fmt.Println(q.size)
}
