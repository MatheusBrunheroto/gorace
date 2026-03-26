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

	for _, v := range websites {

		data := url.Values{}
		for data_name, data_value := range v.Data {
			data.Set(data_name, data_value)
		}

		body := strings.NewReader(data.Encode())

		request, err := http.NewRequest(v.Method, v.Url, body)
		if err != nil {
			fmt.Printf("client: could not create request: %s\n", err)
			os.Exit(1)
		}
		client := &http.Client{}
		request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		for header_name, header_value := range v.Headers {
			request.Header.Set(header_name, header_value)
		}
		for cookie_name, cookie_value := range v.Cookies {
			request.AddCookie(&http.Cookie{Name: cookie_name, Value: cookie_value})
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
