package readfile

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	tp "typedef"
)

//ReadPar Read Parameters
func ReadPar(par *tp.Par) {
	CtrlFile := "eltod.ctl"
	file, err := os.Open(CtrlFile)
	if err != nil {
		log.Fatal(err)
	}
	parTemp := make(map[string]string)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() { // internally, it advances token based on sperator
		key := strings.TrimSpace(scanner.Text()[:40])
		value := strings.TrimSpace(scanner.Text()[40:])
		parTemp[key] = value
		//fmt.Println(key + ".." + parTemp[key])
	}
	par.Link = parTemp["LINK_FILE"]
	par.Node = parTemp["NODE_FILE"]
	par.TripTable = parTemp["TRIP_FILE"]
	par.TollConstFile = parTemp["TOLL_CONSTANT_FILE"]
	par.ChoiceLogFile = parTemp["NEW_MODEL_DATA_FILE"]
	par.NodeTypeZone = ParsetoInt(parTemp["ZONE_NODE_TYPE"])
	par.NodeTypeEntry = ParsetoArrayInt(parTemp["EXPRESS_ENTRY_TYPES"])
	par.NodeTypeExit = ParsetoArrayInt(parTemp["EXPRESS_EXIT_TYPES"])
	par.LinkTypeEL = ParsetoArrayInt(parTemp["EXPRESS_FACILITY_TYPES"])
	par.ELEntry = ParsetoInt(parTemp["EXPRESS_LINK_ENTRY_TYPE"])
	par.ELExit = ParsetoInt(parTemp["EXPRESS_LINK_EXIT_TYPE"])
	par.NumThreads = ParsetoInt(parTemp["NUMBER_OF_THREADS"])
	par.MaxIter = ParsetoInt(parTemp["MAXIMUM_ITERATIONS"])
	par.VOT = ParsetoArrayFloat(parTemp["COST_VALUE"])
	par.MinSplit = ParsetoFloat(parTemp["MINIMUM_TRIP_SPLIT"])

	//choice model par
	par.CMTime = ParsetoArrayFloat(parTemp["MODEL_TIME_FACTOR"])
	par.CMToll = ParsetoArrayFloat(parTemp["MODEL_TOLL_FACTOR"])
	par.CMReliaRatio = ParsetoFloat(parTemp["MODEL_RELIABILITY_RATIO"])
	par.CMReliaTime = ParsetoFloat(parTemp["MODEL_RELIABILITY_TIME"])
	par.CMReliaDist = ParsetoFloat(parTemp["MODEL_RELIABILITY_DISTANCE"])
	par.CMPercTime = ParsetoFloat(parTemp["MODEL_PERCEIVED_TIME"])
	par.CMPercMidVC = ParsetoFloat(parTemp["MODEL_PERCEIVED_MID_VC"])
	par.CMPercMaxVC = ParsetoFloat(parTemp["MODEL_PERCEIVED_MAX_VC"])
	par.CMELWeight = ParsetoFloat(parTemp["MODEL_EXPRESS_WEIGHT"])
	par.CMScaleLen = ParsetoFloat(parTemp["MODEL_SCALE_LENGTH"])
	par.CMMaxCir = ParsetoFloat(parTemp["MODEL_MAX_CIRCUITY"])
	par.CLogTS = ParsetoArrayInt(parTemp["SELECT_MODEL_PERIODS"])
	par.CLogIter = ParsetoArrayInt(parTemp["SELECT_MODEL_ITERATIONS"])
	par.CLogNode = ParsetoArrayInt(parTemp["SELECT_MODEL_NODES"])
	par.CLogVCLS = ParseFields(parTemp["SELECT_MODEL_MODES"])

	//parse toll policy
	tpCodes := ParsetoArrayInt(parTemp["TOLL_POLICY_CODES"])
	tpMin := ParsetoArrayFloat(parTemp["MINIMUM_TOLL"])
	tpMax := ParsetoArrayFloat(parTemp["MAXIMUM_TOLL"])
	tpVCOffset := ParsetoArrayFloat(parTemp["VC_RATIO_OFFSET"])
	tpExp := ParsetoArrayFloat(parTemp["TOLL_EXPONENT"])
	par.TollPolicy = make(map[int]map[string]float64)
	for i := range tpCodes {
		plc := make(map[string]float64)
		plc["MinToll"] = tpMin[i]
		plc["MaxToll"] = tpMax[i]
		plc["Offset"] = tpVCOffset[i]
		plc["Exp"] = tpExp[i]
		par.TollPolicy[i+1] = plc
	}

	file.Close()

	//read toll constant
	file, err = os.Open(par.TollConstFile)
	if err != nil {
		log.Fatal(err)
	}
	scanner = bufio.NewScanner(file)
	i := 0
	var NodeID float64
	var cst []float64
	cstpar := make(map[int][]float64)
	NodeID = -1
	for scanner.Scan() { // internally, it advances token based on sperator
		if i > 0 {
			txt := strings.Split(scanner.Text(), ",")
			n, _ := strconv.ParseFloat(strings.TrimSpace(txt[0]), 64)
			value, _ := strconv.ParseFloat(strings.TrimSpace(txt[2]), 64)
			if NodeID != n {
				if i > 1 {
					cstpar[int(NodeID)] = cst
					cst = nil
				}
				NodeID = n
			}
			cst = append(cst, value)
		}
		i++
	}
	cstpar[int(NodeID)] = cst
	par.TollConst = cstpar
}

