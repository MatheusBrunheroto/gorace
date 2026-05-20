package request

import (
	"gorace/input"
	"io"
	"net/http"
	"net/url"
	"strings"
)

func missingHeaders(request *http.Request) {
	if request.UserAgent() == "" {
		request.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	}
}

// Converts data to put bellow headers, the request body
func getBody(rawData []input.Pair) io.Reader {
	var hasData bool = false

	data := url.Values{}
	for _, d := range rawData {
		if d.Key == "" {
			continue
		}
		data.Set(d.Key, d.Value)
		hasData = true
	}

	var body *strings.Reader = nil
	if hasData {
		body = strings.NewReader(data.Encode()) // Turns k1:v1 and k2:v2 to k1=v1&k2=v2
	}

	return body
}

func buildRequest(w input.Website) (*http.Request, error) {

	// DATA - Not mandatory, but the only way to insert in the request is by creating a body
	request, err := http.NewRequest(w.Method, w.Url, getBody(w.Data)) // get body vem aqui
	if err != nil {
		return &http.Request{}, err
	}

	// HEADERS - Mandatory, if none are informed, common headers will be added
	for _, h := range w.Headers {
		if h.Key == "" {
			continue
		}
		request.Header.Set(h.Key, h.Value)
	}
	missingHeaders(request)
	//	fmt.Println(request)
	// melhorar o filtro desses valores

	// COOKIES - Not mandatory
	for _, c := range w.Cookies {
		if c.Key == "" {
			continue
		}
		request.AddCookie(&http.Cookie{Name: c.Key, Value: c.Value})
	}

	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return request, nil

}

// Preciso adicionar no minimo valores vaizos de request
