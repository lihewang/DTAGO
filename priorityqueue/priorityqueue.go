package priorityqueue

import (
	"container/heap"
	tp "typedef"
)

//PriorityQueue pq
type PriorityQueue []*tp.PQItem

//Len pq length
func (pq PriorityQueue) Len() int { return len(pq) }

//Less min
func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].IMP < pq[j].IMP
}

//Swap swap
func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].PqIndex = i
	pq[j].PqIndex = j
}

//Push insert
func (pq *PriorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*tp.PQItem)
	item.PqIndex = n
	*pq = append(*pq, item)
}

//Pop pop
func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	item.PqIndex = -1 // for safety
	*pq = old[0 : n-1]
	return item
}

//Update decrease key
func (pq *PriorityQueue) Update(item *tp.PQItem) {
	heap.Fix(pq, item.PqIndex)
}
