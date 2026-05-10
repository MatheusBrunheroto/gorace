package request

import (
	"context"
	"fmt"
	"gorace/input"
	"gorace/log"
	"gorace/request/cache"
	"io"
	"net/http"
	"strings"

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

// Always ends up doing N threads to the first website, and N for the other
// Receives a copy, so there is no need to thread lock
func worker(start <-chan struct{}, w input.Website, chans WorkerChans) {

	sent := chans.Progress.Sent
	completed := chans.Progress.Completed

	var request *http.Request
	var err error

	// xSlice(w.Headers, w.Cookies, w.Data)
	hash := computeHash(w)
	copy := cache.Get(hash, chans.CacheChan)

	if copy != nil {
		request = copy.Clone(context.Background())

	} else {
		request, err = buildRequest(w)
		if err != nil {
			return
		}
		cache.Insert(hash, request, chans.CacheChan)
	}

	<-start
	sent <- 1

	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		completed <- 1
		return
	}
	respbody, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		completed <- 1
		return
	}
	resp.Body.Close()

	//fmt.Println(resp.Status, resp.ContentLength)
	if !strings.Contains(string(respbody), "Invalid username or password.") {
		//	fmt.Println(w.Data, resp.Header)
	}
	//We Read the response body on the line below.
	completed <- 1

}
