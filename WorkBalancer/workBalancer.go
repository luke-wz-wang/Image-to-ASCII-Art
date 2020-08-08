package WorkBalancer

import (
	"container/heap"
	"fmt"
	"proj2/ConversionWorker"
	"proj2/WorkerPool"
)

type WorkBalancer struct {
	// a pool containing a number of #numOfWorkers workers
	WorkerPool WorkerPool.WorkerPool
	// conversion work done channel
	Done chan *ConversionWorker.Worker
}

// initialize the work balancer with #numOfWorkers workers
func Init(allFinished *bool, magnitude int, numOfWorkers int ) *WorkBalancer {
	done := make(chan *ConversionWorker.Worker, numOfWorkers)
	// create numOfWorkers channels
	wb := &WorkBalancer{make(WorkerPool.WorkerPool, 0, numOfWorkers), done}
	for i := 0; i < numOfWorkers; i++ {
		w := &ConversionWorker.Worker{OwnRequests: make(chan ConversionWorker.Request, magnitude)}
		// push created worker to heap
		heap.Push(&wb.WorkerPool, w)
		// send worker to start their conversion work
		go w.ConversionWork(wb.Done, allFinished)
	}
	return wb
}

func (wb *WorkBalancer) BalanceWork(req chan ConversionWorker.Request, allFinished *bool) {
	for {
		if *allFinished {
			return
		}
		select {
		// retrieve conversion request from req channel
		case request := <-req:
			wb.DistributeWork(request)
		// retrieve work complete msg from Done channel
		case w := <-wb.Done:
			wb.Renew(w)
		}
		// analyze the status of work balancing
		wb.Analyze()
	}
}

// distribute work to the worker with the least work
func (wb *WorkBalancer) DistributeWork(req ConversionWorker.Request) {
	// pop worker with the least work
	workerWithLeastWork := heap.Pop(&wb.WorkerPool).(*ConversionWorker.Worker)
	workerWithLeastWork.OwnRequests <- req
	workerWithLeastWork.WaitingRequests++
	// push it back into heap
	heap.Push(&wb.WorkerPool, workerWithLeastWork)
}

// renew a worker by removing worker from heap and pushing it back again
func (wb *WorkBalancer) Renew(worker *ConversionWorker.Worker) {
	worker.WaitingRequests--
	// remove worker from heap and then push it back
	heap.Remove(&wb.WorkerPool, worker.Index)
	heap.Push(&wb.WorkerPool, worker)
}

//analyze the status of work balancing
func (wb *WorkBalancer) Analyze() {
	for _, worker := range wb.WorkerPool {
		fmt.Printf("%d ", worker.WaitingRequests)
	}
	fmt.Printf("\n")
}
