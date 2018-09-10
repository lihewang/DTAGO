package priorityqueue

import (
	"container/heap"
	tp "typedef"
)

//PriorityQueue priority queue
type PriorityQueue []*tp.PQItem

func (pq PriorityQueue) Len() int {
	return len(pq)
}
func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].IMP < pq[j].IMP
}

//Pop pop item
func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	item.Index = -1
	*pq = old[0 : n-1]
	return item
}

//Push push item
func (pq *PriorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*tp.PQItem)
	item.Index = n
	*pq = append(*pq, item)
}
func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].Index = i
	pq[j].Index = j
}

//Update update priority queue
func (pq *PriorityQueue) Update(item *tp.PQItem, value *tp.Node, priority float64) {
	//item.value = value
	item.IMP = priority
	heap.Fix(pq, item.Index)
}
