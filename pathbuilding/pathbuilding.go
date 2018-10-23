package pathbuilding

import (
	"container/heap"
	"encoding/csv"

	//"fmt"
	"io"
	"math"
	p "priorityqueue"
	tp "typedef"
)

//SPath build shortest path
func SPath(par *tp.Par, task *tp.PathTask, cLog io.Writer) map[*tp.Node][]*tp.PQItem {
	//count paths build
	if task.BanLink == nil {
		par.PathBuildCount[0]++
	} else {
		par.PathBuildCount[1]++
	}
	numJ := len(task.TRIP)
	pqItemAll := make(map[*tp.Node]*tp.PQItem)
	newpqItem := tp.PQItem{Nd: task.I, TP: task.TP, Visited: false}
	pqItemAll[task.I] = &newpqItem
	pq := make(p.PriorityQueue, 0)
	heap.Init(&pq)
	CurrpqItem := &newpqItem
	//loop until all zone nodes have been visited
	for {
		CurrpqItem.Visited = true
		CurrNd := CurrpqItem.Nd
		if _, ok := task.TRIP[CurrNd]; ok { //check if all destination nodes have been visited
			numJ--
			if numJ == 0 {
				break
			}
		}
		if task.DisttoDest > 0 && CurrpqItem.Dist > task.DisttoDest { //current path has exceed the max length
			return nil
		}
		if !(CurrNd.TYPE == par.NodeTypeZone && CurrNd != task.I) { //skip zone node
			for _, link := range CurrNd.DNLINKS {
				if link == task.BanLink { //for decision node path
					continue
				}
				var pqNextNode *tp.PQItem
				if nd, ok := pqItemAll[link.NODEB]; ok { //checked before
					if nd.Visited {
						continue
					} else {
						pqNextNode = nd
					}
				} else { //not checked before
					NewNd := tp.PQItem{Nd: link.NODEB, ParentItem: nil, IMP: 999999.9, TP: task.TP, Visited: false}
					pqItemAll[link.NODEB] = &NewNd
					heap.Push(&pq, &NewNd)
					pqNextNode = &NewNd
				}
				//update impedence and skims
				var htime float64
				if len(task.TRIP) == 1 {
					var jNd *tp.Node
					NextNd := pqNextNode.Nd
					for v := range task.TRIP {
						jNd = v
					}
					htime = math.Sqrt((NextNd.X-jNd.X)*(NextNd.X-jNd.X)+(NextNd.Y-jNd.Y)*(NextNd.Y-jNd.Y)) / 5280
				}
				timestep := CurrpqItem.TP - 1
				var imptoll float64
				imptoll = link.TollRate[timestep] * par.VOTFactor[task.VCLS]
				if task.DisttoDest > 0 { //alternative path building
					imptoll = 0
				}
				impedent := CurrpqItem.IMP + link.CgTime[timestep]*link.TimeWeight[timestep] + imptoll + htime
				if impedent < pqNextNode.IMP {
					pqNextNode.Time = CurrpqItem.Time + link.CgTime[timestep]
					timepd := int(pqNextNode.Time/15) + task.TP
					if timepd > 96 {
						timepd = timepd - 96
					}
					pqNextNode.TP = timepd
					pqNextNode.IMP = impedent
					pqNextNode.Dist = CurrpqItem.Dist + link.DIST
					pqNextNode.Toll = CurrpqItem.Toll + link.TollRate[timestep]
					pqNextNode.FFTime = CurrpqItem.FFTime + link.FFTime
					pqNextNode.ParentItem = CurrpqItem
					pqNextNode.Link = link					
				}
				pq.Update(pqNextNode)
			}
		} 
		if len(pq) > 0 {
			CurrpqItem = heap.Pop(&pq).(*tp.PQItem)
		} else {
			return nil
		}
	}

	//Get paths
	//var path map[*tp.Node][]*tp.PQItem
	path := make(map[*tp.Node][]*tp.PQItem)
	for nd := range task.TRIP {
		pqI := pqItemAll[nd]
		var nodelist []*tp.PQItem
		nodelist = append(nodelist, pqI)
		//var pathWrite []int
		//pathWrite = append(pathWrite, pqI.Nd.N)
		for pqI.ParentItem != nil {
			pqI = pqI.ParentItem
			nodelist = append(nodelist, pqI)
			//pathWrite = append(pathWrite, pqI.Nd.N)
		}
		//if task.I.N == 51{
			//fmt.Printf("path %v\n", pathWrite)		
		//}
		path[nd] = nodelist
	}
	//fmt.Printf("Number of nodes searched = %v\n", len(pqItemAll))	
	//fmt.Printf("----------------\n")
	return path
}

