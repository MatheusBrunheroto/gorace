package input

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Flag struct {
	name      string
	parameter []string
	exists    bool
}
type Field struct {
	flag      *Flag
	pairs     []Pair
	delimiter string
}

func newFlag(name string) Flag {
	var flag Flag
	flag.name = name
	flag.parameter = []string{}
	flag.exists = false
	return flag
}
func NewField(flag *Flag, pairs []Pair, delimiter string) Field {

	var field Field

	field.flag = flag
	field.pairs = pairs
	field.delimiter = delimiter

	return field
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
	flag.parameter = append(flag.parameter, arg)

	return nil
}

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
func fixEmpty(flag *Flag, urlAmount int, name string) error {

	if flag.exists == false {

		parameter := ""

		if strings.Contains(name, "--method") {
			parameter = "GET"
		}
		if strings.Contains(name, "--threads") {
			parameter = "1"
		}

		flag.parameter = append(flag.parameter, parameter)
		return nil
	}

	// Flag exists, may have more than one
	flagAmount := len(flag.parameter)
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

func parseCLI() ([]Website, error) {

	var websites []Website
	args := os.Args[1:]

	// Initializes all the flags
	urlFlag := newFlag("-u || --url")
	methodFlag := newFlag("-m || --method")
	headersFlag := newFlag("-H || --headers")
	cookiesFlag := newFlag("-c || --cookies")
	dataFlag := newFlag("-d || --data")
	wordlistsFlag := newFlag("-w || --wordlists")
	threadsFlag := newFlag("-t || --threads")

	flags := [7]*Flag{&urlFlag, &methodFlag, &headersFlag, &cookiesFlag, &dataFlag, &wordlistsFlag, &threadsFlag}
	flagMap := map[string]*Flag{
		"-u": &urlFlag, "--url": &urlFlag,
		"-X": &methodFlag, "--method": &methodFlag,
		"-H": &headersFlag, "--headers": &headersFlag,
		"-c": &cookiesFlag, "--cookies": &cookiesFlag,
		"-d": &dataFlag, "--data": &dataFlag,
		"-w": &wordlistsFlag, "--wordlist": &wordlistsFlag,
		"-t": &threadsFlag, "--threads": &threadsFlag,
	}

	// Read the arguments
	var urlAmount int
	for i := 0; i < len(args); i++ {

		flag, exist := flagMap[args[i]]
		if !exist {
			continue
		}
		// Starts to read Flags for new URL in case of double endpoint (ignores the first URL)
		urlAmount = len(urlFlag.parameter)
		if (args[i] == "-u" || args[i] == "--url") && urlAmount != 0 {
			for _, f := range flags[1:] { // Does not include urlArgs on pourpose
				fixEmpty(f, urlAmount, f.name) // checks for flag.Exists, if not, append empty
				f.exists = false
			}

		}

		if err := readFlag(flag, i, args); err != nil {
			return []Website{}, err
		}
		i++

		flag.exists = true
	}
	//RESUMIR ISSO TAMBEM POSSIVELMENTE
	for _, f := range flags[1:] { // Does not include urlArgs on pourpose
		fixEmpty(f, urlAmount, f.name) // checks for flag.Exists, if not, append empty
		f.exists = false
	}

	headers := NewField(&headersFlag, []Pair{}, ":")
	cookies := NewField(&cookiesFlag, []Pair{}, "=")
	data := NewField(&dataFlag, []Pair{}, "=")
	wordlists := NewField(&wordlistsFlag, []Pair{}, "=")
	fields := []*Field{&headers, &cookies, &data, &wordlists}

	for i := 0; i < urlAmount; i++ {

		if err := filterUrl(&urlFlag.parameter[i]); err != nil { // CUIDADO AQUI
			return []Website{}, err
		}
		filterMethod(&methodFlag.parameter[i]) // Doesn't need error to be returned, worst case scenario, GET is used

		threads, err := strconv.Atoi(threadsFlag.parameter[i])
		if err != nil {
			return []Website{}, err
		}

		// Parse keys into the headers, cookies, data and wordlists fields
		for _, field := range fields {
			if i >= len(field.flag.parameter) || field.flag.parameter[i] == "" {
				continue
			}
			if err := filterKeys(field.flag.parameter[i], &field.pairs, field.delimiter); err != nil {
				return []Website{}, err
			}
		}

		// If any wordlist was registered, all the headers, cookies and data placeholders registered before will be replaced
		if len(wordlists.pairs) > 0 {

			filteredHeaders, expandedHeaders, err := handleWordlist(headers.pairs, wordlists.pairs)
			if err != nil {
				return []Website{}, err
			}
			filteredCookies, expandedCookies, err := handleWordlist(cookies.pairs, wordlists.pairs)
			if err != nil {
				return []Website{}, err
			}
			filteredData, expandedData, err := handleWordlist(data.pairs, wordlists.pairs)
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
							Url:     urlFlag.parameter[i],
							Method:  methodFlag.parameter[i],
							Headers: newHeaders,
							Cookies: newCookies,
							Data:    newData,
							Threads: threads,
						}
						websites = append(websites, w)

					}
				}
			}
		} else {

			w := Website{
				Url:     urlFlag.parameter[i],
				Method:  methodFlag.parameter[i],
				Headers: headers.pairs,
				Cookies: cookies.pairs,
				Data:    data.pairs,
				Threads: threads,
			}
			websites = append(websites, w)

		}

	}

	return websites, nil
}

// Using args := os.Args[:2], in the loop, args[i] = flag, args[i+1] = parameter
func RunCLI() ([]Website, error) {

	websites, err := parseCLI()
	if err != nil {
		return []Website{}, err
	}
	fmt.Printf("[+] Input read successfully!\n\n")

	return websites, nil

}
