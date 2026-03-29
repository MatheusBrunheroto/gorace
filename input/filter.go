package input

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
)

func splitSymbol(raw string, symbol string) (string, string, error) {
	key := strings.SplitN(raw, symbol, 2)
	if len(key) != 2 || key[0] == "" || key[1] == "" {
		return "", "", errors.New("Invalid key! -> " + raw + "\nCheck examples with gorace --help.")
	}
	key[0] = strings.TrimSpace(key[0])
	key[1] = strings.TrimSpace(key[1])
	return key[0], key[1], nil
}
func filterKeys(raw string, key_map map[int]map[string]string, symbol string) (int, error) {

	// PROVAVELMENTE REMOVER TUDO ISSO AQUI
	// This also avoids string being empty
	if !strings.HasPrefix(raw, "{") || !strings.HasSuffix(raw, "}") {
		return -1, errors.New("Invalid key! -> " + raw + "\nCheck examples with gorace --help.")
	}
	raw = strings.Trim(raw, "{}") // Removes curly braces to avoid it invading a header
	// ATE AQUI, FICAR MAIS FACIL DE ESCREVER DEPOIS
	// If has multiple headers, splits ':'
	var keys []string
	if !strings.Contains(raw, ",") {

		key_name, key_value, err := splitSymbol(raw, symbol)
		if err != nil {
			return -1, err
		}
		if key_map[0] == nil {
			key_map[0] = make(map[string]string)
		}
		key_map[0][key_name] = key_value

	} else {

		keys = strings.Split(raw, ",")
		for i, k := range keys {

			key_name, key_value, err := splitSymbol(k, symbol)
			if err != nil {
				return -1, err
			}
			if key_map[i] == nil {
				key_map[i] = make(map[string]string)
			}
			key_map[i][key_name] = key_value

			if strings.Contains(key_name, "WORDLIST") || strings.Contains(key_value, "WORDLIST") {
				wordlist_size++
			}
		}

	}

	return wordlist_size, nil
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

func filterPath(path string) string {

	key := strings.SplitN(raw, "=", 2)

	if len(key) != 2 || key[0] == "" || key[1] == "" {
		return "", "", errors.New("Invalid key! -> " + raw + "\nCheck examples with gorace --help.")
	}
	key[0] = strings.TrimSpace(key[0])
	key[1] = strings.TrimSpace(key[1])
	return key[0], key[1], nil

}
