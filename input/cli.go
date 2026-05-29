package input

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

/*
Lets assume that flagAmount =

	"--url": 1,
	"--method": 1,
	"--headers": 0
	"--cookies": 0
	"--data": 1
	"--threads": 0
	"--wordlist": 0

Every quantity should be the same as the number of url, so the append doesn't create an undesired parameter on
the array, causing a "desync".
Example, if the user wants to send data only in the 3rd URL, the data array will have 2 empty elements before it.
*/
func fillDefault(flag *Flag, urlAmount int, name string) error {

	if flag.exists == false {

		parameter := ""

		if strings.Contains(name, "--method") {
			parameter = "GET"
		}
		if strings.Contains(name, "--threads") {
			parameter = "1"
		}
		if strings.Contains(name, "--delay") {
			parameter = "0"
		}

		flag.raw = append(flag.raw, parameter)
		return nil
	}

	// Flag exists, may have more than one
	flagAmount := len(flag.raw)
	if flagAmount > urlAmount {
		return errors.New("Two or more equal flags detected! -> " + flag.name)
	}

	return nil
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

func syncFlag(flags map[string]*Flag, syncPoint int) {
	for k, f := range flags {
		if k == "--url" {
			continue
		}
		fillDefault(f, syncPoint, f.name)
		f.exists = false
	}
}

func parseArgs(flags map[string]*Flag, args []string) error {

	nonUrlFlags := flags
	delete(nonUrlFlags, "--url")

	var urlAmount int
	for i := 0; i < len(args); i++ {

		flag, exist := flags[args[i]]
		if !exist {
			continue
		}
		// Starts to read Flags for new URL in case of double endpoint (ignores the first URL)
		urlAmount = len(flags["--url"].raw)
		if (args[i] == "--url") && urlAmount != 0 {
			syncFlag(flags, urlAmount)
		}

		if err := readFlag(flag, i, args); err != nil {
			return err
		}
		i++

		flag.exists = true
	}
	syncFlag(flags, urlAmount) // Last URL

	return nil
}

func parseCLI(args []string) ([]Website, error) {

	var websites []Website

	flags := initFlags()
	normalizeFlag(&args) // -f to --flag
	parseArgs(flags, args)

	// Read the arguments
	// readFlags()

	//RESUMIR ISSO TAMBEM POSSIVELMENTE

	fields := initFields(flags)
	for i := 0; i < len(flags["--url"].raw); i++ {

		// Reset, so website 1 pairs doesn't reflect in website 2
		fields["headers"].pairs = []Pair{}
		fields["cookies"].pairs = []Pair{}
		fields["data"].pairs = []Pair{}
		fields["wordlists"].pairs = []Pair{}

		if err := filterUrl(&flags["--url"].raw[i]); err != nil { // CUIDADO AQUI
			return []Website{}, err
		}
		filterMethod(&flags["--method"].raw[i]) // Doesn't need error to be returned, worst case scenario, GET is used

		threads, err := strconv.Atoi(flags["--threads"].raw[i])
		if err != nil {
			return []Website{}, err
		}
		delay, err := strconv.Atoi(flags["--delay"].raw[i])

		// Parse keys into the headers, cookies, data and wordlists fields
		for _, field := range fields {
			if i >= len(field.flag.raw) || field.flag.raw[i] == "" {
				continue
			}
			if err := filterKeys(field.flag.raw[i], &field.pairs, field.delimiter); err != nil {
				return []Website{}, err
			}
		}

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
func RunCLI(args []string) ([]Website, string, error) {

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

	websites, err := parseCLI(args)
	if err != nil {
		return []Website{}, "", err
	}
	fmt.Printf("[+] Input read successfully!\n\n")

	return websites, mode, nil

}
