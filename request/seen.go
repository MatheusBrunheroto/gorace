package request

import (
	"net/http"
)

func findRequest()

func getRequest()

func requestExists() {

}

func lookupRequest() {

	hashes := make(map[string]*http.Request) // NA VERDADE AQUI SERIA UM UM HASH PRA UM WEBSITE
	for {

		requestedHash := <-hashChan
		for h, r := range hashes {
			if h == requestedHash {
				chan<- r
			}

		}

	}

}

func requestRegistry() {

	hash := make(chan string) // CHANNEL NO NOME
	request := make(chan *http.Request)
	hashMap := make(map[string]*http.Request) // hash -> url + ... + threads

	for {
		select {
		case x := <-input:

		case y := <-request:
			h, exist := hashMap[y]
			if !exist {
				// ADD
			} else {
				// RETURN
			}

			output <- result
		}

	}

	hashes := make(map[string]*http.Request)

}
