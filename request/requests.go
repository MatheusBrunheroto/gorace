package request

import (
	"fmt"
	"gorace/input"
	"net/http"
	"os"
)

func RequestThread(websites []input.Website) {

	var methods = [7]string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"}

	for _, v := range websites {

		req, err := http.NewRequest(methods[v.Method], v.Link, nil)
		if err != nil {
			fmt.Printf("client: could not create request: %s\n", err)
			os.Exit(1)
		}
		fmt.Println(req)
		fmt.Println(v.Link)
	}

}
