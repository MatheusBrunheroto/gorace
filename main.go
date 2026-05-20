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

// usage example: gorace -u 'https://website.com' -h '{header_name:header_value, h2_name:WORDLIST1}' -c '{WORDLIST2:WORDLIST3}' -t 50 --no-filter
/*
	MODES:
		-c --cli, enables **command line mode**
			The CLI mode will take only one website a time, unless a wordlist is provided,
			in this case, will use two different modes, pitchfork (default) or clusterbomb

		-i --interative, enables **interative mode**
			The interative mode is not very pratical, but can be used to generate your own requests and save them in a
			custom wordlist

		-w --wordlist, enables **wordlist mode**
			It's different from the --cli wordlists, it takes pre-made requests and make them, it's recommended to create the
			wordlists using the interative mode

	OPTIONS:
		--no-filter, skips the data filtering, "wrong" parameters will be able to pass, for example:
		using a header with ":HEADER_NAME:::::Content%!%!%@123123" will be passed in it's raw form to the request"

		-t --threads, number of workers or agents to be used in the test, default=50

go main.go run -u '1.com' --threads 10 -u '2.com' --threads 20


*/
// Fazer SINGLEPACKET, apenas pra modos FLOOD
// Starts output

func main() {

	// CACHE
	cacheChan := make(chan cache.Operation)
	go cache.Run(cacheChan) // OwO

	progress := log.Progress{
		Total:     make(chan int),
		Sent:      make(chan int),
		Succeeded: make(chan int),
		Failed:    make(chan int),
	}

	// DISPLAY ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
	session := display.NewSession(progress.Reader())
	go display.Run(session)
	<-session.Ready

	// Read and filter the CLI inputs
	// session.Draw <- "⸺" // ⸺⸺⸺⸺⸺⸺⸺⸺⸺⸺⸺
	websites, mode, err := input.RunCLI(os.Args[1:])
	if err != nil {
		fmt.Println(err, "\nExiting...\n")
		return
	}
	// session.Draw <- "⸺" // ⸺⸺⸺⸺⸺⸺⸺⸺⸺⸺⸺

	// This is a memory that runs with the code, avoiding buildRequest to generate the same request multiple times

	workerChans := request.WorkerChans{
		Progress:  progress.Writer(),
		CacheChan: cacheChan,
	}
	request.InitWorkers(websites, mode, workerChans)
	// ADD val padrao de threads pra 0
	fmt.Printf("\n\n")
	<-session.Finished // Waits for display output to finish

}

// TODO, MODO de input direto de wordlist, MODOS DE rodar,
// AO inves de retornar os erros e tentar fazer funcionar denovo basta colocar o panic ao inves do erros.New
