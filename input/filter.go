package input

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
)

/*
The two following functions limit which inputs are acceptable. Examples of correct inputs:

	-H --headers  -> 'h1_name:h1_value,h2_name:h2_value'
	-b --cookies  -> 'c1_name = c1_value, c2_name= c2_value'
	-d --data 	  -> 'd1_name =d1_value'
	-w --wordlist -> 'WORDLIST1=PATH1,WORDLIST2=PATH2'

Anything that goes against the structure 'KEY=VALUE' for -b, -d and -w;
or 'KEY:VALUE' for -H is rejected.
*/

func findDelimiter(s string, values []string) (string, bool) {
	for _, v := range values {
		if strings.Contains(s, v) {
			return v, true
		}
	}
	return "", false
}

func parseKeyValue(raw string, pairs *[]Pair) error {

	if delimiter, found := findDelimiter(raw, []string{":", "="}); found {
		parts := strings.SplitN(raw, delimiter, 2)

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		if key == "" || value == "" {
			return errors.New("[-] Invalid key! -> " + raw + "\nCheck examples with gorace --help.")
		}

		*pairs = append(*pairs, Pair{Key: key, Value: value})
		return nil
	}

	return errors.New("[-] No valid delimiter found in: " + raw)

}

func parsePairs(raw string) ([]Pair, error) {

	var parsedPairs []Pair

	// If contais ',' or '&' (has multiple keys), splits ':' or '='
	if delimeter, found := findDelimiter(raw, []string{",", "&"}); found {

		pairs := strings.Split(raw, delimeter)

		for _, pair := range pairs {

			if err := parseKeyValue(pair, &parsedPairs); err != nil {
				return []Pair{}, err
			}

		}
	} else {
		if err := parseKeyValue(raw, &parsedPairs); err != nil {
			return []Pair{}, err
		}
	}

	return parsedPairs, nil
}

func normalizeUrl(target *string) error {

	if *target == "" {
		return errors.New("[-] No Website URL was informed (-U or --url)")
	}
	if strings.HasPrefix(*target, "-") {
		return errors.New("[-] Invalid URL for -U or --url -> " + *target)
	}

	if !strings.HasPrefix(*target, "http://") && !strings.HasPrefix(*target, "https://") {
		*target = "https://" + *target
		fmt.Println("[!] URL must start with http:// or https:// -> New url: " + *target)
	}

	u, err := url.Parse(*target)
	if err != nil || u.Host == "" {
		return errors.New("[-] Invalid URL -> " + *target)
	}

	return nil
}

func normalizeMethod(method string) string {

	methods := []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS", "TRACE", "CONNECT"}

	if method == "" {
		fmt.Println("[!] No method informed (-X or --method), \"GET\" will be used...")
		return "GET"
	}

	method = strings.ToUpper(method)

	if strings.Contains(method, " ") {
		fmt.Println("[!] Inserted method contains a SPACE character, removing...")
		method = strings.ReplaceAll(method, " ", "")
	}

	for _, m := range methods {
		if method == m {
			return method
		}
	}

	fmt.Printf("[!] Method \"%s\" not recognized within [GET, POST, PUT, PATCH, DELETE, HEAD, OPTIONS, TRACE, CONNECT], proceeding anyways...\n", method)
	return "GET"
	// COLOCAR REGEX
}
