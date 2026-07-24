package input

import (
	"gorace/log"
)

func parse(args []string, log chan<- log.Entry) []Config {

	configs := getConfigs(initConfigFlags(), args, log) // May return files with Wordlists

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

// Using args := os.Args[:2], in the loop, args[i] = flag, args[i+1] = parameter
func CLI(args []string, global *GlobalFlags, logChan chan log.Entry) []Config {

	normalizeInputFlags(&args) // Converts -f to --flag, to simplify comparsions

	checkFlagsConsistency(args)

	readGlobalFlags(args, global, logChan) // Flags that aren't related to the website request struct

	configs := parse(args, logChan)

	logChan <- log.Entry{Text: "[+] Input read successfully!\n", Verbosity: 1}

	return configs

}
