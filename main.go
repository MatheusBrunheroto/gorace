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

	// Cache (will later on avoid buildRequest generating the same request multiple times)
	cacheChan := make(chan cache.Operation)
	go cache.Run(cacheChan) // OwO

	progress := log.Progress{
		Total:     make(chan int),
		Sent:      make(chan int),
		Succeeded: make(chan int),
		Failed:    make(chan int),
	}

	// Display
	session := display.NewSession(progress.Reader())
	go display.Run(session)
	<-session.Ready // Waits for ascii_art.go

	// CLI (Read and filter the inputs)
	websites, mode, err := input.RunCLI(os.Args[1:])
	if err != nil {
		fmt.Println(err, "\nExiting...")
		fmt.Println("")
		return
	}
	for _, w := range websites {
		fmt.Println(w)
		fmt.Println()
	}

	// session.Draw <- "⸺" // ⸺⸺⸺⸺⸺⸺⸺⸺⸺⸺⸺

	// Worker
	workerChans := request.WorkerChans{
		Progress:  progress.Writer(),
		CacheChan: cacheChan,
	}
	request.InitWorkers(websites, mode, workerChans)

	fmt.Printf("\n\n")
	<-session.Finished // Waits for display output of the current session to finish

}

// ADD val padrao de threads pra 0
// TODO, MODO de input direto de wordlist, MODOS DE rodar,
// AO inves de retornar os erros e tentar fazer funcionar denovo basta colocar o panic ao inves do erros.New
// Fazer SINGLEPACKET, apenas pra modos FLOOD
// Starts output
