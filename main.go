package main

import (
	"fmt"
	readfile "readfile"
	"sync"
	"time"
	tp "typedef"
	pb "pathbuilding"
)

func main() {
	fmt.Println("//****ELToD 4 (c) 2018****//")

	Links := make(map[int]*tp.Link) //map of link pointers
	Nodes := make(map[int]*tp.Node) //map of node pointers
	par := tp.Par{}                 //parameters

	readfile.ReadPar(&par)
	readfile.ReadNode(&par, Nodes)
	readfile.ReadLink(&par, Links, Nodes)

	//make task queue
	par.NumThreads = 1
	TaskQueue := make(chan *tp.PathTask, 100)

	//create workers
	var wg sync.WaitGroup
	wg.Add(par.NumThreads)
	for i := 1; i <= par.NumThreads; i++ {
		go worker(&par, i, TaskQueue, &wg)
	}

	//read trip table
	readfile.ReadTT(&par, Nodes, TaskQueue)

	wg.Wait()
	fmt.Printf("end of program")
}

//worker
func worker(par *tp.Par, id int, Task <-chan *tp.PathTask, wg *sync.WaitGroup) {
	for task := range Task {
		pb.SPath(par, task)
		time.Sleep(1 * time.Second)
		fmt.Printf("worker %v - task %+v\n", id, task)
	}
	wg.Done()
}
