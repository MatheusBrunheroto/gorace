package input

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

var scanner = bufio.NewScanner(os.Stdin)

type Website struct {
	Link    string
	Method  int8 // Converted to [GET, POST, PUT, PATCH, DELETE, HEAD, OPTIONS] later
	Headers map[string]string
	Cookies map[string]string
	Data    map[string]string
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
			Link:    link,
			Method:  method,
			Headers: make(map[string]string),
			Cookies: make(map[string]string),
			Data:    make(map[string]string),
		}

		fmt.Printf("[%d] Insert link -> ", i)
		scanner.Scan()
		link = scanner.Text()
		fmt.Printf("\n")
		if link == "" {
			break
		}

		fmt.Println("\t[GET(0), POST(1), PUT(2), PATCH(3), DELETE(4), HEAD(5), OPTIONS(6)]")
		for {
			fmt.Printf("\tInsert method [0 ... 6] -> ")
			fmt.Scan(&method)
			if method >= 0 && method <= 6 {
				fmt.Printf("\n")
				break
			}
		}

		website.Link = link
		website.Method = method

		website.Headers = readKeyValuePairs("header")
		website.Cookies = readKeyValuePairs("cookie")
		website.Data = readKeyValuePairs("data")

		*websites = append(*websites, website)
	}

}
