package typedef

//Par parameter
type Par struct {
	Link, Node, TripTable, LinkClose, TurnProhibit        string
	NodeTypeZone, NodeTypeEntry, NodeTypeExit, LinkTypeEL []int
	NumZones, NumThreads                                  int
	VehClass                                              []string
	VOT                                                   []float64
	VOTFactor                                             map[string]float64
	TollPolicy                                            map[int]map[string]float64
}

//Link network link
type Link struct {
	ID, A, B                                                  int
	DIST, FFSPEED, ALPHA, BETA, TOLL, TOLLSEGDIST, TIMEFACTOR float64
	CAPACITY                                                  int
	NUMLANES, FTYPE, TOLLPOLICY                               int
	NODEA, NODEB                                              *Node
	FFTime, TimeWeight                                        float64 //calculated attributes
	CgTime, Vol, CgSpeed, TollRate                            []float64
}

//Node network node
type Node struct {
	N       int
	TYPE    int
	DNGRP   int
	DNLINKS []*Link
}

//PathTask path building task for path builder
type PathTask struct {
	TP   int               //time steps
	I    *Node             //origin node
	VCLS string            //vehicle class
	TRIP map[*Node]float64 //list of destinations and trips
}

//PQItem priority queue node
type PQItem struct {
	Nd                            *Node
	Link                          *Link //parent node link
	ParentItem                    *PQItem
	IMP, Time, Dist, Toll, FFTime float64
	TP, Index                     int
	Visited                       bool
}

//PathLink path link
type PathLink struct {
	Link    *Link
	TP      int
	BNdType int
}
