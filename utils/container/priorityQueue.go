// This example demonstrates a Priority queue built using the heap interface.
package container

import (
	"container/heap"
)

//https://pkg.go.dev/container/heap#example__priorityQueue
//https://programmer.help/blogs/simple-priority-queue-implemented-by-golang.html
//https://golangtc.com/t/56447ce3b09ecc72c3000056
//尽管golang不支持方法重载，且接口不具备属性。
//但是我们可以定义一个包含getter，setter的接口，依赖此接口的函数调用对应方法即可。
//https://stackoverflow.com/questions/42702117/interfaces-with-getters-and-setters-in-go-language

// An Item is something we manage in a Priority queue.
type Item struct {
	Value    string // The Value of the item; arbitrary.
	Priority int    // The Priority of the item in the queue.
	// The Index is needed by update and is maintained by the heap.Interface methods.
	Index int // The Index of the item in the heap.
}

// A PriorityQueue implements heap.Interface and holds Items.
type PriorityQueue []*Item

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	// We want Pop to give us the highest, not lowest, Priority so we use greater than here.
	return pq[i].Priority < pq[j].Priority
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].Index = i
	pq[j].Index = j
}

func (pq *PriorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*Item)
	item.Index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // avoid memory leak
	item.Index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}

// update modifies the Priority and Value of an Item in the queue.
func (pq *PriorityQueue) Update(item *Item, value string, priority int) {
	item.Value = value
	item.Priority = priority
	heap.Fix(pq, item.Index)
}
