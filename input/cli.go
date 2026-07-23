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

func removeGlobalFlags(args []string) []string {

	var output []string
	targets := map[string]bool{
		"-m":        true,
		"--mode":    true,
		"-v":        true,
		"--verbose": true,
		"--match":   true,
	}

	for i := 0; i < len(args); i++ {

		arg := args[i]

		if targets[args[i]] {
			if i+1 < len(args) {
				i++
				continue
			}
		}

		if args[i] == "--no-color" {
			continue
		}

		output = append(output, arg)

	}

	return output
}

// Mode and Verbosity
func readGlobalFlags(args []string, global *GlobalFlags, logChan chan<- log.Entry) {

	modes := []string{"sequential", "round-sequential", "cascade", "round-cascade", "flood"}
	var modeExists bool = false
	var matchExists bool = false

	for i, flag := range args {

		if flag == "-m" || flag == "--mode" {

			for _, mode := range modes {
				if mode == args[i+1] {
					modeExists = true
					global.Mode = mode
					break
				}
			}

		}

		if flag == "-v" || flag == "--verbose" {
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
		}

		// The code won't return anything about the responses unless a "match" flag is sent
		if flag == "--match" {
			global.Match = args[i+1]
			matchExists = true
		}

		if flag == "--no-color" {
			global.NoColor = true
		}

	}
	if !modeExists {
		logChan <- log.Entry{Text: "[!] Mode wasn't identified, using \"flood\" as default...", Verbosity: 1}
	}
	if !matchExists {
		logChan <- log.Entry{Text: "[!] Response match wasn't identified, increase verbosity to check for feedback", Verbosity: 1}
	}

}

/* The command MUST follow:

	gorace -u 'url' \
		   -H 'h1_name:h1_value,h2_name:h2_value'	\
		   -b 'c1_name = c1_value, c2_name= c2_value'	\
		   -d 'd1_name =d1_value' \
		   -w 'WORDLIST1=PATH1,WORDLIST2=PATH2' \

		   -u 'url2' \
		   -H 'h1_name:h1_value,h2_name:h2_value'	\
		   -b 'c1_name = c1_value, c2_name= c2_value'	\
		   -d 'd1_name =d1_value' \
		   -w 'WORDLIST1=PATH1,WORDLIST2=PATH2' \

After every "-u" except the first, the flags before it are all added to a []Website,

How parseCLI() works:

	urlFlag []string 					-> Each index represents a url
	flagMap := map[string]*[]string 	-> maps a Flag ("-u") with the address to the urlArgs string array
	So flagMap["-u"] corresponds to the address of the urlArgs string array
*/

// Take the abreviation -f of --flag, and turns it into --flag, because it makes dealing with the flags from initFlags()

func parse(args []string, log chan<- log.Entry) []Config {

	flags := initFlags()
	normalizeInputFlags(&args) // -f to --flag

	configs := getConfigs(flags, args, log) // May return files with Wordlists

	var websites []Config

	for _, c := range configs {

		if len(c.Wordlists) > 0 {
			w := handleWordlist(c.copy()) // If any wordlist was registered, all the headers, cookies and data placeholders registered before will be replaced
			websites = append(websites, w...)
		} else {
			websites = append(websites, c)
		}

	}

	return websites
}

func helpFlag(args []string) {
	for _, f := range args {
		if f == "-h" || f == "--help" {
			fmt.Println(assets.Help)
			os.Exit(0)
		}
	}
}

// Using args := os.Args[:2], in the loop, args[i] = flag, args[i+1] = parameter
func CLI(args []string, global *GlobalFlags, logChan chan log.Entry) []Config {

	helpFlag(args)
	readGlobalFlags(args, global, logChan)
	args = removeGlobalFlags(args)

	if len(args)%2 != 0 {
		panic("[x] A flag is missing a parameter! Exiting...")
	}

	// Flags that aren't related to the website request struct (Mode and Verbosity)

	// bar := strings.Repeat("⸺", 30)
	// logChan <- log.Entry{Text: bar + "\n", Verbosity: 0}

	configs := parse(args, logChan)
	logChan <- log.Entry{Text: "[+] Input read successfully!\n", Verbosity: 1}

	return configs

}
