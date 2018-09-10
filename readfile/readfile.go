package readfile

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	tp "typedef"
)

//ReadPar Read Parameters
func ReadPar(par *tp.Par) {
	CtrlFile := "C:/Users/lihe.wang/go/input/eltod.ctl"
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
		//fmt.Println(key + ".." + par[key])
	}
	par.Link = parTemp["LINK_FILE"]
	par.Node = parTemp["NODE_FILE"]
	par.TripTable = parTemp["TRIP_FILE"]
	par.NodeTypeZone = ParsetoArray(parTemp["ZONE_NODE_TYPE"])
	par.NodeTypeEntry = ParsetoArray(parTemp["EXPRESS_ENTRY_TYPES"])
	par.NodeTypeExit = ParsetoArray(parTemp["EXPRESS_EXIT_TYPES"])
	par.NumThreads = ParsetoInt(parTemp["NUMBER_OF_THREADS"])
	par.VOT = ParseRecord(parTemp["COST_VALUE"])

	//parse toll policy
	tpCodes := ParsetoArray(parTemp["TOLL_POLICY_CODES"])
	tpMin := ParseRecord(parTemp["MINIMUM_TOLL"])
	tpMax := ParseRecord(parTemp["MAXIMUM_TOLL"])
	tpVCOffset := ParseRecord(parTemp["VC_RATIO_OFFSET"])
	tpExp := ParseRecord(parTemp["TOLL_EXPONENT"])
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
	fmt.Println("Control file: " + CtrlFile)
}

//ParsetoInt parse parameters
func ParsetoInt(s string) (par int) {
	p, _ := strconv.ParseInt(s, 10, 0)
	return int(p)
}

//ParsetoArray parse parameters
func ParsetoArray(s string) (par []int) {
	for _, element := range strings.Split(s, ",") {
		p, _ := strconv.ParseInt(element, 10, 0)
		par = append(par, int(p))
	}
	return
}

//ReadNode read node file
func ReadNode(par *tp.Par, nodes map[int]*tp.Node) {
	file, err := os.Open(par.Node)
	//get zone type
	ztp := par.NodeTypeZone[0]
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
			ndtype, err := strconv.ParseInt(strings.TrimSpace(txt[3]), 10, 8)
			dngrp, err := strconv.ParseInt(strings.TrimSpace(txt[4]), 10, 8)
			if err != nil {
				log.Fatal(err)
			}

			node := tp.Node{
				N:     int(n),
				TYPE:  int(ndtype),
				DNGRP: int(dngrp),
			}
			if node.TYPE == ztp {
				if numzones < node.N {
					numzones = node.N
				}
			}
			nodes[int(n)] = &node
		}
		i = i + 1
	}
	par.NumZones = numzones //number of zones
	file.Close()
	fmt.Printf("Read node file. Total of %+v zones.\n", numzones)
}

//ReadLink read link file
func ReadLink(par *tp.Par, links map[int]*tp.Link, nodes map[int]*tp.Node) {
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
			var cgtimes []float64
			var tolls []float64
			for j := 1; j <= 96; j++ {
				cgtimes = append(cgtimes, dist/ffspeed*60)
				tolls = append(tolls, toll)
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
				TimeWeight:  1,
				TollRate:    tolls,
			}

			links[i] = &link
			nodes[int(a)].DNLINKS = append(nodes[int(a)].DNLINKS, &link) //add down links to node A
		}
		i = i + 1
	}
	file.Close()
	fmt.Printf("Read link file. Total of %+v links.\n", i-1)
}

//ReadTT Read trip table
func ReadTT(par *tp.Par, nodes map[int]*tp.Node, TaskQueue chan *tp.PathTask) {
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
			Record := ParseFields(scanner.Text())
			for j, element := range Record[3:] {
				par.VehClass = append(par.VehClass, element)
				//VOT factor by veh class
				par.VOTFactor = make(map[string]float64)
				par.VOTFactor[element] = par.VOT[j]
			}
			fmt.Printf("Read trip table file. Vehicle classes %+v\n", par.VehClass)
		} else { //data
			Record := ParseRecord(scanner.Text()) //parsed to float64
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
		i = i + 1
		if i > 20 {
			close(TaskQueue)
			break
		}
	}
}

//ParseRecord parse parameters
func ParseRecord(s string) []float64 {
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
