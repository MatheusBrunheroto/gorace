package request

import (
	"fmt"
	"gorace/input"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
)

func InitWorker(websites []input.Website, amount int) {

	start := make(chan struct{})
	var wg sync.WaitGroup

	wg.Add(amount)

	for i := 0; i < amount; i++ {
		fmt.Println(i)

		go func() {
			defer wg.Done()

			<-start
			Worker(websites)
		}()
	}

	close(start)

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
		fmt.Println(request.ContentLength)
		//We Read the response body on the line below.
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatalln(err)
		}
		//Convert the body to type string
		sb := string(body)
		log.Printf(sb)
	}

}

//func CreateC
