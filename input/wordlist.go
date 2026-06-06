package input

import (
	"bufio"
	"fmt"
	"os"
)

type expansion struct {
	Field       string
	Placeholder string
	Pairs       []Pair
}

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

func removeWordlistPlaceholder(key string, field *[]Pair, isKey bool) {
	var filtered []Pair
	for _, p := range *field {
		if p.Key != key && isKey {
			filtered = append(filtered, p)
		} else if p.Value != key && !isKey {
			filtered = append(filtered, p)

		}
	}
	*field = filtered
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
func getCombinations(expanded map[string][]Pair, placeholders []string, index int) [][]Pair {

	if index == len(placeholders) {
		return [][]Pair{{}} // base case: uma combinação vazia
	}
	var result [][]Pair
	key := placeholders[index]
	pairs := expanded[key]

	rest := getCombinations(expanded, placeholders, index+1)

	for _, pair := range pairs {
		for _, combo := range rest {
			newCombo := append([]Pair{pair}, combo...)
			result = append(result, newCombo)
		}
	}

	return result

}

func insertsWordlist(original Config, expanded []expansion) []Config {

	var newConfigs []Config

	expandedHeaders := make(map[string][]Pair)
	expandedCookies := make(map[string][]Pair)
	expandedData := make(map[string][]Pair)

	for _, e := range expanded {

		switch e.Field {
		case "Headers":
			expandedHeaders[e.Placeholder] = append(expandedHeaders[e.Placeholder], e.Pairs...)
		case "Cookies":
			expandedCookies[e.Placeholder] = append(expandedCookies[e.Placeholder], e.Pairs...)
		case "Data":
			expandedData[e.Placeholder] = append(expandedData[e.Placeholder], e.Pairs...)
		}

	}

	var headerPlaceholders, cookiePlaceholders, dataPlaceholders []string
	for p := range expandedHeaders {
		headerPlaceholders = append(headerPlaceholders, p)
	}
	for p := range expandedCookies {
		cookiePlaceholders = append(cookiePlaceholders, p)
	}
	for p := range expandedData {
		dataPlaceholders = append(dataPlaceholders, p)
	}

	newHeaders := getCombinations(expandedHeaders, headerPlaceholders, 0)
	newCookies := getCombinations(expandedCookies, cookiePlaceholders, 0)
	newData := getCombinations(expandedData, dataPlaceholders, 0)

	// se vazio, coloca uma combinação vazia pra o loop funcionar
	if len(newHeaders) == 0 {
		newHeaders = [][]Pair{{}}
	}
	if len(newCookies) == 0 {
		newCookies = [][]Pair{{}}
	}
	if len(newData) == 0 {
		newData = [][]Pair{{}}
	}

	for _, h := range newHeaders {
		for _, c := range newCookies {
			for _, d := range newData {
				newFields := original.copy()
				new := Config{
					Url:     original.Url,
					Method:  original.Method,
					Headers: append(newFields.Headers, h...),
					Cookies: append(newFields.Cookies, c...),
					Data:    append(newFields.Data, d...),
					Threads: original.Threads,
					Delay:   original.Delay,
				}
				newConfigs = append(newConfigs, new)
			}
		}
	}
	return newConfigs
}

func handleWordlist(config Config) []Config {

	fields := map[string]*[]Pair{
		"Headers": &config.Headers,
		"Cookies": &config.Cookies,
		"Data":    &config.Data,
	}

	var expansions []expansion
	// var placeholder string
	for _, wordlist := range config.Wordlists {

		words := readWordlists(wordlist.Value) // Path

		// Headers, Cookies and Data
		for name, field := range fields {
			// Field pairs, headers.key headers.value
			for _, pair := range *field {

				if pair.Key == wordlist.Key {

					removeWordlistPlaceholder(pair.Key, field, true)
					e := expansion{
						Field:       name,
						Placeholder: pair.Key,
						Pairs:       parseWordlist(words, pair, true),
					}
					expansions = append(expansions, e)
					continue

				}
				if pair.Value == wordlist.Key {
					removeWordlistPlaceholder(pair.Value, field, false)
					e := expansion{
						Field:       name,
						Placeholder: pair.Value,
						Pairs:       parseWordlist(words, pair, false),
					}
					expansions = append(expansions, e)
					continue
				}

			}

		}

	}

	return insertsWordlist(config, expansions)
}
