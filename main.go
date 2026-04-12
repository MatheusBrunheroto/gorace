package main

import (
	"fmt"
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

*/

func init() {

	fmt.Println("\n")
	fmt.Println(`//  ▄▄ •       ▄▄▄   ▄▄▄·  ▄▄· ▄▄▄ .
// ▐█ ▀ ▪▪     ▀▄ █·▐█ ▀█ ▐█ ▌▪▀▄.▀·
// ▄█ ▀█▄ ▄█▀▄ ▐▀▀▄ ▄█▀▀█ ██ ▄▄▐▀▀▪▄
// ▐█▄▪▐█▐█▌.▐▌▐█•█▌▐█ ▪▐▌▐███▌▐█▄▄▌
// ·▀▀▀▀  ▀█▄▀▪.▀  ▀ ▀  ▀ ·▀▀▀  ▀▀▀ `)
	fmt.Println("\n")

}

func main() {

	args := os.Args[:1]

	var err error
	var websites []input.Website
	if websites, err = input.RunCLI(); err != nil {
		fmt.Println(err)
		return
	}

	for _, w := range websites {
		fmt.Println(w)
	}

	var mode string = "flood"
	modes := []string{"sequential", "parallel", "flood"}

	for i := range args {
		if args[i] == "-m" || args[i] == "--mode" {
			mode = args[i+1]
			break
		}
	}
	var modeExists bool = false
	for _, m := range modes {
		if mode == m {
			modeExists = true
		}
	}
	if !modeExists {
		mode = "flood"
		fmt.Println("Unable to determine mode \"" + mode + "\" using \"flood\" as default...")
	}
	request.InitWorker(websites, mode)
	//close(jobs)

}

// TODO, MODO de input direto de wordlist, MODOS DE rodar,
