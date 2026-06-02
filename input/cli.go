package input

import (
	"fmt"
)

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

func parseCLI(args []string) []Config {

	flags := initFlags()
	normalizeInputFlags(&args) // -f to --flag

	configs := getConfigs(flags, args) // May return files with Wordlists

	var websites []Config

	for _, c := range configs {

		if len(c.Wordlists) > 0 {
			w := handleWordlist(c.copy())     // If any wordlist was registered, all the headers, cookies and data placeholders registered before will be replaced
			websites = append(websites, w...) // pode retornar mais de um
		} else {
			websites = append(websites, c)
		}

	}

	return websites
}

////////////////////////////////////////

// Using args := os.Args[:2], in the loop, args[i] = flag, args[i+1] = parameter
func RunCLI(args []string) ([]Config, string) {

	if len(args)%2 != 0 {
		panic("[!] A flag is missing a parameter! Exiting...")
	}

	var mode string = "flood" // Default
	modes := []string{"sequential", "round-sequential", "cascade", "round-cascade", "flood"}
	var modeExists bool = false

	for i, flag := range args {

		if flag == "-m" || flag == "--mode" {

			verifyMode := args[i+1]
			for _, m := range modes {

				if verifyMode == m {
					modeExists = true
					mode = verifyMode
					break
				}

			}

		}
	}
	if !modeExists {
		fmt.Println("[!] Mode wasn't identified, using \"flood\" as default...")
	}

	configs := parseCLI(args)
	fmt.Printf("[+] Input read successfully!\n\n")

	return configs, mode

}
