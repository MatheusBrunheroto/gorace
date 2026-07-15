package main

import (
	"fmt"
	"gorace/display"
	"gorace/input"
	"gorace/log"
	"gorace/request"
	"gorace/request/cache"
	"os"
)

// BASIC TEST: go run main.go -u '1.com' --threads 10 -u '2.com' --threads 20

// 1. Start Cache, the cache channel will be used to avoid reprocessing websites to workers in worker.go.
// 2. Start Display, the display channel will be used to display any kind of thing in worker.go.
/* (This means a worker will be called with both channels as parameters, cache and display). */
// 3. Read the CLI input, will treat the wordlists, websites, etc. So it can be sent to worker.go later on.
// 4. Start Workers, it reads and writes on Cache, and writes on Display (Progress)
/*
 *                      THREADS
 *                   1. ┌─────┐   CHANNELS
 *                 ┌────┤Cache│◄──────────┐
 *                 │    └─────┘           │
 * ┌────────┐      │ 2. ┌───────┐         │
 * │ main.go├──────┴────┤Display│◄─────┐  │
 * └──┬─▲──┬┘           └───────┘      │  │
 * 3. │ │  │                         ┌─┴──▼─┐
 *    │ │  └─────────────────────────┤Worker│
 *   ┌┴─┴┐    4.  init(websites)     └──────┘
 *   │CLI│
 *   └───┘
 *   return
 *  websites
 */

func main() {

	global := input.GlobalFlags{Mode: "flood", Verbosity: 1}
	progress := log.Progress{
		Total:     make(chan int),
		Sent:      make(chan int),
		Succeeded: make(chan int),
		Failed:    make(chan int),
		Finished:  make(chan struct{}),
	}

	// 1. Cache (will later on avoid buildRequest generating the same request multiple times)
	cacheChan := make(chan cache.Operation)
	go cache.Run(cacheChan) // OwO

	// 2. Logger (default verbosity = 1)
	logChan := make(chan log.Entry)
	go log.Run(logChan, &global.Verbosity) // [x] Panic() are not read inside log, as it could run the error before actually stopping it

	// 3. Display
	display.Run(progress.Reader(), logChan)

	// 4. CLI (Read and Filter)
	websites := input.CLI(os.Args[1:], &global, logChan)

	// 5. Workers
	workerChans := request.WorkerChans{
		Progress:  progress.Writer(),
		CacheChan: cacheChan,
	}
	request.InitWorkers(websites, global.Mode, workerChans, logChan)

	fmt.Printf("\n\n")
	<-progress.Finished // Waits for display output of the current session to finish

}

// Fazer SINGLEPACKET, apenas pra modos FLOOD
// Starts output
// FAZER URL LER WORDLIST, SUPORTAR WORDLISTx, le a string inteira pra ver se contem
