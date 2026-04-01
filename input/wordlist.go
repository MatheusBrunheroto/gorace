package input

import (
	"bufio"
	"errors"
	"fmt"
	"os"
)

/*
Map structures:

	- wordlists_map[placeholder] = path

	Case key_name == placeholder:
		- key_map[placeholder] = key_value
	-> Add the words from the wordlist to new_entry[index][words] = key_value

	Case key_value == placeholder:
		- key_map[key_name] = placeholder
	-> Add the words from the wordlist to new_entry[index][key_name] = words

	Case key_name == placeholder && key_value == placeholder:
		- key_map[placeholder1] = placeholder2
	-> Requires different techniques


	wordlist = {"placeholder":"path"}
	key = {"name":"value"}
	key = {"placeholder":"value"} -> key = {"open.path()":"value"}
*/

func handleWordlist(entry []KeyValue, wordlists []KeyValue) ([]KeyValue, error) {

	var new_entry []KeyValue
	wordlists_map := sliceToMap(wordlists) // If someone uses -w WORDLIST1=path1 and -w WORDLIST1=path2, the path1 is ignored

	// First, copy the values on the key_map that aren't placeholders, the placeholders will go to insertWordlist
	for _, kv := range entry {

		// If key_name or key_value are placeholders, they can be called inside "wordlists"
		_, keyIsPlaceholder := wordlists_map[kv.Key]
		_, valueIsPlaceholder := wordlists_map[kv.Value]

		if keyIsPlaceholder && valueIsPlaceholder {
			fmt.Println("special case")

		} else if keyIsPlaceholder {
			if err := insertWordlist(&new_entry, kv.Value, wordlists_map[kv.Key], "name"); err != nil {
				return []KeyValue{}, err
			}
		} else if valueIsPlaceholder {
			if err := insertWordlist(&new_entry, kv.Key, wordlists_map[kv.Value], "value"); err != nil {
				return []KeyValue{}, err
			}

		} else {

			new_entry = append(new_entry, KeyValue{Key: kv.Key, Value: kv.Value}) // Adds non wordlist keys to the new_entry

		}
	}

	return new_entry, nil
}
func sliceToMap(wordlists []KeyValue) map[string]string {
	wordlists_map := make(map[string]string, len(wordlists))
	for _, kv := range wordlists {
		wordlists_map[kv.Key] = kv.Value
	}
	return wordlists_map
}

func insertWordlist(new_entry *[]KeyValue, placeholder string, path string, method string) error {

	words, err := readWordlists(path)
	if err != nil {
		return err
	}

	for _, w := range words {

		if method == "name" {
			*new_entry = append(*new_entry, KeyValue{Key: w, Value: placeholder})
		} else {
			*new_entry = append(*new_entry, KeyValue{Key: placeholder, Value: w})
		}

	}

	return nil

}

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

	return words, nil
}
