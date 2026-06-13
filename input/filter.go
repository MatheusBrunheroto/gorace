package input

import (
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

// OK
func findDelimiter(s string, values []string) (string, bool) {
	for _, v := range values {
		if strings.Contains(s, v) {
			return v, true
		}
	}
	return "", false
}

// OK
func parseKeyValue(raw string) Pair {

	delimiter, found := findDelimiter(raw, []string{":", "="})
	if !found {
		panic("[x] No valid delimiter found in: " + raw)
	}

	parts := strings.SplitN(raw, delimiter, 2)

	key := strings.TrimSpace(parts[0])
	value := strings.TrimSpace(parts[1])

	if key == "" || value == "" {
		panic("[x] Empty key or value! -> " + raw + "\nCheck examples with gorace --help.")
	}

	return Pair{Key: key, Value: value}
}

// OK
func parsePairs(raw string) []Pair {

	var parsedPairs []Pair

	// If contais ',' or '&' (has multiple keys), splits by ':' or '='
	if delimeter, found := findDelimiter(raw, []string{",", "&"}); found {

		// pairs := strings.SplitSeq(raw, delimeter)
		for pair := range strings.SplitSeq(raw, delimeter) {
			parsedPairs = append(parsedPairs, parseKeyValue(pair))
		}

	} else {
		parsedPairs = append(parsedPairs, parseKeyValue(raw))
	}

	return parsedPairs
}

//////////////////////////////////////
//////////////////////////////////////

func normalizeUrl(target *string) {

	if *target == "" {
		panic("[x] No Website URL was informed (-U or --url)")
	}
	if strings.HasPrefix(*target, "-") {
		panic("[x] Invalid URL for -u or --url -> " + *target)
	}

	if !strings.HasPrefix(*target, "http://") && !strings.HasPrefix(*target, "https://") {
		*target = "https://" + *target
		fmt.Println("[!] URL must start with http:// or https:// -> New url: " + *target)
	}

	u, err := url.Parse(*target)
	if err != nil || u.Host == "" {
		panic("[x] Invalid URL -> " + *target)
	}

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

}
