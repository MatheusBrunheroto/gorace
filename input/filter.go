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
func splitSymbol(raw string, symbol string) (string, string, error) {
	key := strings.SplitN(raw, symbol, 2)
	if len(key) != 2 || key[0] == "" || key[1] == "" {
		return "", "", errors.New("Invalid key! -> " + raw + "\nCheck examples with gorace --help.")
	}
	key[0] = strings.TrimSpace(key[0])
	key[1] = strings.TrimSpace(key[1])
	return key[0], key[1], nil
}
func filterKeys(raw string, entry *[]Pair, delimiter string) error {

	// If has multiple headers, splits ':'
	// por funcao aq
	if !strings.Contains(raw, ",") {

		k, v, err := splitSymbol(raw, delimiter)
		if err != nil {
			return err
		}
		*entry = append(*entry, Pair{Key: k, Value: v})

	} else {

		var pairs []string
		pairs = strings.Split(raw, ",")

		for _, pair := range pairs {

			k, v, err := splitSymbol(pair, delimiter)
			if err != nil {
				return err
			}
			*entry = append(*entry, Pair{Key: k, Value: v})

		}

	}

	return nil
}

func filterUrl(target *string) error {

	if *target == "" {
		return errors.New("No Website URL was informed (-U or --url)")
	}
	if strings.HasPrefix(*target, "-") {
		return errors.New("Invalid URL for -U or --url -> " + *target)
	}

	if !strings.HasPrefix(*target, "http://") && !strings.HasPrefix(*target, "https://") {
		*target = "https://" + *target
		fmt.Println("URL must start with http:// or https:// -> New url: " + *target)
	}

	u, err := url.Parse(*target)
	if err != nil || u.Host == "" {
		return errors.New("invalid URL -> " + *target)
	}

	return nil
}

func filterMethod(method *string) {

	methods := []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS", "TRACE", "CONNECT"}

	if *method == "" {
		fmt.Println("No method informed (-X or --method), \"GET\" will be used...")
		*method = "GET"
		return
	}

	*method = strings.ToUpper(*method)

	if strings.Contains(*method, " ") {
		*method = strings.ReplaceAll(*method, " ", "")
		fmt.Println("Method contains a SPACE character, removing...")
	}

	for _, m := range methods {
		if *method == m {
			return
		}
	}
	fmt.Printf("Method \"%s\" not recognized within [GET, POST, PUT, PATCH, DELETE, HEAD, OPTIONS, TRACE, CONNECT], proceeding anyways...\n", *method)

	// COLOCAR REGEX
}