//LoadLink load links
func LoadLink(par *tp.Par, task *tp.PathTask, path map[*tp.Node][]*tp.PQItem, cLog io.Writer) {
	for nd, loadTrip := range task.TRIP {
		if path[nd] != nil {
			csvWriter := csv.NewWriter(cLog)
			//load link
			if loadTrip <= par.MinSplit || task.VCLS == "TRK" { //trip can't be split
				for i := len(path[nd]) - 2; i >= 0; i-- {
					item := (path[nd])[i]
					item.Link.Add(item.TP, loadTrip, task.VCLS)
				}
			} else { //trip is large enough to be split
				for i := len(path[nd]) - 2; i >= 0; i-- {
					item := (path[nd])[i]
					item.Link.Add(item.TP, loadTrip, task.VCLS)
					if item.Nd.IsELEntry || item.Nd.IsELExit { //decision node
						var tltype int
						var dists [2]float64
						//get current path type
						if item.Nd.IsELEntry {
							if (path[nd])[i-1].Link.IsELEntry {
								tltype = 0 //current path is toll path
							} else {
								tltype = 1
							}
						} else if item.Nd.IsELExit {
							if (path[nd])[i-1].Link.IsELExit {
								tltype = 1 //current path is non-toll path
							} else {
								tltype = 0
							}
						}
						lastItem := (path[nd])[0]
						dists[tltype] = lastItem.Dist - item.Dist
						//alternative path task
						var newTrips [2]float64
						newtrip := make(map[*tp.Node]float64)
						newtrip[nd] = loadTrip
						newtask := tp.PathTask{
							TP:         item.TP,
							I:          item.Nd,
							O:          task.O,
							VCLS:       task.VCLS,
							TRIP:       newtrip,
							BanLink:    (path[nd])[i-1].Link,
							DisttoDest: dists[tltype] * par.CMMaxCir,
						}
						newPath := SPath(par, &newtask, cLog) //build alternative path
						var newtltype int
						if tltype == 0 {
							newtltype = 1
						} else {
							newtltype = 0
						}
						var newItem *tp.PQItem
						if newPath == nil {
							dists[newtltype] = 9999
						} else {
							newPath := path[nd]
							newItem = newPath[0]
							dists[newtltype] = newItem.Dist
						}
						//check distance
						if dists[newtltype]/dists[tltype] > par.CMMaxCir {
							//no split, continue current path
						} else if dists[tltype]/dists[newtltype] > par.CMMaxCir {
							LoadLink(par, &newtask, newPath, cLog) //no split, use new path
						} else { //split
							var times [2]float64
							var pcvtimes [2]float64
							var fftimes [2]float64
							var tolls [2]float64
							var reliabilities [2]float64
							//current path skims
							times[tltype] = lastItem.Time - item.Time
							pcvtimes[tltype] = lastItem.PcvTime - item.PcvTime
							fftimes[tltype] = lastItem.FFTime - item.FFTime
							tolls[tltype] = lastItem.Toll - item.Toll
							reliabilities[tltype] = par.CMReliaTime * (times[tltype] - fftimes[tltype]) *
								math.Pow(dists[tltype], -1*par.CMReliaDist)
							//alternative path skims
							times[newtltype] = newItem.Time
							pcvtimes[newtltype] = newItem.PcvTime
							fftimes[newtltype] = newItem.FFTime
							tolls[newtltype] = newItem.Toll
							reliabilities[newtltype] = par.CMReliaTime * (newItem.Time - newItem.FFTime) *
								math.Pow(newItem.Dist, -1*par.CMReliaDist)

							//get toll constant
							var tollconst float64
							if value, ok := par.TollConst[item.Nd.N]; ok {
								tollconst = value[item.TP-1]
							} else {
								tollconst = par.TollConst[0][item.TP-1]
							}

							//calculate toll share
							utility := tollconst + par.CMTimeFactor[task.VCLS]*(times[0]-times[1]) +
								par.CMTollFactor[task.VCLS]*(tolls[0]-tolls[1]) +
								par.CMReliaRatio*par.CMTimeFactor[task.VCLS]*(reliabilities[0]-reliabilities[1])
							tollshare := 1 / (1 + math.Exp(-1*utility))

							//Check against pervious choices
							if tollshare <= task.LShare { //do not split, use toll path
								if tltype == 0 { //current toll path is toll path
									continue
								} else {
									newtask.TRIP[nd] = loadTrip
									LoadLink(par, &newtask, newPath, cLog) //load new path
									break
								}
							} else if tollshare >= task.HShare { //do not split, use non-toll path
								if tltype == 0 { //current path is toll path, switch to non-toll path
									newtask.TRIP[nd] = loadTrip
									LoadLink(par, &newtask, newPath, cLog) //load new path
									break
								} else {
									continue
								}
							} else { //need to split
								newTrips[0] = task.TRIP[nd] * tollshare   //toll trips
								newTrips[1] = task.TRIP[nd] - newTrips[0] //non-toll
								if tltype == 0 {                          //toll path
									loadTrip = newTrips[0] //current path trips
									newtask.TRIP[nd] = newTrips[1]
									newtask.LShare = tollshare
								} else { //non-toll path
									loadTrip = newTrips[1]
									newtask.TRIP[nd] = newTrips[0]
									newtask.HShare = tollshare
								}
								newtask.NumSplit = task.NumSplit + 1
								LoadLink(par, &newtask, newPath, cLog) //load new path
							}
							//log
							/*if FindinArray(par.Iter, par.CLogIter) && FindinArray(task.TP, par.CLogTS) &&
								FindinArray(item.Nd.N, par.CLogNode) {
								record := []string{strconv.Itoa(par.Iter), strconv.Itoa(task.O.N), strconv.Itoa(task.DestNode.N),
									strconv.Itoa(task.TP), strconv.Itoa(task.TP), task.VCLS, strconv.Itoa(item.Nd.N), strconv.Itoa(task.DestNode.N),
									strconv.FormatFloat(dists[0], 'f', 2, 64), strconv.FormatFloat(dists[1], 'f', 2, 64),
									strconv.FormatFloat(times[0], 'f', 2, 64), strconv.FormatFloat(times[1], 'f', 2, 64),
									strconv.FormatFloat(pcvtimes[0], 'f', 2, 64), strconv.FormatFloat(pcvtimes[1], 'f', 2, 64),
									strconv.FormatFloat(fftimes[0], 'f', 2, 64), strconv.FormatFloat(fftimes[1], 'f', 2, 64),
									strconv.FormatFloat(tolls[0], 'f', 2, 64), strconv.FormatFloat(tolls[1], 'f', 2, 64),
									strconv.FormatFloat(reliabilities[0], 'f', 2, 64), strconv.FormatFloat(reliabilities[1], 'f', 2, 64),
									strconv.FormatFloat(utility, 'f', 2, 64), strconv.FormatFloat(tollshare, 'f', 2, 64),
									strconv.FormatFloat(newTrips[0], 'f', 2, 64), strconv.FormatFloat(newTrips[1], 'f', 2, 64)}
								csvWriter.Write(record)
							}*/
						}
					}
				}
			}
			
			var v int
			if task.NumSplit > 10 {
				v = 10
			} else {
				v = task.NumSplit
			}
			par.SplitCount[v]++
			csvWriter.Flush()
		}
	}
}

//FindinArray find number in array
func FindinArray(v int, a []int) bool {
	result := false
	for _, element := range a {
		if v == element {
			result = true
			break
		}
	}
	return result
}

