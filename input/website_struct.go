package input

type KeyValue struct {
	Key   string
	Value string
}

type Website struct {
	Url     string
	Method  string
	Headers []KeyValue
	Cookies []KeyValue
	Data    []KeyValue
}
