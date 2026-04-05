package request

import (
	"fmt"
	"gorace/input"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
)

func InitWorker(start <-chan struct{}, jobs <-chan input.Website, threads int) {

	var wg sync.WaitGroup

	for i := 0; i < threads; i++ {

		// Wait for all the workers to be initialized, and start the requests at the same time
		wg.Go(func() {
			Worker(start, jobs, &wg)
		})
	}

	wg.Wait()

}

// Always ends up doing N threads to the first website, and N for the other
// Receives a copy, so there is no need to thread lock
func Worker(start <-chan struct{}, jobs <-chan input.Website, wg *sync.WaitGroup) {

	<-start

	fmt.Println(jobs)
	for w := range jobs {
		wg.Add(1)
		defer wg.Done()
		//		largest := maxSlice(w.Headers, w.Cookies, w.Data)

		client := &http.Client{}
		var data url.Values
		var body *strings.Reader
		var request *http.Request
		var err error

		switch w.Method {

		case "POST", "PUT", "PATCH": // Has body

			data = url.Values{}
			for _, d := range w.Data {
				data.Set(d.Key, d.Value)
			}
			body = strings.NewReader(data.Encode()) // Turns k1:v1 and k2:v2 to k1=v1&k2=v2

			request, err = http.NewRequest(w.Method, w.Url, body)
			request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		default:
			body = nil
		}

		if err != nil {
			fmt.Printf("client: could not create request: %s\n", err)
			os.Exit(1)
		}

		// Won't do anything if v.Headers or v.Cookies are empty, no need to check
		for _, h := range w.Headers {
			if h.Key == "" {
				continue
			}
			request.Header.Set(h.Key, h.Value)
		}
		// melhorar o filtro desses valores
		for _, c := range w.Cookies {
			if c.Key == "" {
				continue
			}
			request.AddCookie(&http.Cookie{Name: c.Key, Value: c.Value})
		}

		resp, err := client.Do(request)
		if err != nil {
			panic(err)
		}
		respbody, err := io.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}
		resp.Body.Close()

		fmt.Println(resp.Status, resp.ContentLength)
		if !strings.Contains(string(respbody), "Invalid username or password.") {
			fmt.Println(w.Data, resp.Header)
		}
		//We Read the response body on the line below.

	}

}
