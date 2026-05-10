package request

import (
	"net/http"
)

func findRequest()

func getRequest()

func requestExists() {

}

func lookupRequest(hashChannel chan string) {

	s
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

func requestRegistry(requestChannel chan *http.Request) {

	hashMap := make(map[string]*http.Request) // hash -> url + ... + threads

	for {
		select {
		case candidateHash := <-hashChannel:
			_, exist := hashMap[candidateHash]
			if !exist{
				hashMap[candidateHash] = // URL
			}

		case requestedHash := <-requestChannel:

			h, exist := hashMap[requestedHash]
			if exist {
				// RETURN
			}
			output <- result
		}

	}

	hashes := make(map[string]*http.Request)

}
