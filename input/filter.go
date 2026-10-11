package input

import (
	"gorace/log"
	"net/url"
	"strings"
)

/*
The two following functions limit which inputs are acceptable. Examples of correct inputs:

	-H --headers  -> 'h1_name:h1_value' (repeat -H for multiple headers)
	-b --cookies  -> 'c1_name=c1_value;c2_name=c2_value'
	-d --data 	  -> 'd1_name=d1_value&d2_name=d2_value' (URL-encoded)
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

func parseKeyValue(raw string, delimiter string) Pair {

	_, found := findDelimiter(raw, []string{delimiter})
	if !found {
		panic("[x] No valid delimiter found in: " + raw)
	}

	parts := strings.SplitN(raw, delimiter, 2)

	key := strings.TrimSpace(parts[0])
	value := strings.TrimSpace(parts[1])

	if key == "" || value == "" {
		panic("[x] Empty key or value! -> " + raw + "\nCheck examples with gorace --help")
	}

	return Pair{Key: key, Value: value}
}

func parsePairs(raw string, flag string) []Pair {

	var parsedPairs []Pair
	delimiter, separator := "=", ""
	switch flag {
	case "--headers":
		delimiter = ":"
	case "--cookies":
		separator = ";"
	case "--data":
		separator = "&"
	case "--wordlist":
		separator = ","
	}

	if separator != "" {

		for pair := range strings.SplitSeq(raw, separator) {
			parsedPairs = append(parsedPairs, parseKeyValue(pair, delimiter))
		}

	} else {
		parsedPairs = append(parsedPairs, parseKeyValue(raw, delimiter))
	}

	if flag == "--data" {
		for i, pair := range parsedPairs {
			key, keyErr := url.QueryUnescape(pair.Key)
			value, valueErr := url.QueryUnescape(pair.Value)
			if keyErr != nil || valueErr != nil {
				panic("[x] Invalid URL-encoded data -> " + raw)
			}
			parsedPairs[i] = Pair{Key: key, Value: value}
		}
	}

	return parsedPairs
}

//////////////////////////////////////
//////////////////////////////////////

func normalizeUrl(target *string, logChan chan<- log.Entry) {

	if *target == "" {
		panic("[x] No Website URL was informed (-U or --url)")
	}
	if strings.HasPrefix(*target, "-") {
		panic("[x] Invalid URL for -u or --url -> " + *target)
	}

	if !strings.HasPrefix(*target, "http://") && !strings.HasPrefix(*target, "https://") {
		*target = "https://" + *target
		logChan <- log.Entry{Text: "[!] URL must start with http:// or https:// -> New url: " + *target, Verbosity: 1}
	}

	u, err := url.Parse(*target)
	if err != nil || u.Host == "" {
		panic("[x] Invalid URL -> " + *target)
	}

}

func normalizeMethod(method string, logChan chan<- log.Entry) string {

	methods := []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS", "TRACE", "CONNECT"}

	if method == "" {
		logChan <- log.Entry{Text: "[!] No method informed (-X or --method), \"GET\" will be used...", Verbosity: 1}
		return "GET"
	}

	method = strings.ToUpper(method)

	if strings.Contains(method, " ") {
		logChan <- log.Entry{Text: "[!] Inserted method contains a SPACE character, removing...", Verbosity: 1}
		method = strings.ReplaceAll(method, " ", "")
	}

	for _, m := range methods {
		if method == m {
			return method
		}
	}
	logChan <- log.Entry{Text: "[!] Method \"%s\" not recognized within [GET, POST, PUT, PATCH, DELETE, HEAD, OPTIONS, TRACE, CONNECT], proceeding anyways...\n", Verbosity: 1}
	return "GET"

}
