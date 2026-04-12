package input

/* gorace -u 'url' -H 'Content-Type: application/json'

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
}
