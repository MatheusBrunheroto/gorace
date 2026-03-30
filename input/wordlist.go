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
	-> Add the words from the wordlist to new_key_map[index][words] = key_value

	Case key_value == placeholder:
		- key_map[key_name] = placeholder
	-> Add the words from the wordlist to new_key_map[index][key_name] = words

	Case key_name == placeholder && key_value == placeholder:
		- key_map[placeholder1] = placeholder2
	-> Requires different techniques


	wordlist = {"placeholder":"path"}
	key = {"name":"value"}
	key = {"placeholder":"value"} -> key = {"open.path()":"value"}
*/

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

func insertWordlist(new_key_map map[int]map[string]string, index *int, placeholder string, path string, method string) error {

	words, err := readWordlists(path)
	if err != nil {
		return err
	}

	if method == "name" {
		for _, w := range words {
			if new_key_map[*index] == nil {
				new_key_map[*index] = make(map[string]string)
			}
			new_key_map[*index][w] = placeholder
			*index++
		}

	} else if method == "value" {
		for _, w := range words {
			if new_key_map[*index] == nil {
				new_key_map[*index] = make(map[string]string)
			}
			new_key_map[*index][placeholder] = w
			*index++
		}
	}

	return nil

}

func handleWordlist(key_map map[string]string, new_key_map map[int]map[string]string, wordlists map[string]string) error {

	var index int = 0

	// First, copy the values on the key_map that aren't placeholders, the placeholders will go to insertWordlist
	for k, v := range key_map {

		// If key_name or key_value are placeholders, they can be called inside "wordlists"
		_, key_is_placeholder := wordlists[k]
		_, value_is_placeholder := wordlists[v]

		if key_is_placeholder && value_is_placeholder {
			fmt.Println("gozada insana")
			//insertWordlist(new_key_map, &index, k, wordlists[k])

		} else if key_is_placeholder {
			insertWordlist(new_key_map, &index, k, wordlists[k], "name") // BOTAR FVERIFICACAO DE ERRO AQUI
			delete(wordlists, k)

		} else if value_is_placeholder {
			insertWordlist(new_key_map, &index, v, wordlists[v], "value")
			delete(wordlists, v)

		} else {
			if new_key_map[index] == nil {
				new_key_map[index] = make(map[string]string)
			}
			new_key_map[index][k] = v // Adds non wordlist keys to the new_key_map
			index++
		}
	}

	return nil
}
