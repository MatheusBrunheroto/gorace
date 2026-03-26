package request

import (
	"fmt"
	"gorace/input"
	"net/http"
	"os"
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

	var methods = [7]string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"}

	for _, v := range websites {

		request, err := http.NewRequest(methods[v.Method], v.Link, nil)
		if err != nil {
			fmt.Printf("client: could not create request: %s\n", err)
			os.Exit(1)
		}
		client := &http.Client{}

		resp, err := client.Do(request)
		if err != nil {
			panic(err)
		}
		fmt.Println(resp.ContentLength)
		//We Read the response body on the line below.

	}

}

//func CreateC
