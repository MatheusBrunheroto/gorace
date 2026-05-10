package request

import (
	"fmt"
	"gorace/input"
	"net/http"

	"github.com/cespare/xxhash/v2"
)

type RegistryOp struct {
	Hash      uint64
	Request   *http.Request
	ReplyChan chan *http.Request
}

func computeHash(w input.Website) uint64 {
	code := fmt.Sprintf("%s%s%s%s%s%d", w.Url, w.Method, w.Headers, w.Cookies, w.Data, w.Threads)
	return xxhash.Sum64String(code)
}

func getRequest(hash uint64, registryChan chan<- RegistryOp) *http.Request {
	reply := make(chan *http.Request)
	registryChan <- RegistryOp{
		Hash:      hash,
		Request:   nil,
		ReplyChan: reply,
	}
	return <-reply // Reply could be either nil or the request, this is verified back in the worker.go
}
func insertRequest(hash uint64, request *http.Request, registryChan chan<- RegistryOp) {
	registryChan <- RegistryOp{
		Hash:      hash,
		Request:   request,
		ReplyChan: nil,
	}
}

func RunRegistry(lookupChannel <-chan RegistryOp) {

	hashMap := make(map[uint64]*http.Request) // hash -> url + ... + threads

	for {
		r := <-lookupChannel

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
