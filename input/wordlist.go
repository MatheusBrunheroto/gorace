package input

import (
	"bufio"
	"errors"
	"os"
)

func readWordlists(path string) ([]string, error) {

	file, err := os.Open(path)
	if err != nil {
		return []string{}, errors.New("Invalid path -> " + path)
	}
	defer file.Close()

	var words []string

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		words = append(words, scanner.Text())
	}
	if err := scanner.Err(); err != nil { // ✅ agora detecta
		return []string{}, err
	}
	return words, nil
}

func parseWordlist(words []string, field []Pair, pair Pair, isKey bool) []Pair {

	var expanded []Pair

	for _, w := range words {

		if isKey {
			expanded = append(expanded, Pair{Key: w, Value: pair.Value})
		} else {
			expanded = append(expanded, Pair{Key: pair.Key, Value: w})
		}
	}
	return expanded
}

func removeWordlistPlaceholder(key string, field *[]Pair) {
	for i := range *field {
		if (*field)[i].Key == key {
			*field = append((*field)[:i], (*field)[i+1:]...)
			i--
		}
	}
}

func insertsWordlist(original Config, expanded map[string][]Pair) []Config {

	var newConfigs []Config
	if len(expanded["Headers"]) == 0 {
		expanded["Headers"] = []Pair{{}}
	}
	if len(expanded["Cookies"]) == 0 {
		expanded["Cookies"] = []Pair{{}}
	}
	if len(expanded["Data"]) == 0 {
		expanded["Data"] = []Pair{{}}
	}
	for _, h := range expanded["Headers"] {
		for _, c := range expanded["Cookies"] {
			for _, d := range expanded["Data"] {

				new := Config{
					Url:     original.Url,
					Method:  original.Method,
					Headers: append(original.Headers, h),
					Cookies: append(original.Cookies, c),
					Data:    append(original.Data, d),
					Threads: original.Threads,
					Delay:   original.Delay,
				}
				newConfigs = append(newConfigs, new)

			}
		}
	}

	return newConfigs
}

/*
-H 'key:value'

-H 'WORDLIST:value'
-H 'key:WORDLIST'

*/
// na verdade isso ta tudo errado, preciso salva separado

func handleWordlist(config Config) ([]Config, error) {

	fields := map[string]*[]Pair{
		"Headers": &config.Headers,
		"Cookies": &config.Cookies,
		"Data":    &config.Data,
	}
	expanded := map[string][]Pair{
		"Headers": []Pair{},
		"Cookies": []Pair{},
		"Data":    []Pair{},
	}

	for _, wordlist := range config.Wordlists {

		words, err := readWordlists(wordlist.Value) // Path
		if err != nil {
			return []Config{}, err
		}
		// Headers, Cookies and Data
		for name, field := range fields {

			// Field pairs, headers.key headers.value
			for _, pair := range *field {

				if pair.Key == wordlist.Key {
					removeWordlistPlaceholder(pair.Key, field)
					e := parseWordlist(words, *field, pair, true)
					expanded[name] = append(expanded[name], e...)
					continue
				}
				if pair.Value == wordlist.Key {
					removeWordlistPlaceholder(pair.Value, field)
					e := parseWordlist(words, *field, pair, true)
					expanded[name] = append(expanded[name], e...)
					continue
				}

			}

		}

	}
	new := insertsWordlist(config, expanded)

	return new, nil
}
