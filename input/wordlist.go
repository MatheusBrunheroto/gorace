package input

import (
	"bufio"
	"math/rand"
	"os"
	"strings"
	"time"
)

type expansion struct {
	Field       string
	Placeholder string
	Pairs       []Pair
}

func readWordlists(path string) []string {

	file, err := os.Open(path)
	if err != nil {
		panic("[x] Invalid path -> " + path)
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

	return words
}

func removeWordlistPlaceholder(key string, field *[]Pair, isKey bool) {
	var filtered []Pair
	removed := false
	for _, p := range *field {
		match := (isKey && p.Key == key) || (!isKey && p.Value == key)
		if match && !removed {
			removed = true
			continue
		}
		filtered = append(filtered, p)
	}
	*field = filtered
}
func parseWordlist(words []string, pair Pair, wrap [2]string, isKey bool) []Pair {

	var expanded []Pair

	for _, word := range words {

		w := wrap[0] + word + wrap[1]

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
		return [][]Pair{{}}
	}
	key := placeholders[index]
	pairs := expanded[key]
	rest := getCombinations(expanded, placeholders, index+1)

	var result [][]Pair
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

func filterPlaceholder(placeholder string, wordlistPlaceholder string) [2]string {
	before, after, _ := strings.Cut(placeholder, wordlistPlaceholder)
	return [2]string{before, after}
}
func handleExpansion(pair Pair, fieldName string, words []string, wordlistKey string, isKey bool) expansion {

	target := pair.Key
	if !isKey {
		target = pair.Value
	}
	return expansion{
		Field:       fieldName,
		Placeholder: randomString(8),
		Pairs:       parseWordlist(words, pair, filterPlaceholder(target, wordlistKey), isKey),
	}

}
func handleWordlist(config Config) []Config {

	fields := map[string]*[]Pair{
		"Headers": &config.Headers,
		"Cookies": &config.Cookies,
		"Data":    &config.Data,
	}

	var expansions []expansion

	for _, wordlist := range config.Wordlists {

		words := readWordlists(wordlist.Value) // Path

		// Headers, Cookies and Data
		for name, field := range fields {
			// Field pairs -> headers.key headers.value...
			for _, pair := range *field {

				if strings.Contains(pair.Key, wordlist.Key) {

					removeWordlistPlaceholder(pair.Key, field, true)
					expansions = append(expansions, handleExpansion(pair, name, words, wordlist.Key, true))
					continue

				}
				if strings.Contains(pair.Value, wordlist.Key) {

					removeWordlistPlaceholder(pair.Value, field, false)
					expansions = append(expansions, handleExpansion(pair, name, words, wordlist.Key, false))
					continue

				}

			}

		}

	}

	return insertsWordlist(config, expansions)
}

func randomString(n int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, n)
	for i := range b {
		b[i] = charset[r.Intn(len(charset))]
	}
	return string(b)
}
