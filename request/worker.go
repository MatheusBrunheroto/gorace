package request

import (
	"context"
	"gorace/input"
	"gorace/log"
	"gorace/log/verbose"
	"gorace/request/cache"
	"net/http"
	"time"
)

type WorkerChans struct {
	Progress  log.ProgressWriter
	CacheChan chan cache.Operation
	LogChan   chan<- log.Entry
}

// Checks for request existence in cache, if it doesn't exist, create a new and insert in cache
func getOrBuildRequest(w input.Config, cacheChan chan cache.Operation) (*http.Request, uint64, bool, error) {

	var request *http.Request
	var err error

	hash := cache.ComputeHash(w)
	hit := false

	if copy := cache.Get(hash, cacheChan); copy != nil {
		request = copy
		hit = true
	} else {
		if request, err = buildRequest(w); err != nil {
			cache.Insert(hash, nil, cacheChan)
			return nil, hash, false, err
		}
		cache.Insert(hash, request, cacheChan)
	}

	request = request.Clone(context.Background())
	if request.GetBody != nil {
		if request.Body, err = request.GetBody(); err != nil {
			return nil, hash, hit, err
		}
	}
	return request, hash, hit, nil

}

// Always ends up doing N threads to the first Config, and N for the other
// Receives a copy, so there is no need to thread lock
func worker(start <-chan struct{}, w input.Config, match string, chans WorkerChans) {

	request, hash, hit, err := getOrBuildRequest(w, chans.CacheChan)
	if err != nil {
		chans.LogChan <- log.Entry{Text: err.Error(), Verbosity: 1}
		chans.Progress.Failed <- 1
		return
	}

	<-start
	verbose.Worker(w, hash, hit, chans.LogChan)

	time.Sleep(time.Duration(w.Delay) * time.Millisecond)
	chans.Progress.Sent <- 1

	client := &http.Client{Timeout: w.Timeout}
	resp, err := client.Do(request)
	if err != nil {
		chans.Progress.Failed <- 1
		verbose.WorkerError(hash, err.Error(), chans.LogChan)
		return
	}

	if err := verbose.WorkerResponse(hash, resp, match, chans.LogChan); err != nil {
		chans.Progress.Failed <- 1
		return
	}

	chans.Progress.Succeeded <- 1

}
