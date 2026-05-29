package input

import (
	"errors"
	"strings"
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

type Website struct {
	Url     string
	Method  string
	Headers []Pair
	Cookies []Pair
	Data    []Pair
	Threads int
	Delay   int
}

func (w Website) New() Website {
	return Website{
		Url:     "",
		Method:  "GET",
		Headers: []Pair{},
		Cookies: []Pair{},
		Data:    []Pair{},
		Threads: 1,
		Delay:   0,
	}

}
func fillWebsiteDefaults(flag *Flag, urlAmount int, name string) error {

	if flag.exists == false {

		parameter := ""

		if strings.Contains(name, "--method") {
			parameter = "GET"
		}
		if strings.Contains(name, "--threads") {
			parameter = "1"
		}
		if strings.Contains(name, "--delay") {
			parameter = "0"
		}

		flag.raw = append(flag.raw, parameter)
		return nil
	}

	// Flag exists, may have more than one
	flagAmount := len(flag.raw)
	if flagAmount > urlAmount {
		return errors.New("Two or more equal flags detected! -> " + flag.name)
	}

	return nil
}
