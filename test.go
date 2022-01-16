
type (
	node struct {
		data int
		prev *node
		next *node
	}
)

type (
	LinkList struct {
		size int32
		head node
		tail node
	}
)

func NewLinkList() *LinkList {
	list := new(LinkList)
	list.head.next = &list.tail
	list.head.prev = &list.tail
	list.tail.prev = &list.head
	list.tail.next = &list.head
	return list
}

var (
	a, b int
)

func (list *LinkList) Push(data int) {
	newNode := &node{data: data}
	for {
		last := list.tail.next
		if last.prev == nil {
			if last.next.prev != last {
				runtime.Gosched()
				continue
			}
			atomic.CompareAndSwapPointer((*unsafe.Pointer)(unsafe.Pointer(&last.prev)), nil, unsafe.Pointer(&list.tail))
		}

		newNode.next = last
		if !atomic.CompareAndSwapPointer((*unsafe.Pointer)(unsafe.Pointer(&list.tail.next)), unsafe.Pointer(last), unsafe.Pointer(newNode)) {
			runtime.Gosched()
			continue
		}

		atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&last.prev)), unsafe.Pointer(newNode))
		atomic.CompareAndSwapPointer((*unsafe.Pointer)(unsafe.Pointer(&newNode.prev)), nil, unsafe.Pointer(&list.tail))
		atomic.AddInt32(&list.size, 1)

		return
	}
}


const (
	NumGroutines = 1000
	NumOpertions = 100000
)

func main() {
	var wg sync.WaitGroup
	var list = NewLinkList()
	// var list = &list.List{}
	// list = list.Init()
	// var mu sync.Mutex
	for i := 0; i < NumGroutines; i++ {
		wg.Add(1)
		go func() {
			defer func() {
				wg.Done()
			}()

			for j := 0; j < NumOpertions; j++ {
				list.Push(j)
			}
		}()
	}
	wg.Wait()
	// var count int
	// for curr := list.tail.next; curr != &list.head; curr = curr.next {
	// 	count++
	// }
	// fmt.Println(a, b, count)
}
