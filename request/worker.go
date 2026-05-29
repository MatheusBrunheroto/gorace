package request

import (
	"context"
	"fmt"
	"gorace/input"
	"gorace/log"
	"gorace/request/cache"
	"io"
	"net/http"
	"time"

	"github.com/cespare/xxhash/v2"
)

type WorkerChans struct {
	Progress  log.ProgressWriter
	CacheChan chan cache.Operation
}

func computeHash(w input.Website) uint64 {
	code := fmt.Sprintf("%s%s%s%s%s%d", w.Url, w.Method, w.Headers, w.Cookies, w.Data, w.Threads)
	return xxhash.Sum64String(code)
}

// Checks for request existence in cache, if it doesn't exist, create a new and insert in cache
func getOrBuildRequest(w input.Website, cacheChan chan cache.Operation) (*http.Request, error) {

	var request *http.Request
	var err error

	hash := computeHash(w)

	if copy := cache.Get(hash, cacheChan); copy != nil {
		request = copy.Clone(context.Background()) // Does not clone BODY
		return request, nil
	}

	if request, err = buildRequest(w); err != nil {
		return nil, err
	}
	cache.Insert(hash, request, cacheChan)
	return request, nil

}

// Always ends up doing N threads to the first website, and N for the other
// Receives a copy, so there is no need to thread lock
func worker(start <-chan struct{}, w input.Website, chans WorkerChans) {

	request, err := getOrBuildRequest(w, chans.CacheChan)
	if err != nil {
		fmt.Println(err)
		return
	}

	<-start
	time.Sleep(time.Duration(w.Delay) * time.Millisecond)
	chans.Progress.Sent <- 1

	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		chans.Progress.Failed <- 1
		fmt.Println(err)
		return
	}
	_ = resp
	respbody, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(w)
	fmt.Println(string(respbody))
	resp.Body.Close()

	chans.Progress.Succeeded <- 1

}
