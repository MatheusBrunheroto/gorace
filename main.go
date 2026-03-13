package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

var scanner = bufio.NewScanner(os.Stdin)

type Website struct {
	link    string
	method  int8 // Converted to [GET, POST, PUT, PATCH, DELETE, HEAD, OPTIONS] later
	headers map[string]string
	cookies map[string]string
	data    map[string]string
}

// Simplify the header, cookies and data collection
func readKeyValuePairs(text string) map[string]string {

	key_map := make(map[string]string)
	fmt.Printf("\tInsert \"%s_name : %s_value\"\n", text, text)

	for i := 1; ; i++ {

		fmt.Printf("\t  [%d] -> ", i)
		scanner.Scan()
		raw := scanner.Text()

		// This avoids string being empty, because it is mandatory to have ":"
		if !strings.Contains(raw, ":") {
			break
		}

		key := strings.SplitN(raw, ":", 2)
		if key[0] == "" || key[1] == "" {
			break
		}
		key[0] = strings.TrimSpace(key[0])
		key[1] = strings.TrimSpace(key[1])
		key_map[key[0]] = key[1] // key_map[key_name] = key_value

	}

	n := len(key_map)
	switch n {
	case 0:
		fmt.Printf("\t  Empty or invalid! Skipping... \n\n")
	case 1:
		fmt.Printf("\t  Empty or invalid! Only one %s was saved, skipping... \n\n", text)
	default:
		fmt.Printf("\t  Empty or invalid! First %d were saved, skipping... \n\n", n)
	}

	return key_map
}

func GetTargetInfo(websites *[]Website) {

	var link string
	var method int8

	for i := 1; ; i++ {

		// Reset temporary website to suit the struct, and avoid wrong values in future empty entries
		website := Website{
			link:    link,
			method:  method,
			headers: make(map[string]string),
			cookies: make(map[string]string),
			data:    make(map[string]string),
		}

		fmt.Printf("[%d] Insert link -> ", i)
		scanner.Scan()
		link = scanner.Text()
		if link == "" {
			break
		}

		fmt.Printf("\tInsert method [0 ... 6] -> ")
		fmt.Scan(&method)

		website.link = link
		website.method = method

		website.headers = readKeyValuePairs("header")
		website.cookies = readKeyValuePairs("cookie")
		website.data = readKeyValuePairs("data")

		*websites = append(*websites, website)
	}

}

func main() {

	// If opçao de wordlist, nao passar pelo gettarget, sim pelo readwordlist

	var websites = []Website{}
	GetTargetInfo(&websites)
	//fmt.Println(websites)

}
