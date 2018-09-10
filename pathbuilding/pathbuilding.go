package pathbuilding

import (
	"container/heap"
	"fmt"
	p "priorityqueue"
	tp "typedef"
)

//SPath build shortest path
func SPath(par *tp.Par, task *tp.PathTask) (paths map[*tp.Node][]*tp.PathLink) {
	numJ := len(task.TRIP)
	zoneVisited := 1

	pq := make(p.PriorityQueue, 1)
	CurrNd := task.I
	pqItem := tp.PQItem{Nd: CurrNd, ParentItem: nil, IMP: 0, Index: 0, Visited: true}
	pq[0] = &pqItem
	heap.Init(&pq)
	pqItemAll := make(map[*tp.Node]*tp.PQItem)
	for zoneVisited < int(par.NumZones) {
		CurrpqItem := heap.Pop(&pq).(*tp.PQItem)
		CurrpqItem.Visited = true
		CurrNd = CurrpqItem.Nd
		if _, ok := task.TRIP[CurrNd]; ok { //check if all destination nodes have been visited
			numJ--
			if numJ == 0 {
				break
			}
		}
		if CurrNd.TYPE == par.NodeTypeZone[0] { //skip zone node
			if zoneVisited > 1 {
				continue
			} else {
				zoneVisited++
			}
		}
		for _, link := range CurrNd.DNLINKS {
			if element, ok := pqItemAll[link.NODEB]; ok { //node was checked
				if element.Visited {
					continue
				} else {
					impedent := pqItem.IMP + link.CgTime[pqItem.TP]*link.TimeWeight + link.TollRate[pqItem.TP]*par.VOTFactor[task.VCLS]
					if impedent < element.IMP {
						element.Time = pqItem.Time + link.CgTime[pqItem.TP]
						time := pqItem.Time + link.CgTime[pqItem.TP]
						timepd := int(time/15) + task.TP
						if timepd > 96 {
							timepd = timepd - 96
						}
						element.TP = timepd
						element.IMP = impedent
						element.ParentItem = CurrpqItem
						element.Link = link
					}
				}
			} else { //node was not checked
				impedent := pqItem.IMP + link.CgTime[pqItem.TP]*link.TimeWeight + link.TollRate[pqItem.TP]*par.VOTFactor[task.VCLS]
				time := pqItem.Time + link.CgTime[pqItem.TP]
				timepd := int(time/15) + task.TP
				if timepd > 96 {
					timepd = timepd - 96
				}
				newpqItem := tp.PQItem{Nd: link.NODEB, Link: link, ParentItem: CurrpqItem, IMP: impedent, Time: time,
					TP: timepd, Index: 0, Visited: false}
				pqItemAll[link.NODEB] = &newpqItem
				heap.Push(&pq, &newpqItem)
			}
			pq.Update(&pqItem, pqItem.Nd, pqItem.IMP)
		}
	}

	//Get paths
	for nd := range task.TRIP {
		pqI := pqItemAll[nd]
		var p []int
		var pathLinks []*tp.PathLink
		p = append(p, pqI.Nd.N)
		pathLinks = append(pathLinks, &tp.PathLink{Link: pqI.Link, TP: pqI.TP, BNdType: 0})
		for pqI.ParentItem != nil {
			pqI = pqI.ParentItem
			p = append(p, pqI.Nd.N)
			var ndtp int
			if pqI.Nd.TYPE == par.NodeTypeEntry[0] {
				ndtp = 1
			} else if pqI.Nd.TYPE == par.NodeTypeExit[0] {
				ndtp = 2
			} else {
				ndtp = 0
			}
			pathLinks = append(pathLinks, &tp.PathLink{Link: pqI.Link, TP: pqI.TP, BNdType: ndtp})
		}
		paths[nd] = pathLinks

		fmt.Printf("path %v", p)
	}
	return
}

//LoadLink load links
func LoadLink(task *tp.PathTask, paths map[*tp.Node][]*tp.PathLink) {
	for i, path := range paths {
		sNd := paths[i]
		//load link
		for i := range path {
			path[i].Link.Vol[path[i].TP] = path[i].Link.Vol[path[i].TP] + task.TRIP[nd]
		}
	}
}
