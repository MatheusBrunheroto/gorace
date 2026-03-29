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

	var wordlist_path multiple_flags

	var total_wordlists int = 0
	flag.Var(&wordlist_path, "w", "Wordlists ('NAME=path')")

	url := flag.String("u", "", "Website URL")
	method := flag.String("X", "", "HTTP Method")
	headers := flag.String("H", "", "HTTP Headers")
	cookies := flag.String("b", "", "Cookies")
	data := flag.String("d", "", "POST Data")
	threads := flag.Int("t", 50, "POST Data")
	flag.Parse()
	// flag.String(&data, "no-filter", "Doesn't filter the DATA field")

	if *threads <= 0 {
		return fmt.Errorf("Invalid amount of threads -> %d", *threads)
	}
	*thread_amount = *threads

	/* By reading the wordlists first, it's possible to separate the key_names from the path, and
	   check later if a header/cookie/data contains this name as a key_name/key_value */
	wordlists := make(map[string]string)
	if len(wordlist_path) > 0 {
		readWordlists(wordlist_path, wordlists)
	}

	if err := filterUrl(url); err != nil {
		return err
	}
	filterMethod(method) // Doesn't need error to be returned, worst case scenario, GET is used

	headers_map := make(map[int]map[string]string)
	cookies_map := make(map[int]map[string]string)
	data_map := make(map[int]map[string]string)
	var unfiltered_data bool

	for _, h := range *headers {

		wordlist_size, err := filterKeys(h, headers_map, ":")
		total_wordlists += wordlist_size
		for i := 0; i < wordlist_size; i++ {
			insertWordlist(headers_map, wordlists[i])
		}

		if err != nil {
			return err
		}
	}
	for _, c := range cookies {
		err := filterKeys(c, cookies_map, "=")
		if err != nil {
			fmt.Println(err)
			return err
		}
	}

	if unfiltered_data {
		data_map["unfiltered"] = data[0]
	} else {
		for _, d := range data {
			err := filterKeys(d, data_map)
			if err != nil {
				fmt.Println(err)
				return err
			}
		}
	}

	if total_wordlists != len(wordlists) {
		return fmt.Errorf("number of WORDLIST placeholders (%d) does not match provided wordlists (%d)", total_wordlists, len(wordlists))
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
