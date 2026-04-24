package main

import (
	"fmt"
	"gorace/display"
	"gorace/input"
	"gorace/request"
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

go main.go -u '1.com' --threads 10 -u '2.com' --threads 20


*/

// Starts output

func main() {

	// fazer array de canais, um pra avisar que a wordlist ta pronta e mandar a quantidade de palavras, outro pra sent e ouj ro pra completed
	//progressChannel := make(chan int, 2)
	progressChannel := [3]chan int{make(chan int), make(chan int), make(chan int)}
	display.Display(progressChannel) // call for the amount of words too

	// Reads and filter the CLI inputs
	websites, mode, err := input.RunCLI(os.Args[1:])
	if err != nil {
		fmt.Println(err, "\nExiting...\n")
		return
	}
	var totalRequests int = 0
	for _, w := range websites {
		totalRequests += w.Threads
	}
	progressChannel[0] <- totalRequests // Initializes progress bar inside display.go

	request.InitWorker(progressChannel[:1], websites, mode)

	fmt.Println("\n")
}

// TODO, MODO de input direto de wordlist, MODOS DE rodar,
// AO inves de retornar os erros e tentar fazer funcionar denovo basta colocar o panic ao inves do erros.New
