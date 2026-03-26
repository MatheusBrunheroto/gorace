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
	request.InitWorker(websites, 20) // ADICIONOAR OPCAO DE THREAD

}
