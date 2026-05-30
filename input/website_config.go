package input

import "strconv"

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
	Url     string
	Method  string
	Headers []Pair
	Cookies []Pair
	Data    []Pair
	Threads int
	Delay   int
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
func writeCurrentConfig(current Config, flag string, value string) error {

	switch flag {

	// Only one flag can be called
	case "--url":
		if err := filterUrl(&value); err != nil {
			return err
		}
		current.Url = value

	case "--method":
		current.Method = filterMethod(value)
	case "--threads":
		current.Threads, _ = strconv.Atoi(value)
	case "--delay":
		current.Delay, _ = strconv.Atoi(value)

	// One or more flags can be called
	case "--headers", "--cookies", "--data", "--wordlists":
		handleField(current, flag, value)

	}

	return nil

}
func getConfigs(flags map[string]string, args []string) ([]Config, error) {

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

		// First URL CRIAR FUNCAO
		if alias == "--url" && i > 1 {
			configs = append(configs, current.copy()) // Save Current Config
			current = defaultConfig()                 // Set Current Config to Default
		}

		writeCurrentConfig(current, alias, value)

	}
	configs = append(configs, current.copy()) // Save Current Config
	current = defaultConfig()

}
