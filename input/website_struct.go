package input

type Website struct {
	Link    string
	Method  int8 // Converted to [GET, POST, PUT, PATCH, DELETE, HEAD, OPTIONS] later
	Headers map[string]string
	Cookies map[string]string
	Data    map[string]string
}
