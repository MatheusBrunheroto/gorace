package input

import (
	"strings"
)

func readFlagValue(index int, args []string) string {

	if index >= len(args) {
		panic("Missing parameter for flag -> ") // PASSAR ALGUMA STRING PRA CA PRA RETORNAR CERTINHO
	}
	arg := args[index]

	if strings.HasPrefix(arg, "-") {
		panic("Wrong parameter usage! -> ")
	}

	return arg
}

// Take the abreviation -f of --flag, and turns it into --flag, because it makes dealing with the flags from initFlags()
func normalizeInputFlags(args *[]string) {

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

func initFlags() map[string]string {

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
