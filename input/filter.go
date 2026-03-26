package input

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
)

func filterKeys(raw string, key_map map[string]string) error {

	// This avoids string being empty, because it is mandatory to have ":"
	if !strings.Contains(raw, ":") {
		return errors.New("Invalid key! -> " + raw)
	}
	key := strings.SplitN(raw, ":", 2)
	if key[0] == "" || key[1] == "" {
		return errors.New("Invalid key! -> " + raw)
	}

	key[0] = strings.TrimSpace(key[0])
	key[1] = strings.TrimSpace(key[1])
	key_map[key[0]] = key[1] // key_map[key_name] = key_value
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
		strings.ReplaceAll(*method, " ", "")
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
