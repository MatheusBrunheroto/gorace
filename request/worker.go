package request

import (
	"context"
	"fmt"
	"gorace/input"
	"gorace/log"
	"io"
	"net/http"
	"strings"
)

// Always ends up doing N threads to the first website, and N for the other
// Receives a copy, so there is no need to thread lock
func worker(start <-chan struct{}, w input.Website,
	progress log.ProgressWriter) {

	sent := progress.Sent
	completed := progress.Completed

	var request *http.Request
	var err error

	// xSlice(w.Headers, w.Cookies, w.Data)
	hash := computeHash(w)
	copy := getRequest(hash, registryChan)

	if copy != nil {
		request = copy.Clone(context.Background())

	} else {
		request, err = buildRequest(w)
		if err != nil {
			return
		}
		insertRequest(hash, request, entryChan)
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
	progress.Completed <- 1

}
