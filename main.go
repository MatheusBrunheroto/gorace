package main

import (
	"fmt"
	"gorace/input"
	"gorace/request"
	"os"
)

func error() {
	fmt.Println("Wrong usage, please specify mode: \"-c\", \"-w\", \"-i\"; or \"--cli\", \"--wordlist\", \"--iterative\".")
	fmt.Println("For more info use: \"-h\" or \"--help\".")
}

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

	// Guarantee that only one "mode" flag is passed
	modes := 0
	if os.Args[1] == "-c" || os.Args[1] == "--cli" {
		modes++
	}
	if os.Args[1] == "-w" || os.Args[1] == "--wordlist" {
		modes++
	}
	if os.Args[1] == "-i" || os.Args[1] == "--interative" {
		modes++
	}
	if modes != 1 {
		error()
		return
	}

	// Only gets to here if a "mode" flag exists
	var websites []input.Website

	switch os.Args[1] {

	case "-c", "--cli":
		if err := input.RunCLI(&websites); err != nil {
			fmt.Println(err)
			return
		}

	case "-w", "--wordlist":
		fmt.Println("Modo Wordlist")

	case "-i", "--iterative":
		input.GetTargetInfo(&websites)

	}

	fmt.Println(websites)
	request.InitWorker(websites, 50) // ADICIONOAR OPCAO DE THREAD

}
