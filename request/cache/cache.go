package cache

import (
	"fmt"
	"gorace/input"
	"net/http"

	"github.com/cespare/xxhash/v2"
)

type Operation struct {
	Hash      uint64
	Request   *http.Request
	ReplyChan chan *http.Request
}

func ComputeHash(w input.Config) uint64 {
	code := fmt.Sprintf("%s%s%s%s%s%d", w.Url, w.Method, w.Headers, w.Cookies, w.Data, w.Threads)
	return xxhash.Sum64String(code)
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

	hashMap := make(map[uint64]*http.Request)
	waiters := make(map[uint64][]chan *http.Request) // quem está esperando esse hash ficar pronto
	pending := make(map[uint64]bool)                 // hash já está sendo construído por alguém

	for {

		r := <-ch

		if r.Request == nil { // If has no "Request" in r, the operation was "getRequest"

			if req, ok := hashMap[r.Hash]; ok {
				r.ReplyChan <- req

			} else if pending[r.Hash] { // Hash already being built, wait
				waiters[r.Hash] = append(waiters[r.Hash], r.ReplyChan)

			} else {
				pending[r.Hash] = true
				r.ReplyChan <- nil
			}

		} else { // If has "Request" in r, the operation was "insertRequest"

			hashMap[r.Hash] = r.Request

			delete(pending, r.Hash)
			for _, w := range waiters[r.Hash] {
				w <- r.Request // Cache Hit on requests with same hash that were in the "Insertion Queue"
			}
			delete(waiters, r.Hash)

		}
	}
}
