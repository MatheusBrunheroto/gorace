package input

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"
)

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

	link := flag.String("U", "", "Website URL")
	method := flag.String("X", "", "HTTP Method")

	flag.Parse()

	if *link == "" {
		return errors.New("No Website URL was informed (-U or --url)")
	}
	if *method == "" {
		fmt.Println("No method informed (-X or --method), \"GET\" will be used...")
		*method = "GET"
	}

	headers_map := make(map[string]string)
	cookies_map := make(map[string]string)
	data_map := make(map[string]string)

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
	for _, d := range data {
		err := filterKeys(d, data_map)
		if err != nil {
			fmt.Println(err)
			return err
		}
	}

	website := Website{
		Link:    *link,
		Method:  *method,
		Headers: headers_map,
		Cookies: cookies_map,
		Data:    data_map,
	}
	fmt.Println(*link)
	fmt.Println(*method)
	fmt.Println(headers)
	fmt.Println(cookies)
	fmt.Println(data)
	*websites = append(*websites, website)
	return nil

}
