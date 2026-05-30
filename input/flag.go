package input

import (
	"errors"
	"strings"
)

/*
Primeiro, cria cada site, nao importa o placeholder de WORDLIST, faz append num temporario
depois que fizer esse append e acabar de ler os args, desiste deles e trabalha direto na wordlist
a cada URL começa um novo, os defaults ja estao setados do newCOnfig.
E SE tiver 2 flags iguais? trato como mesmo site 2x, ou junto elas?
- URL1 -> H [1] -H [2] -> -H [1,2] <- se for array de pares, append simples, se for METODO Ou outra coisa, ERRO BRUTAL
- URL1 -> H[1] -URL1.1 -> -H[2] -C // se quiser isso poe wordlist la né amigao

*/
// returns raw arguments later
// OK

func readFlagValue(index int, args []string) (string, error) {

	if index >= len(args) {
		return "", errors.New("Missing parameter for flag -> ") // PASSAR ALGUMA STRING PRA CA PRA RETORNAR CERTINHO
	}
	arg := args[index]

	if strings.HasPrefix(arg, "-") {
		return "", errors.New("Wrong parameter usage! -> ")
	}

	return arg, nil
}

// OK
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

// OK
func initFlags() map[string]string {

	// flags := [8]*Flag{&urlFlag, &methodFlag, &headersFlag, &cookiesFlag, &dataFlag, &wordlistsFlag, &threadsFlag, &delayFlag}
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
