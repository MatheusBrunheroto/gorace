package input

import (
	"flag"
	"fmt"
	"strings"
)

// IF diverse inputs are provided, creates only one "website", but for the wordlist, change the

type multiple_flags []string

// Required to use "flag.Var(...)"
func (f *multiple_flags) String() string {
	return strings.Join(*f, ", ")
}
func (f *multiple_flags) Set(value string) error {
	*f = append(*f, value)
	return nil
}

func RunCLI(websites *[]Website, thread_amount *int) error {

	threadsFlag := flag.Int("t", 50, "Amount of Threads")
	urlFlag := flag.String("u", "", "Website URL")
	methodFlag := flag.String("X", "", "HTTP Method")
	headersFlag := flag.String("H", "", "HTTP Headers")
	cookiesFlag := flag.String("c", "", "Cookies")
	dataFlag := flag.String("d", "", "POST Data")

	var wordlistPath multiple_flags
	flag.Var(&wordlistPath, "w", "Wordlists ('NAME=path')")

	flag.Parse() // flag.String(&data, "no-filter", "Doesn't filter the DATA field")

	if *threadsFlag <= 0 {
		return fmt.Errorf("Invalid amount of threads -> %d", *threadsFlag)
	}
	*thread_amount = *threadsFlag

	if err := filterUrl(urlFlag); err != nil { // CUIDADO AQUI
		return err
	}
	filterMethod(methodFlag) // Doesn't need error to be returned, worst case scenario, GET is used

	/* By reading the wordlists first, it's possible to separate the key_names from the path, and
	   check later if a header/cookie/data contains this name as a key_name/key_value */

	// key_map[n][key_name] = key_value, this allows repeated key_names
	var headers, cookies, data, wordlists []KeyValue

	if err := filterKeys(*headersFlag, &headers, ":"); err != nil {
		return err
	}
	if err := filterKeys(*cookiesFlag, &cookies, "="); err != nil {
		return err
	}
	if err := filterKeys(*dataFlag, &data, "="); err != nil {
		return err
	}
	for _, w := range wordlistPath {
		if err := filterKeys(w, &wordlists, "="); err != nil {
			return err
		}
	}

	// Now with all headers, cookies and data, search for the wordlists placeholders registered before
	// loop for aqui, pra iterar sobre todos os map

	if len(wordlistPath) > 0 {
		var err error

		if headers, err = handleWordlist(headers, wordlists); err != nil {
			return err
		}
		if cookies, err = handleWordlist(cookies, wordlists); err != nil {
			return err
		}
		if data, err = handleWordlist(data, wordlists); err != nil {
			return err
		}
	}

	website := Website{
		Url:     *urlFlag,
		Method:  *methodFlag,
		Headers: headers,
		Cookies: cookies,
		Data:    data,
	}

	*websites = append(*websites, website)
	return nil

}