//ParsetoInt parse parameters
func ParsetoInt(s string) (par int) {
	p, _ := strconv.ParseInt(s, 10, 0)
	return int(p)
}

//ParsetoFloat parse parameters
func ParsetoFloat(s string) (par float64) {
	par, _ = strconv.ParseFloat(s, 64)
	return
}

//ParsetoArrayInt parse parameters
func ParsetoArrayInt(s string) (par []int) {
	for _, element := range strings.Split(s, ",") {
		p, _ := strconv.ParseInt(element, 10, 0)
		par = append(par, int(p))
	}
	return
}

//ReadNode read node file
func ReadNode(par *tp.Par, nodes map[int]*tp.Node, wr io.Writer) {
	file, err := os.Open(par.Node)
	//get zone type
	ztp := par.NodeTypeZone
	if err != nil {
		log.Fatal(err)
	}
	var i int32
	var numzones int
	i = 0
	numzones = 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if i > 0 {
			txt := strings.Split(scanner.Text(), ",")
			n, err := strconv.ParseInt(strings.TrimSpace(txt[0]), 10, 32)
			x := ParsetoFloat(strings.TrimSpace(txt[1]))
			y := ParsetoFloat(strings.TrimSpace(txt[2]))
			ndtype, err := strconv.ParseInt(strings.TrimSpace(txt[3]), 10, 8)
			dngrp, err := strconv.ParseInt(strings.TrimSpace(txt[4]), 10, 8)
			if err != nil {
				log.Fatal(err)
			}
			elentry := false
			for _, v := range par.NodeTypeEntry {
				if int(ndtype) == v {
					elentry = true
					break
				}
			}
			elexit := false
			for _, v := range par.NodeTypeExit {
				if int(ndtype) == v {
					elexit = true
					break
				}
			}
			node := tp.Node{
				N:         int(n),
				X:         x,
				Y:         y,
				TYPE:      int(ndtype),
				DNGRP:     int(dngrp),
				IsELEntry: elentry,
				IsELExit:  elexit,
			}
			if node.TYPE == ztp {
				if numzones < node.N {
					numzones = node.N
				}
			}
			nodes[int(n)] = &node
		}
		i++
	}
	par.NumZones = numzones //number of zones
	file.Close()
	fmt.Fprintf(wr, "Read node file. Total of %+v zones.\n", numzones)
}

//ReadLink read link file
func ReadLink(par *tp.Par, links map[int]*tp.Link, nodes map[int]*tp.Node, wr io.Writer) {
	file, err := os.Open(par.Link)
	if err != nil {
		log.Fatal(err)
	}
	var i int
	i = 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if i > 0 {
			txt := strings.Split(scanner.Text(), ",")
			a, err := strconv.ParseInt(strings.TrimSpace(txt[0]), 10, 0)
			b, err := strconv.ParseInt(strings.TrimSpace(txt[1]), 10, 0)
			dist, err := strconv.ParseFloat(strings.TrimSpace(txt[2]), 64)
			capacity, err := strconv.ParseInt(strings.TrimSpace(txt[3]), 10, 0)
			lanes, err := strconv.ParseInt(strings.TrimSpace(txt[4]), 10, 0)
			ftype, err := strconv.ParseInt(strings.TrimSpace(txt[5]), 10, 0)
			ffspeed, err := strconv.ParseFloat(strings.TrimSpace(txt[6]), 64)
			alpha, err := strconv.ParseFloat(strings.TrimSpace(txt[7]), 64)
			beta, err := strconv.ParseFloat(strings.TrimSpace(txt[8]), 64)
			timefactor, err := strconv.ParseFloat(strings.TrimSpace(txt[9]), 64)
			tollpolicy, err := strconv.ParseInt(strings.TrimSpace(txt[11]), 10, 0)
			toll, err := strconv.ParseFloat(strings.TrimSpace(txt[12]), 64)
			tolldist, err := strconv.ParseFloat(strings.TrimSpace(txt[13]), 64)
			if err != nil {
				log.Fatal(err)
			}
			//EL min toll. if TOLL attribute > 0, use TOLL; else use toll policy
			if tlp, ok := par.TollPolicy[int(tollpolicy)]; ok {
				if toll == 0 {
					toll = tlp["MinToll"]
				}
			}
			//CgTime
			var cgtimes [96]float64
			var tolls [96]float64
			var weight [96]float64
			for j := 0; j <= 95; j++ {
				cgtimes[j] = dist / ffspeed * 60
				tolls[j] = toll
				weight[j] = 1
			}
			//EL
			elentry := false
			if int(ftype) == par.ELEntry {
				elentry = true
			}
			elexit := false
			if int(ftype) == par.ELExit {
				elexit = true
			}

			link := tp.Link{
				ID:          i,
				A:           int(a),
				B:           int(b),
				DIST:        dist,
				CAPACITY:    int(capacity),
				NUMLANES:    int(lanes),
				FTYPE:       int(ftype),
				FFSPEED:     ffspeed,
				ALPHA:       alpha,
				BETA:        beta,
				TIMEFACTOR:  timefactor,
				TOLLPOLICY:  int(tollpolicy),
				TOLL:        toll,
				TOLLSEGDIST: tolldist,
				NODEA:       nodes[int(a)],
				NODEB:       nodes[int(b)],
				CgTime:      cgtimes,
				FFTime:      dist / ffspeed * 60,
				TimeWeight:  weight,
				TollRate:    tolls,
				IsELEntry:   elentry,
				IsELExit:    elexit,
			}

			links[i] = &link
			nodes[int(a)].DNLINKS = append(nodes[int(a)].DNLINKS, &link) //add down links to node A
		}
		i = i + 1
	}
	file.Close()
	fmt.Fprintf(wr, "Read link file. Total of %+v links.\n", i-1)
}

