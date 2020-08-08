package ConversionWorker

import "proj2/ImageConversion"

type Request struct {
	// img src
	Source string
	// ascii code output path
	OutPath string
	// all request generation done flag
	AllDone bool
	// requests finish channel
	ReqDone chan bool
	Flag chan bool
}

type Worker struct {
	// heap index
	Index int
	// its own requests channel
	OwnRequests chan Request
	// number of waiting request for this worker
	WaitingRequests int
}

func (w *Worker) ConversionWork(done chan *Worker, allFinished *bool) {
	for {
		// return if all work is finished
		if *allFinished{
			return
		}
		// retrieve conversion request from its own request channel
		req := <-w.OwnRequests
		// write to its resp channel
		source := req.Source
		outpath := req.OutPath
		allDone := req.AllDone
		req.ReqDone <- ImageConversion.ImageToAscii(source, outpath, allDone, allFinished)
		done <- w
	}
}
