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

type wordlist struct {
	Placeholder string
	Words       []string
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
		panic("[x] " + err.Error())
	}

	return words
}

func parseWordlist(pair Pair, lists []wordlist) []Pair {

	var expanded []Pair

	var expand func(int, []string)
	expand = func(index int, replacements []string) {

		if index == len(lists) {
			// Replace in the original pair so inserted words are not expanded again.
			replacer := strings.NewReplacer(replacements...)
			expanded = append(expanded, Pair{
				Key:   replacer.Replace(pair.Key),
				Value: replacer.Replace(pair.Value),
			})
			return
		}

		for _, word := range lists[index].Words {
			expand(index+1, append(replacements, lists[index].Placeholder, word))
		}

	}
	expand(0, nil)

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
					Timeout: original.Timeout,
				}
				newConfigs = append(newConfigs, new)
			}
		}
	}
	return newConfigs
}

func handleExpansion(pair Pair, fieldName string, lists []wordlist) expansion {

	return expansion{
		Field:       fieldName,
		Placeholder: randomString(8),
		Pairs:       parseWordlist(pair, lists),
	}

}
func handleWordlist(config Config) []Config {

	fields := map[string]*[]Pair{
		"Headers": &config.Headers,
		"Cookies": &config.Cookies,
		"Data":    &config.Data,
	}

	var expansions []expansion
	var lists []wordlist

	for _, list := range config.Wordlists {
		lists = append(lists, wordlist{Placeholder: list.Key, Words: readWordlists(list.Value)})
	}

	// Headers, Cookies and Data
	for name, field := range fields {
		var filtered []Pair
		for _, pair := range *field {
			var matching []wordlist
			for _, list := range lists {
				if strings.Contains(pair.Key, list.Placeholder) || strings.Contains(pair.Value, list.Placeholder) {
					matching = append(matching, list)
				}
			}

			if len(matching) > 0 {
				expansions = append(expansions, handleExpansion(pair, name, matching))
				continue
			}
			filtered = append(filtered, pair)
		}
		*field = filtered
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
