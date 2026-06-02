package input

import (
	"bufio"
	"fmt"
	"os"
)

// OK
func removeWordlistPlaceholder(key string, field *[]Pair) {
	var filtered []Pair
	for _, p := range *field {
		if p.Key != key {
			filtered = append(filtered, p)
		}
	}
	*field = filtered
}

// ok
func readWordlists(path string) []string {

	file, err := os.Open(path)
	if err != nil {
		panic("Invalid path -> " + path)
	}
	defer file.Close()

	var words []string

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		words = append(words, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
	fmt.Println(words)
	return words
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

				newFields := original.copy()

				new := Config{
					Url:     original.Url,
					Method:  original.Method,
					Headers: append(newFields.Headers, h),
					Cookies: append(newFields.Cookies, c),
					Data:    append(newFields.Data, d),
					Threads: original.Threads,
					Delay:   original.Delay,
				}
				newConfigs = append(newConfigs, new)

			}
		}
	}

	return newConfigs
}

func parseWordlist(words []string, pair Pair, isKey bool) []Pair {

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

type Expansion struct {
	Field       string
	Placeholder string
	Pairs       []Pair
}

func handleWordlist(config Config) []Config {

	fields := map[string]*[]Pair{
		"Headers": &config.Headers,
		"Cookies": &config.Cookies,
		"Data":    &config.Data,
	}
	expanded := make(map[string]map[string][]Pair)

	for _, wordlist := range config.Wordlists {

		words := readWordlists(wordlist.Value) // Path

		// Headers, Cookies and Data
		for name, field := range fields {

			// Field pairs, headers.key headers.value
			for _, pair := range *field {

				if pair.Key == wordlist.Key {
					removeWordlistPlaceholder(pair.Key, field)
					e := parseWordlist(words, pair, true)
					expanded[name][pair.Key] = append(expanded[name][pair.Key], e...)
					continue
				}
				if pair.Value == wordlist.Key {
					removeWordlistPlaceholder(pair.Value, field)
					e := parseWordlist(words, pair, false)
					expanded[name][pair.Value] = append(expanded[name][pair.Value], e...)
					continue
				}

			}

		}

	}
	fmt.Println(expanded)
	return insertsWordlist(config, expanded)

}
