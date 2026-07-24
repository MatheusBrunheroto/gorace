package input

import (
	"fmt"
	"gorace/assets"
	"gorace/log"
	"os"
	"strconv"
)

type GlobalFlags struct {
	Mode      string
	Match     string
	Verbosity int
	NoColor   bool
}

func readGlobalFlags(args []string, global *GlobalFlags, logChan chan<- log.Entry) {

	modes := []string{"sequential", "round-sequential", "cascade", "round-cascade", "flood"}
	var modeExists bool = false
	var matchExists bool = false

	for i, flag := range args {

		switch flag {
		case "--mode":
			for _, mode := range modes {
				if mode == args[i+1] {
					modeExists = true
					global.Mode = mode
					break
				}
			}

		case "--verbose":
			if v, err := strconv.Atoi(args[i+1]); err == nil {

				switch {
				case v < 0:
					global.Verbosity = 0
				case v >= 4:
					global.Verbosity = 4
				default:
					global.Verbosity = v
				}

			}
		case "--match":
			global.Match = args[i+1]
			matchExists = true
		case "--no-color":
			global.NoColor = true

		case "--help":
			fmt.Println(assets.Help)
			os.Exit(0)
		}

		// The code won't return anything about the responses unless a "match" flag is sent
	}
	if !modeExists {
		logChan <- log.Entry{Text: "[!] Mode wasn't identified, using \"flood\" as default...", Verbosity: 1}
	}
	if !matchExists {
		logChan <- log.Entry{Text: "[!] Response match wasn't identified, increase verbosity to check for feedback", Verbosity: 1}
	}
	if !global.NoColor {
		logChan <- log.Entry{Text: "[!] It's recommended to run \"--no-color\" if you want to output the results to a file", Verbosity: 1}
	}

}

func checkFlagsConsistency(args []string) {

	// Append
	flags := initGlobalFlags()
	for key, value := range initConfigFlags() {
		flags[key] = value
	}

	for i, arg := range args {

		if _, ok := flags[arg]; ok {

			if arg == "--no-color" || arg == "--help" {
				continue
			}

			if i+1 < len(args) {
				if _, ok := flags[args[i+1]]; ok {
					panic(fmt.Sprintf("[x] Expected a value, got another flag instead in \"%s %s\"! Exiting...", arg, args[i+1]))
				}
				continue
			}

			panic(fmt.Sprintf("[x] A flag is missing a parameter in \"%s\"! Exiting...", arg))

		}

		// nao le caso onde tem 1 arg sem flag a mais

	}

}

// Take the abreviation -f of --flag, and turns it into --flag, because it makes dealing with the flags from initFlags()
func normalizeInputFlags(args *[]string) {

	table := map[string]string{
		// General
		"-h": "--help",

		// Request
		"-u": "--url",
		"-X": "--method",
		"-H": "--headers",
		"-b": "--cookies",
		"-d": "--data",
		"-w": "--wordlist",

		// Execution
		"-m": "--mode",
		"-t": "--threads",
		"-D": "--delay",
		"-M": "--match",
		"-v": "--verbose",
	}

	for i := 0; i < len(*args); i++ {
		if normalized, ok := table[(*args)[i]]; ok {
			(*args)[i] = normalized
		}
	}

}

func initGlobalFlags() map[string]string {
	return map[string]string{
		"--help":     "", // -h
		"--no-color": "",
		"--mode":     "", // -m
		"--match":    "", // -M
		"--verbose":  "", // -v
	}
}
func initConfigFlags() map[string]string {
	return map[string]string{
		"--url":      "", // -u
		"--method":   "", // -X
		"--headers":  "", // -H
		"--cookies":  "", // -b
		"--data":     "", // -d
		"--wordlist": "", // -w
		"--threads":  "", // -t
		"--delay":    "", // -D
	}
}
