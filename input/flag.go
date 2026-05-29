package input

import (
	"errors"
	"strings"
)

type Flag struct {
	name   string
	raw    []string
	exists bool
}

func newFlag(name string) Flag {
	var flag Flag
	flag.name = name
	flag.raw = []string{}
	flag.exists = false
	return flag
}

func readFlag(flag *Flag, index int, args []string) error {

	index++
	if index >= len(args) {
		return errors.New("Missing parameter for flag -> ") // PASSAR ALGUMA STRING PRA CA PRA RETORNAR CERTINHO
	}
	arg := args[index]

	if strings.HasPrefix(arg, "-") {
		return errors.New("Wrong parameter usage! -> ")
	}
	flag.raw = append(flag.raw, arg)

	return nil
}

// Take the abreviation -f of --flag, and turns it into --flag, because it makes dealing with the flags from initFlags()
func normalizeFlag(args *[]string) {

	table := map[string]string{
		"-u": "--url",
		"-X": "--method",
		"-H": "--headers",
		"-b": "--cookies",
		"-d": "--data",
		"-w": "--wordlist",
		"-t": "--threads",
		"-D": "--delay",
	}

	for i := 0; i < len(*args); i++ {
		if normalized, ok := table[(*args)[i]]; ok {
			(*args)[i] = normalized
		}
	}

}

func initFlags() map[string]*Flag {

	urlFlag := newFlag("--url")
	methodFlag := newFlag("--method")
	headersFlag := newFlag("--headers")
	cookiesFlag := newFlag("--cookies")
	dataFlag := newFlag("--data")
	wordlistsFlag := newFlag("--wordlist")
	threadsFlag := newFlag("--threads")
	delayFlag := newFlag("--delay")

	// flags := [8]*Flag{&urlFlag, &methodFlag, &headersFlag, &cookiesFlag, &dataFlag, &wordlistsFlag, &threadsFlag, &delayFlag}
	return map[string]*Flag{
		"--url":      &urlFlag,       // -u
		"--method":   &methodFlag,    // -X
		"--headers":  &headersFlag,   // -H
		"--cookies":  &cookiesFlag,   // -b
		"--data":     &dataFlag,      // -d
		"--wordlist": &wordlistsFlag, // -w
		"--threads":  &threadsFlag,   // -t
		"--delay":    &delayFlag,     // -D
	}

}
