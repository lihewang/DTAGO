package main

import (
	"fmt"
	"io"
	"os"
	pbld "pathbuilding"
	readfile "readfile"
	"runtime"
	"sync"
	"time"
	tp "typedef"
	//"github.com/pkg/profile"
)

func main() {
	//command line: go tool pprof --pdf run.exe cpu.pprof > cpu_profile.pdf
	//defer profile.Start(profile.ProfilePath(".")).Stop()
	par := tp.Par{} //parameters
	runtime.GOMAXPROCS(runtime.NumCPU())
	logFile, _ := os.Create("log.txt")
	defer logFile.Close()
	mwriter := io.MultiWriter(logFile, os.Stdout)
	fmt.Fprintln(mwriter, "-----ELToD v4.0 (c) 2018-----")
	fmt.Fprintf(mwriter, "Start time : %v\n", time.Now().Format("01-02-2006 15:04:05"))
	fmt.Fprintf(mwriter, "Number of CPUs : %v\n", runtime.NumCPU())

	//log file
	choiceLogFile, _ := os.Create(par.ChoiceLogFile)
	defer choiceLogFile.Close()
	fmt.Fprintln(choiceLogFile, "ITERATION,ORIGIN,DESTINATION,START,PERIOD,MODE,FROM_NODE,TO_NODE,DISTANCE1,DISTANCE2,"+
		"TIME1,TIME2,PERCEIVE1,PERCEIVE2,FFTIME1,FFTIME2,TOLL1,TOLL2,RELIABLE1,RELIABLE2,UTILITY,SHARE1,TRIPS1,TRIPS2")

	Links := make(map[int]*tp.Link) //map of link pointers
	Nodes := make(map[int]*tp.Node) //map of node pointers
	prg := []int{0, 0, 0}

	readfile.ReadPar(&par)
	fmt.Fprintf(mwriter, "Number of threads %v.\n", par.NumThreads)

	var wg sync.WaitGroup

	readfile.ReadNode(&par, Nodes, mwriter)
	readfile.ReadLink(&par, Links, Nodes, mwriter)
	readfile.CheckNetwork(&par, Links, Nodes, mwriter)

	//status print
	go func() {
		for range time.NewTicker(time.Second).C {
			fmt.Printf("\b\b\b   \rProcessing tp %v zone %v record %v", prg[0], prg[1], prg[2])
		}
	}()

	for iter := 1; iter <= par.MaxIter; iter++ {
		start := time.Now()
		fmt.Fprintf(mwriter, "\nIteration = %v\n", iter)
		par.Iter = iter
		wg.Add(par.NumThreads)
		TaskQueue := make(chan *tp.PathTask, 100)
		//create workers
		for i := 1; i <= par.NumThreads; i++ {
			go worker(&par, TaskQueue, &wg, mwriter, choiceLogFile, &prg)
		}
		readfile.ReadTT(&par, Nodes, TaskQueue, mwriter)
		wg.Wait()
		pbld.Volsmoothing(&par, Links, mwriter)
		pbld.UpdateLink(&par, Links, mwriter)
		fmt.Printf("\r")
		fmt.Fprintf(mwriter, "Time used = %.2f min                     \n", time.Since(start).Minutes())
		//break
	}
	//Number of splits
	for i, v := range par.SplitCount {
		fmt.Fprintf(mwriter, "Trip Split %v = %v\n", i, v)
	}
	//path build
	fmt.Fprintf(mwriter, "Initial path build  = %v\n", par.PathBuildCount[0])
	fmt.Fprintf(mwriter, "Alternative path build  = %v\n", par.PathBuildCount[1])
	fmt.Fprintf(mwriter, "End time : %v\n", time.Now().Format("01-02-2006 15:04:05"))
}

//worker
func worker(par *tp.Par, Task <-chan *tp.PathTask, wg *sync.WaitGroup, wr io.Writer, cLog io.Writer, prg *[]int) {
	for task := range Task {
		pbld.LoadLink(par, task, pbld.SPath(par, task, cLog), cLog)
		(*prg)[0] = task.TP
		(*prg)[1] = task.I.N
		(*prg)[2] = task.Rcd
	}
	wg.Done()
}
