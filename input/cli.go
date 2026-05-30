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

func handleField(current Config, flag string, value string) error {

	var pairs []Pair
	var err error

	fields := map[string][]Pair{
		"--headers":   current.Headers,
		"--cookies":   current.Cookies,
		"--data":      current.Data,
		"--wordlists": current.Wordlists,
	}

	field, _ := fields[flag]
	if pairs, err = parsePairs(value); err != nil { // parse pairs do anything with : or =
		return err
	}
	field = append(field, pairs...)
	return nil

}

func parseCLI(args []string) ([]Config, error) {

	flags := initFlags()
	normalizeInputFlags(&args) // -f to --flag

	configs, err := getConfigs(flags, args) // Will return non-wordlists vs wordlists
	if err != nil {
		return []Config{}, err
	}

	////////////////////////////////////////

	for i := 0; i < len(flags["--url"].raw); i++ {

		// Reset, so website 1 pairs doesn't reflect in website 2
		// LEMBRAR OQ EU TAVA FAZENDO COM WORDLIST PELO GIT
		fields["wordlists"].pairs = []Pair{}

		// If any wordlist was registered, all the headers, cookies and data placeholders registered before will be replaced
		if len(fields["wordlists"].pairs) > 0 {

			filteredHeaders, expandedHeaders, err := handleWordlist(fields["headers"].pairs, fields["wordlists"].pairs)
			if err != nil {
				return []Website{}, err
			}
			filteredCookies, expandedCookies, err := handleWordlist(fields["cookies"].pairs, fields["wordlists"].pairs)
			if err != nil {
				return []Website{}, err
			}
			filteredData, expandedData, err := handleWordlist(fields["data"].pairs, fields["wordlists"].pairs)
			if err != nil {
				return []Website{}, err
			}

			// Avoid loop not starting
			if len(expandedHeaders) == 0 {
				expandedHeaders = []Pair{{}}
			}
			if len(expandedCookies) == 0 {
				expandedCookies = []Pair{{}}
			}
			if len(expandedData) == 0 {
				expandedData = []Pair{{}}
			}
			for _, h := range expandedHeaders {
				for _, c := range expandedCookies {
					for _, d := range expandedData {

						// Adding filteredKey before avoids newField being empty in case k = []
						newHeaders := filteredHeaders
						newCookies := filteredCookies
						newData := filteredData

						if h.Key != "" && h.Value != "" {
							newHeaders = append(newHeaders, h)
						}
						if c.Key != "" && c.Value != "" {
							newCookies = append(newCookies, c)
						}
						if d.Key != "" && d.Value != "" {
							newData = append(newData, d)
						}

						w := Website{
							Url:     flags["--url"].raw[i],    // Unique
							Method:  flags["--method"].raw[i], // Unique
							Headers: newHeaders,
							Cookies: newCookies,
							Data:    newData,
							Threads: threads,
							Delay:   delay,
						}
						websites = append(websites, w)

					}
				}
			}
		} else {

			w := Website{
				Url:     flags["--url"].raw[i],
				Method:  flags["--method"].raw[i],
				Headers: fields["headers"].pairs,
				Cookies: fields["cookies"].pairs,
				Data:    fields["data"].pairs,
				Threads: threads,
				Delay:   delay,
			}
			websites = append(websites, w)
		}

	}

	return websites, nil
}

// Using args := os.Args[:2], in the loop, args[i] = flag, args[i+1] = parameter
func RunCLI(args []string) ([]Config, string, error) {

	var mode string = "flood"
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

	configs, err := parseCLI(args)
	if err != nil {
		return []Config{}, "", err
	}
	fmt.Printf("[+] Input read successfully!\n\n")

	return configs, mode, nil

}
