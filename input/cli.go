package input

import (
	"flag"
	"fmt"
	"os"
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

func RunCLI(websites *[]Website) error {

	os.Args = append(os.Args[:1], os.Args[2:]...) // Removes the first argument (-c, -i, -w)

	var headers, cookies, data multiple_flags
	flag.Var(&headers, "H", "HTTP Headers")
	flag.Var(&cookies, "C", "Cookies")
	flag.Var(&data, "D", "POST Data")
	flag.Var(&data, "no-filter", "Doesn't filter the DATA field")

	url := flag.String("U", "", "Website URL")
	method := flag.String("X", "", "HTTP Method")

	flag.Parse()

	if err := filterUrl(url); err != nil {
		return err
	}
	filterMethod(method) // Doesn't need error to be returned, worst case scenario, GET is used

	headers_map := make(map[string]string)
	cookies_map := make(map[string]string)
	data_map := make(map[string]string)
	var unfiltered_data bool

	for _, h := range headers {
		err := filterKeys(h, headers_map)
		if err != nil {
			fmt.Println(err)
			return err
		}
	}
	for _, c := range cookies {
		err := filterKeys(c, cookies_map)
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
