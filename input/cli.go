package input

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

func readFlag(flag *[]string, index *int, args []string) error {

	(*index)++
	if *index >= len(args) {
		return errors.New("Missing parameter for flag -> ") // PASSAR ALGUMA STRING PRA CA PRA RETORNAR CERTINHO
	}
	arg := args[*index]

	if strings.HasPrefix(arg, "-") {
		return errors.New("Wrong parameter usage! -> ")
	}
	*flag = append(*flag, arg)

	return nil
}

/*
Lets suppose that flagAmount =

	"--url": 0,
	"--method": 0,
	"--headers": 0
	"--cookies": 0
	"--data": 0
	"--threads": 0
	"--wordlist": 0

Every quantity should be the same as the number of url, so the append doesn't create an undesired parameter on
the array, causing a "desync".
Example, if the user wants to send data only in the 3rd URL, the data array will have 2 empty elements before it.
*/
func fixEmpty(flag *[]string, url *[]string) error {

	flagAmount := len(*flag)
	urlAmount := len(*url)

	if flagAmount < urlAmount {
		*flag = append(*flag, "")
	} else if flagAmount > urlAmount {
		return errors.New("Two or more equal flags detected!")
	}

	return nil
}

/*
The command MUST follow:

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

After every but the first "-u", the flags before it are all added to a []Website,

How parseCLI() works:

	urlArgs []string 					-> Each index represents a url
	flagMap := map[string]*[]string 	-> maps a Flag ("-u") with the address to the urlArgs string array
	So flagMap["-u"] corresponds to the address of the urlArgs string array
*/
func parseCLI() ([]Website, error) {

	var websites []Website

	var (
		urlArgs      []string
		methodArgs   []string
		headersArgs  []string
		cookiesArgs  []string
		dataArgs     []string
		threadsArgs  []string
		wordlistArgs []string
	)
	parameters := [6]*[]string{&methodArgs, &headersArgs, &cookiesArgs, &dataArgs, &threadsArgs, &wordlistArgs}
	flagMap := map[string]*[]string{
		"-u": &urlArgs, "--url": &urlArgs,
		"-X": &methodArgs, "--method": &methodArgs,
		"-H": &headersArgs, "--headers": &headersArgs,
		"-c": &cookiesArgs, "--cookies": &cookiesArgs,
		"-d": &dataArgs, "--data": &dataArgs,
		"-t": &threadsArgs, "--threads": &threadsArgs,
		"-w": &wordlistArgs, "--wordlist": &wordlistArgs,
	}
	args := os.Args[1:]

	for i := 0; i < len(args); i++ {

		if flagPtr, ok := flagMap[args[i]]; ok {

			if (args[i] == "-u" || args[i] == "--url") && len(urlArgs) != 0 {
				for _, n := range parameters {
					fixEmpty(n, &urlArgs)
				}
			}

			if err := readFlag(flagPtr, &i, args); err != nil {
				return []Website{}, err
			}

		}

	}

	for i := 0; i < len(urlArgs); i++ {

		var headers, cookies, data, wordlists []KeyValue

		if err := filterUrl(&urlArgs[i]); err != nil { // CUIDADO AQUI
			return []Website{}, err
		}
		filterMethod(&methodArgs[i]) // Doesn't need error to be returned, worst case scenario, GET is used

		/* By reading the wordlists first, it's possible to separate the key_names from the path, and
		   check later if a header/cookie/data contains this name as a key_name/key_value */

		// key_map[n][key_name] = key_value, this allows repeated key_names
		if i < len(headersArgs) && headersArgs[i] != "" {
			if err := filterKeys(headersArgs[i], &headers, ":"); err != nil {
				return []Website{}, err
			}
		}
		if i < len(cookiesArgs) && cookiesArgs[i] != "" {
			if err := filterKeys(cookiesArgs[i], &cookies, "="); err != nil {
				return []Website{}, err
			}
		}
		if i < len(dataArgs) && dataArgs[i] != "" {
			if err := filterKeys(dataArgs[i], &data, "="); err != nil {
				return []Website{}, err
			}
		}
		if i < len(wordlistArgs) && wordlistArgs[i] != "" {
			if err := filterKeys(wordlistArgs[i], &wordlists, "="); err != nil {
				return []Website{}, err
			}
		}

		// Now with all headers, cookies and data, search for the wordlists placeholders registered before
		if len(wordlists) > 0 {

			filteredHeaders, expandedHeaders, err := handleWordlist(headers, wordlists)
			if err != nil {
				return []Website{}, err
			}
			filteredCookies, expandedCookies, err := handleWordlist(cookies, wordlists)
			if err != nil {
				return []Website{}, err
			}
			filteredData, expandedData, err := handleWordlist(data, wordlists)
			if err != nil {
				return []Website{}, err
			}

			// avoid loop stoping for no reason
			if len(expandedHeaders) == 0 {
				expandedHeaders = []KeyValue{{}}
			}
			if len(expandedCookies) == 0 {
				expandedCookies = []KeyValue{{}}

			}
			if len(expandedData) == 0 {
				expandedData = []KeyValue{{}}
			}
			for _, h := range expandedHeaders {
				for _, c := range expandedCookies {
					for _, d := range expandedData {

						// Adding filteredKey before avoids newKey being empty in case k = []
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
							Url:     urlArgs[i],
							Method:  methodArgs[i],
							Headers: newHeaders,
							Cookies: newCookies,
							Data:    newData,
						}
						websites = append(websites, w)

					}
				}
			}
		} else {

			w := Website{
				Url:     urlArgs[i],
				Method:  methodArgs[i],
				Headers: headers,
				Cookies: cookies,
				Data:    data,
			}
			websites = append(websites, w)

		}

	}
	return websites, nil
}

// Using args := os.Args[:2], in the loop, args[i] = flag, args[i+1] = argument
func RunCLI(start chan<- struct{}, jobs chan<- Website, thread_amount *int) error {

	websites, err := parseCLI()
	if err != nil {
		return err
	}

	//fmt.Println(websites)
	for _, v := range websites {
		fmt.Println(v)
	}
	return nil

}

// TENHO QUE CRIAR O SITE COM OS DADOS DA COISA
/*

Website
- data = {
			{key:value}
			{key:value}
			{key:WORDLIST}
		 }

if the wordlist has 3 words, only the last part should repeat



*/
