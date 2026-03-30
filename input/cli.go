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

	threads := flag.Int("t", 50, "Amount of Threads")
	url := flag.String("u", "", "Website URL")
	method := flag.String("X", "", "HTTP Method")
	headers := flag.String("H", "", "HTTP Headers")
	cookies := flag.String("b", "", "Cookies")
	data := flag.String("d", "", "POST Data")

	var wordlist_path multiple_flags
	flag.Var(&wordlist_path, "w", "Wordlists ('NAME=path')")

	flag.Parse() // flag.String(&data, "no-filter", "Doesn't filter the DATA field")

	if *threads <= 0 {
		return fmt.Errorf("Invalid amount of threads -> %d", *threads)
	}
	*thread_amount = *threads

	if err := filterUrl(url); err != nil {
		return err
	}
	filterMethod(method) // Doesn't need error to be returned, worst case scenario, GET is used

	/* By reading the wordlists first, it's possible to separate the key_names from the path, and
	   check later if a header/cookie/data contains this name as a key_name/key_value */

	// key_map[n][key_name] = key_value, this allows repeated key_names
	headers_map := make(map[string]string)
	cookies_map := make(map[string]string)
	data_map := make(map[string]string)
	wordlists_map := make(map[string]string)

	if err := filterKeys(*headers, headers_map, ":"); err != nil {
		return err
	}
	if err := filterKeys(*cookies, cookies_map, "="); err != nil {
		return err
	}
	if err := filterKeys(*data, data_map, "="); err != nil {
		return err
	}
	for _, w := range wordlist_path {
		if err := filterKeys(w, wordlists_map, "="); err != nil {
			return err
		}
	}

	// Now with all headers, cookies and data, search for the wordlists placeholders registered before
	// loop for aqui, pra iterar sobre todos os map
	new_headers_map := make(map[int]map[string]string)
	new_cookies_map := make(map[int]map[string]string)
	new_data_map := make(map[int]map[string]string)
	if len(wordlist_path) > 0 {
		if err := handleWordlist(headers_map, new_headers_map, wordlists_map); err != nil {
			return err
		}
		if err := handleWordlist(cookies_map, new_cookies_map, wordlists_map); err != nil {
			return err
		}
		if err := handleWordlist(data_map, new_data_map, wordlists_map); err != nil {
			return err
		}
	}

	website := Website{
		Url:     *url,
		Method:  *method,
		Headers: headers_map,
		Cookies: cookies_map,
		Data:    data_map,
	}

	*websites = append(*websites, website)
	return nil

}
