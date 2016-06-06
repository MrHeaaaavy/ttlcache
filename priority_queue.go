package ttlcache

import (
	"container/heap"
	"sync"
)

func newPriorityQueue() *priorityQueue {
	queue := &priorityQueue{}
	heap.Init(queue)
	return queue
}

type priorityQueue struct {
	mutex sync.Mutex
	items []*Item
}

func (pq *priorityQueue) update(item *Item) {
	heap.Fix(pq, item.QueueIndex)
}

func (pq *priorityQueue) push(item *Item) {
	heap.Push(pq, item)
}

func (pq *priorityQueue) pop() *Item {
	if pq.Len() == 0 {
		return nil
	}
	return heap.Pop(pq).(*Item)
}

func (pq *priorityQueue) remove(item *Item) {
	heap.Remove(pq, item.QueueIndex)
}

func (pq priorityQueue) Len() int {
	pq.mutex.Lock()
	length := len(pq.items)
	pq.mutex.Unlock()
	return length
}

func (pq priorityQueue) Less(i, j int) bool {
	pq.mutex.Lock()
	less := pq.items[i].ExpireAt.Before(pq.items[j].ExpireAt)
	pq.mutex.Unlock()
	return less
}

func (pq priorityQueue) Swap(i, j int) {
	pq.mutex.Lock()
	pq.items[i], pq.items[j] = pq.items[j], pq.items[i]
	pq.items[i].QueueIndex = i
	pq.items[j].QueueIndex = j
	pq.mutex.Unlock()
}

func (pq *priorityQueue) Push(x interface{}) {
	pq.mutex.Lock()
	item := x.(*Item)
	item.QueueIndex = len(pq.items)
	pq.items = append(pq.items, item)
	pq.mutex.Unlock()
}

func (pq *priorityQueue) Pop() interface{} {
	pq.mutex.Lock()
	old := pq.items
	n := len(old)
	item := old[n-1]
	item.QueueIndex = -1
	pq.items = old[0 : n-1]
	pq.mutex.Unlock()
	return item
}