//CheckNetwork check network coding
func CheckNetwork(par *tp.Par, links map[int]*tp.Link, nodes map[int]*tp.Node, wr io.Writer) {
	//decision nodes need at least two downstream links
	fatal := 0
	numDNd := 0
	for _, nd := range nodes {
		if nd.TYPE == par.NodeTypeEntry[0] || nd.TYPE == par.NodeTypeExit[0] {
			HasRamp := false
			numDNd++
			if len(nd.DNLINKS) < 2 {
				fmt.Fprintf(wr, "!Fatal network coding error: decision node %v has no alternative path.\n", nd.N)
				fatal++
			}
			for _, dlnk := range nd.DNLINKS {
				if dlnk.IsELEntry || dlnk.IsELExit {
					HasRamp = true
				}
			}
			if !HasRamp {
				fmt.Fprintf(wr, "!Fatal network coding error: decision node %v has no express path.\n", nd.N)
				fatal++
			}
		}
	}
	fmt.Fprintf(wr, "Number of decision node %v.\n", numDNd)
	if fatal > 0 {
		os.Exit(1)
	}
}

//ReadTT Read trip table
func ReadTT(par *tp.Par, nodes map[int]*tp.Node, TaskQueue chan *tp.PathTask, wr io.Writer) {
	file, err := os.Open(par.TripTable)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	var i int
	var ts, iNodeID float64 //time step
	jNodes := []*tp.Node{}
	tripRecs := [][]float64{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if i == 0 { //first record, field names
			line := ParseFields(scanner.Text())
			//Put pars into maps
			par.VOTFactor = make(map[string]float64)
			par.CMTimeFactor = make(map[string]float64)
			par.CMTollFactor = make(map[string]float64)
			for j, element := range line[3:] {
				par.VehClass = append(par.VehClass, element)
				par.VOTFactor[element] = par.VOT[j]
				par.CMTimeFactor[element] = par.CMTime[j]
				par.CMTollFactor[element] = par.CMToll[j]
			}
		} else { //data
			Record := ParsetoArrayFloat(scanner.Text()) //parsed to float64
			tsCurr := Record[2]
			iNodeIDCurr := Record[0]
			if i > 1 && (ts != tsCurr || iNodeID != iNodeIDCurr) { //diffrent origin node or time step
				//create task for each vehicle class
				for col, element := range par.VehClass {
					trips := make(map[*tp.Node]float64)
					for j := 0; j < len(tripRecs); j++ {
						if tripRecs[j][col] > 0 {
							trips[jNodes[j]] = tripRecs[j][col]
						}
					}
					if len(trips) > 0 {
						task := tp.PathTask{
							TP:   int(ts),
							I:    nodes[int(iNodeID)],
							VCLS: element,
							TRIP: trips,
							Rcd: i,
						}
						TaskQueue <- &task //send task to task queue
						//fmt.Printf("task sent to queue %+v\n", task)
					}
				}
				jNodes = nil
				tripRecs = nil
			}
			ts = tsCurr
			iNodeID = iNodeIDCurr
			jNodeID := Record[1]
			dNd := nodes[int(jNodeID)]
			jNodes = append(jNodes, dNd)
			tripsList := []float64{}
			for j := range par.VehClass {
				tripsList = append(tripsList, Record[j+3])
			}
			tripRecs = append(tripRecs, tripsList)
			//fmt.Printf("record %+v\n", Record)
		}
		i++
	}
	close(TaskQueue)
}

//ParsetoArrayFloat parse parameters
func ParsetoArrayFloat(s string) []float64 {
	result := []float64{}
	for _, element := range strings.Split(s, ",") {
		p, _ := strconv.ParseFloat(strings.TrimSpace(element), 64)
		result = append(result, p)
	}
	return result
}

//ParseFields parse parameters
func ParseFields(s string) (par []string) {
	par = strings.Split(s, ",")
	return
}
