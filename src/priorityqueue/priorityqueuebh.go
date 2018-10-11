package priorityqueue

import (tp "typedef")

//Pqheap heap
type Pqheap struct{
	nds []*tp.PQItem
}

//MakeHeap init heap
func MakeHeap (x *tp.PQItem) *Pqheap{
	var h Pqheap
	h.nds = make([]*tp.PQItem, 2)
	h.nds[1] = x
	return &h
}

//Len length
func (h *Pqheap) Len() int {
	return len(h.nds) - 1 // heap is 1-indexed
}

//Insert insert node
func (h *Pqheap) Insert(x *tp.PQItem) {
	(*h).nds = append(h.nds, x)
	h.bubbleUp(len(h.nds) - 1)
}

func (h *Pqheap) bubbleUp(k int) {
	p, ok := parent(k)
	if !ok {
		return // k is root node
	}
	if h.nds[p].IMP > h.nds[k].IMP {
		h.nds[k], h.nds[p] = h.nds[p], h.nds[k]
		h.bubbleUp(p)
	}
}

func parent(k int) (int, bool) {
	if k == 1 {
		return 0, false
	}
	return k / 2, true
}

func left(k int) int {
	return 2 * k
}

func right(k int) int {
	return 2*k + 1
}

//Pop pop node
func (h *Pqheap) Pop() (*tp.PQItem, bool) {
	if h.Len() == 0 {
		return nil, false
	}
	v := h.nds[1]
	h.nds[1] = h.nds[h.Len()]
	(*h).nds = h.nds[:h.Len()]
	h.bubbleDown(1)
	return v, true
}

func (h *Pqheap) bubbleDown(k int) {
	min := k
	c := left(k)

	// find index of minimum value (k, k's left child, k's right child)
	for i := 0; i < 2; i++ {
		if (c + i) <= h.Len() {
			if h.nds[min].IMP > h.nds[c+i].IMP {
				min = c + i
			}
		}
	}
	if min != k {
		h.nds[k], h.nds[min] = h.nds[min], h.nds[k]
		h.bubbleDown(min)
	}
}
