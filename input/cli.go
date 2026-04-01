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

func RunCLI(start chan<- struct{}, jobs chan<- Website, thread_amount *int) error {

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
	if len(wordlistPath) > 0 {

		filteredHeaders, expandedHeaders, err := handleWordlist(headers, wordlists)
		if err != nil {
			return err
		}
		filteredCookies, expandedCookies, err := handleWordlist(cookies, wordlists)
		if err != nil {
			return err
		}
		filteredData, expandedData, err := handleWordlist(data, wordlists)
		if err != nil {
			return err
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

					newHeaders := append(filteredHeaders, h)
					newCookies := append(filteredCookies, c)
					newData := append(filteredData, d)

					website := Website{
						Url:     *urlFlag,
						Method:  *methodFlag,
						Headers: newHeaders,
						Cookies: newCookies,
						Data:    newData,
					}

					fmt.Println(website)
					jobs <- website

				}
			}
		}
	} else {

		website := Website{
			Url:     *urlFlag,
			Method:  *methodFlag,
			Headers: headers,
			Cookies: cookies,
			Data:    data,
		}
		jobs <- website
	}

	// no need to deal with empty values
	close(start)
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
