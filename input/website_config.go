package input

import (
	"gorace/log"
	"math"
	"strconv"
	"time"
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
	Timeout   time.Duration
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
		Timeout:   c.Timeout,
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
		Timeout: 30 * time.Second,
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
		threads, err := strconv.Atoi(raw)
		if err != nil || threads <= 0 {
			panic("[x] Invalid threads (use a positive integer) -> " + raw)
		}
		current.Threads = threads
	case "--delay":
		delay, err := strconv.Atoi(raw)
		if err != nil || delay < 0 || int64(delay) > math.MaxInt64/int64(time.Millisecond) {
			panic("[x] Invalid delay (use nonnegative milliseconds within the supported range) -> " + raw)
		}
		current.Delay = delay
	case "--timeout":
		timeout, err := time.ParseDuration(raw)
		if err != nil || timeout < 0 {
			panic("[x] Invalid timeout (use 30s, 2m, or 0) -> " + raw)
		}
		current.Timeout = timeout

	// One or more flags
	case "--headers", "--cookies", "--data", "--wordlist":

		pairs := parsePairs(raw, flag)

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
		if alias == "--no-color" || alias == "--help" {
			continue
		}
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
		if alias == "--url" && current.Url != "" {
			configs = append(configs, current.copy()) // Save Current Config
			current = defaultConfig()                 // Set Current Config to Default
		}

		writeConfig(&current, alias, value, logChan)

	}
	if current.Url == "" {
		panic("[x] No Website URL was informed (-u or --url)")
	}
	configs = append(configs, current) // To save the Last URL

	return configs

}
