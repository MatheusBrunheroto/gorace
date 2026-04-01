package request

import (
	"fmt"
	"gorace/input"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
)

func InitWorker(websites []input.Website, amount int) {

	start := make(chan struct{})
	var wg sync.WaitGroup

	wg.Add(amount)

	for i := 0; i < amount; i++ {

		// Wait for all the workers to be initialized, and start the requests at the same time
		go func() {
			defer wg.Done()
			<-start
			Worker(websites)
		}()
	}

	close(start) // Closing the channel is the trigger

	wg.Wait()
}

// Always ends up doing N threads to the first website, and N for the other
// Receives a copy, so there is no need to thread lock
func Worker(websites []input.Website) {

	fmt.Println(websites)

	for _, v := range websites {

		client := &http.Client{}
		var data url.Values
		var body *strings.Reader
		var request *http.Request
		var err error

		switch v.Method {

		case "POST", "PUT", "PATCH": // Has body

			data = url.Values{}
			for _, d := range v.Data {
				data.Set(d.Key, d.Value)
			}
			body = strings.NewReader(data.Encode()) // Turns k1:v1 and k2:v2 to k1=v1&k2=v2

			request, err = http.NewRequest(v.Method, v.Url, body)
			request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		default:
			body = nil
		}

		if err != nil {
			fmt.Printf("client: could not create request: %s\n", err)
			os.Exit(1)
		}

		// Won't do anything if v.Headers or v.Cookies are empty, no need to check
		for _, h := range v.Headers {
			request.Header.Set(h.Key, h.Value)
		}
		for _, c := range v.Cookies {
			request.AddCookie(&http.Cookie{Name: c.Key, Value: c.Value})
		}

		resp, err := client.Do(request)
		if err != nil {
			panic(err)
		}
		fmt.Println(resp.Status, resp.ContentLength)
		//We Read the response body on the line below.

	}

}

//func CreateC
