package main

import (
	"fmt"
	"gorace/input"
	"gorace/request"
)

func error() {
	fmt.Println("Wrong usage, please specify mode: \"-c\", \"-w\", \"-i\"; or \"--cli\", \"--wordlist\", \"--iterative\".")
	fmt.Println("For more info use: \"-h\" or \"--help\".")
}

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

*/

func main() {

	jobs := make(chan input.Website, 1000)
	start := make(chan struct{})

	var threads int
	if err := input.RunCLI(start, jobs, &threads); err != nil {
		fmt.Println(err)
		return
	}

	// fmt.Println(websites)// ADICIONOAR OPCAO DE THREAD7
	request.InitWorker(start, jobs, threads)
	close(jobs)
	fmt.Println("PORRAS RECHEADAS")
}
