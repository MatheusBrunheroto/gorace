package input

type Website struct {
	Url     string
	Method  string // [GET, POST, PUT, PATCH, DELETE, HEAD, OPTIONS]
	Headers map[string]string
	Cookies map[string]string
	Data    map[string]string
}
