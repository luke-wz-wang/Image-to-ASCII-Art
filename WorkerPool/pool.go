package WorkerPool

import "proj2/ConversionWorker"

type WorkerPool []*ConversionWorker.Worker

/* override heap interface */

func (wp *WorkerPool) Swap(m, n int) {

	a := *wp
	tmp := a[n]
	a[n] = a[m]
	a[m] = tmp
	a[m].Index = m
	a[n].Index = n
}

func (wp WorkerPool) Len() int {

	return len(wp)
}

func (wp WorkerPool) Less(m, n int) bool {

	return wp[m].WaitingRequests < wp[n].WaitingRequests
}

func (wp *WorkerPool) Pop() interface{} {

	past := *wp
	num := len(past)
	*wp = past[0 : num-1]

	elem := past[num-1]
	elem.Index = -1
	return elem
}

func (wp *WorkerPool) Push(p interface{}) {

	num := len(*wp)
	wpElem := p.(*ConversionWorker.Worker)
	wpElem.Index = num
	*wp = append(*wp, wpElem)
}

