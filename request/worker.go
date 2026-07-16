package request

import (
	"context"
	"gorace/input"
	"gorace/log"
	"gorace/log/verbose"
	"gorace/request/cache"
	"io"
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

	if copy := cache.Get(hash, cacheChan); copy != nil {
		request = copy.Clone(context.Background()) // !!! Does not clone BODY
		return request, hash, true, nil
	}

	if request, err = buildRequest(w); err != nil {
		return nil, 0, false, err
	}
	cache.Insert(hash, request, cacheChan)
	return request.Clone(context.Background()), hash, false, nil

}

// Always ends up doing N threads to the first Config, and N for the other
// Receives a copy, so there is no need to thread lock
func worker(start <-chan struct{}, w input.Config, chans WorkerChans) {

	request, hash, hit, err := getOrBuildRequest(w, chans.CacheChan)
	if err != nil {
		chans.LogChan <- log.Entry{Text: err.Error(), Verbosity: 1}
		return
	}

	<-start
	verbose.Worker(w, hash, hit, chans.LogChan)

	time.Sleep(time.Duration(w.Delay) * time.Millisecond)
	chans.Progress.Sent <- 1

	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		chans.Progress.Failed <- 1
		verbose.WorkerError(hash, err.Error(), chans.LogChan)
		return
	}
	_ = resp
	respbody, err := io.ReadAll(resp.Body)
	if err != nil {
		verbose.WorkerError(hash, err.Error(), chans.LogChan)
		return
	}

	chans.LogChan <- log.Entry{Text: string(respbody), Verbosity: 3}
	resp.Body.Close()

	chans.Progress.Succeeded <- 1

}
