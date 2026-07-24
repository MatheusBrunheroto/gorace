package input

import (
	"gorace/log"
	"strconv"
)

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

func writeConfig(current *Config, flag string, raw string, logChan chan<- log.Entry) {

	switch flag {

	// Only one flag
	case "--url":
		normalizeUrl(&raw, logChan)
		current.Url = raw
	case "--method":
		current.Method = normalizeMethod(raw, logChan)
	case "--threads":
		current.Threads, _ = strconv.Atoi(raw)
	case "--delay":
		current.Delay, _ = strconv.Atoi(raw)

	// One or more flags
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

func getConfigs(flags map[string]string, args []string, logChan chan<- log.Entry) []Config {

	var configs []Config
	current := defaultConfig()

	var alias, value string
	for i := 0; i < len(args); i++ {

		alias = args[i]
		i++
		if i >= len(args) {
			continue
		}
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

		writeConfig(&current, alias, value, logChan)

	}
	configs = append(configs, current) // To save the Last URL

	return configs

}
