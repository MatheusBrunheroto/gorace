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

func maxThreads(websites []input.Website) int {
	n := 1
	for _, w := range websites {
		if w.Threads > n {
			n = w.Threads
		}
	}
	return n
}
func initChan(n int) []chan struct{} {
	c := make([]chan struct{}, n) // Each chan here is a "thread"
	for i := range c {
		c[i] = make(chan struct{})
	}
	return c
}

func normalWorker(start chan struct{}, website input.Website, wg *sync.WaitGroup) {
	for range website.Threads {
		wg.Go(func() { Worker(start, website, wg) })
	}
	close(start)
}

func roundWorker(start chan struct{}, websites []input.Website, wg *sync.WaitGroup) {
	for _, w := range websites {
		wg.Go(func() { Worker(start, w, wg) })
	}
	close(start)
}

func initNormal(websites []input.Website, isSequential bool) {

	startChans := initChan(len(websites))
	var cascateWg sync.WaitGroup

	for i, w := range websites {

		if isSequential {
			var sequentialWg sync.WaitGroup
			normalWorker(startChans[i], w, &sequentialWg)
			sequentialWg.Wait()
			continue
		}
		normalWorker(startChans[i], w, &cascateWg)

	}

	if isSequential {
		return
	}
	cascateWg.Wait()

}

func initRound(websites []input.Website, isSequential bool) {

	startChannels := initChan(maxThreads(websites))
	var cascateWg sync.WaitGroup

	for _, c := range startChannels {

		if isSequential {
			var sequentialWg sync.WaitGroup
			roundWorker(c, websites, &sequentialWg)
			sequentialWg.Wait()
			continue

		}
		roundWorker(c, websites, &cascateWg)

	}

	if isSequential {
		return
	} else {
		cascateWg.Wait()
	}

}

func InitWorker(websites []input.Website, mode string) {

	switch mode {

	// After N threads of an URL requests were sent to worker, waits for them to finish before starting next URL requests
	case "sequential":
		initNormal(websites, true)
	// Same as sequential, but doesn't wait for its requests to finish before starting the next URL requests
	case "cascade":
		initNormal(websites, false)

	// Sequential's behaviour, but cycles through the URLs requests for N times, N = largest amount of threads informed
	case "round-sequential":
		initRound(websites, true)
	// Cascade's behaviour, but cycles through the URLs requests for N times, N = largest amount of threads informed
	case "round-cascade":
		initRound(websites, false)

	// This is the default mode "flood", group all the requests and fire them at the exact same moment
	default:
		start := make(chan struct{})
		var wg sync.WaitGroup
		for _, w := range websites {
			for range w.Threads {
				wg.Go(func() { Worker(start, w, &wg) })
			}
		}
		close(start)
		wg.Wait()
	} // WG GO INTERNO PRA GERAR COISA, ESPERAR POR DENTRO, ALGO DO TIPO, E UM POR FORA SEI LA

}

// Always ends up doing N threads to the first website, and N for the other
// Receives a copy, so there is no need to thread lock
func Worker(start <-chan struct{}, w input.Website, wg *sync.WaitGroup) {

	<-start

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
