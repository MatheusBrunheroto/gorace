package input

import (
	"fmt"
	"strconv"
)

/*
	 gorace -u 'url' -H 'Content-Type: application/json'

		Arguments -> "-u", "url", "-H", "Content-Type: application/json"
		Flags -> "-u", "-H"
		Parameter ("-u") -> "url"

		Headers -> pair
		(Header) KeyName -> "Content-Type"
		(Header) Pair -> "application/json"

So, the struct PAIR will have:
  - The address of the PARAMETERS provided by the FLAGS
  - The key and key PAIR
*/

type Pair struct {
	Key   string
	Value string
}

type Config struct {
	Url       string
	Method    string
	Headers   []Pair
	Cookies   []Pair
	Data      []Pair
	Wordlists []Pair
	Threads   int
	Delay     int
}

func (c Config) copy() Config {
	return Config{
		Url:       c.Url,
		Method:    c.Method,
		Headers:   append([]Pair(nil), c.Headers...),
		Cookies:   append([]Pair(nil), c.Cookies...),
		Data:      append([]Pair(nil), c.Data...),
		Wordlists: append([]Pair(nil), c.Wordlists...),
		Threads:   c.Threads,
		Delay:     c.Delay,
	}
}

func defaultConfig() Config {
	return Config{
		Url:     "",
		Method:  "GET",
		Headers: []Pair{},
		Cookies: []Pair{},
		Data:    []Pair{},
		Threads: 1,
		Delay:   0,
	}
}

func writeConfig(current *Config, flag string, raw string) {

	switch flag {

	// Only one flag can be called
	case "--url":
		normalizeUrl(&raw)
		current.Url = raw

	case "--method":
		current.Method = normalizeMethod(raw)
	case "--threads":
		current.Threads, _ = strconv.Atoi(raw)
	case "--delay":
		current.Delay, _ = strconv.Atoi(raw)

	// One or more flags can be called
	case "--headers", "--cookies", "--data", "--wordlist":

		pairs := parsePairs(raw)

		switch flag {
		case "--headers":
			current.Headers = append(current.Headers, pairs...)
		case "--cookies":
			current.Cookies = append(current.Cookies, pairs...)
		case "--data":
			current.Data = append(current.Data, pairs...)
		case "--wordlist":
			current.Wordlists = append(current.Wordlists, pairs...)
		}

	}

}

func getConfigs(flags map[string]string, args []string) []Config {

	var configs []Config
	current := defaultConfig()

	// se flag url, modifica Config.URL
	var alias, value string
	for i := 0; i < len(args); i++ {

		// i % 2 is verified beforehand
		alias = args[i]
		i++
		value = args[i]

		_, exist := flags[alias]
		if !exist {
			continue
		}

		// First URL
		if alias == "--url" && i > 1 {
			configs = append(configs, current.copy()) // Save Current Config
			current = defaultConfig()                 // Set Current Config to Default
		}

		writeConfig(&current, alias, value)
		fmt.Println(current)

	}
	configs = append(configs, current) // To save the Last URL

	return configs

}
