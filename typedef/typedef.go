package typedef

import (
	"sync"
)

//Par parameter
type Par struct {
	Link, Node, TripTable, LinkClose, TurnProhibit, TollConstFile      string
	ChoiceLogFile                                                      string
	VehClass, CLogVCLS                                                 []string
	NodeTypeZone, NumZones, NumThreads, MaxIter, Iter, ELEntry, ELExit int
	NodeTypeEntry, NodeTypeExit, LinkTypeEL                            []int
	CLogTS, CLogIter, CLogNode                                         []int
	SplitCount                                                         [11]int
	PathBuildCount                                                     [2]int
	MinSplit, CMReliaRatio, CMReliaTime, CMReliaDist, CMPercTime       float64
	CMPercMidVC, CMPercMaxVC, CMELWeight, CMScaleLen, CMMaxCir         float64
	VOT, CMTime, CMToll                                                []float64
	VOTFactor, CMTimeFactor, CMTollFactor                              map[string]float64
	TollPolicy                                                         map[int]map[string]float64
	TollConst                                                          map[int][]float64
	SmoothForwardPercent, SmoothBackworkPercent, SmoothIter            int
}

//Link network link
type Link struct {
	ID, A, B                                                  int
	DIST, FFSPEED, ALPHA, BETA, TOLL, TOLLSEGDIST, TIMEFACTOR float64
	CAPACITY                                                  int
	NUMLANES, FTYPE, TOLLPOLICY                               int
	NODEA, NODEB                                              *Node
	FFTime                                                    float64 //calculated attributes
	CgTime, Vol, Volpre, CgSpeed, TollRate, TimeWeight        [96]float64
	VolCls, VolClsPre                                         map[string]*[96]float64
	IsELEntry, IsELExit                                       bool
}

//Add add volume
func (l *Link) Add(tp int, v float64, cls string) {
	var mu sync.RWMutex
	mu.RLock()
	defer mu.RUnlock()
	l.VolCls[cls][tp-1] = l.VolCls[cls][tp-1] + v
}

//Node network node
type Node struct {
	N                   int
	X, Y                float64
	TYPE                int
	DNGRP               int
	DNLINKS             []*Link
	IsELEntry, IsELExit bool
}

//PathTask path building task for path builder
type PathTask struct {
	TP, Rcd        int               //time steps
	I, O           *Node             //path origin node, packet start zone
	VCLS           string            //vehicle class
	TRIP           map[*Node]float64 //list of destinations and trips
	BanLink        *Link             //banned link for decision node path
	DisttoDest     float64           //distance to destination of the current path
	NumSplit       int               //number of split times
	LShare, HShare float64           //share boundary of packet
}

//PQItem priority queue node
type PQItem struct {
	Nd                                     *Node
	Link                                   *Link //parent node link
	ParentItem                             *PQItem
	IMP, Time, PcvTime, Dist, Toll, FFTime float64
	TP, PqIndex                            int
	Visited                                bool
}
