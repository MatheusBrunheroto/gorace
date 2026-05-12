package cache

import (
	"net/http"
)

type Operation struct {
	Hash      uint64
	Request   *http.Request
	ReplyChan chan *http.Request
}

func Get(hash uint64, ch chan<- Operation) *http.Request {
	reply := make(chan *http.Request)
	ch <- Operation{
		Hash:      hash,
		Request:   nil,
		ReplyChan: reply,
	}
	return <-reply // Reply could be either nil or the request, this is verified back in worker.go
}
func Insert(hash uint64, request *http.Request, ch chan<- Operation) {
	ch <- Operation{
		Hash:      hash,
		Request:   request,
		ReplyChan: nil,
	}
}
func Run(ch <-chan Operation) {

	hashMap := make(map[uint64]*http.Request) // hash -> url + ... + threads

	for {
		r := <-ch

		if r.Request == nil { // If has no "Request" in r, the operation was "getRequest"
			request, exists := hashMap[r.Hash]
			if !exists {
				r.ReplyChan <- nil
			} else {
				r.ReplyChan <- request
			}

		} else { // If has "Request" in r, the operation was "insertRequest"
			hashMap[r.Hash] = r.Request
		}
	}

}
