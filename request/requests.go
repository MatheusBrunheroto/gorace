package request

import (
	"fmt"
	"gorace/input"
	"gorace/log"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/cespare/xxhash/v2"
)

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

func missingHeaders(request *http.Request) {
	if request.UserAgent() == "" {
		request.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	}
}

func buildRequest(w input.Website) (*http.Request, error) {

	// DATA - Not mandatory, but the only way to insert in the request is by creating a body
	request, err := http.NewRequest(w.Method, w.Url, nil)
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

func verifyExistence(w input.Website) bool {
	

	return request
} // É UMA GO FUNCTION, 
func getBuiltRequests(){
	
	hashes := make(map[string]*http.Request) // NA VERDADE AQUI SERIA UM UM HASH PRA UM WEBSITE

	for {

		requestedHash := <-hashChan		
		for h, r := range hashes{
			if h == requestedHash {
				chan <- r
			}

		}
		
	}


}

// Always ends up doing N threads to the first website, and N for the other
// Receives a copy, so there is no need to thread lock
func worker(progressChannel log.Progress, start <-chan struct{}, w input.Website) {

	//		largest :=
	// xSlice(w.Headers, w.Cookies, w.Data)
	// O REQUEST É FEITO MULTIPLAS VEZES, TALVEZ ARRUMAR ISSO COM O request.clone
	if exists := verifyExistence() {

	}
	else{
		request, err := buildRequest(w)
	}
	websiteCode := fmt.Sprintf("%s%s%s%s%s%d", w.Url, w.Method, w.Headers, w.Cookies, w.Data, w.Threads)
	hash := xxhash.Sum64String(websiteCode)

	fmt.Printf("%x\n", hash)

	
	<-start
	progressChannel.Sent <- 1

	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		fmt.Println("\033[K ASS")
		progressChannel.Completed <- 1
		return
	}
	respbody, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		progressChannel.Completed <- 1
		return
	}
	resp.Body.Close()

	//fmt.Println(resp.Status, resp.ContentLength)
	if !strings.Contains(string(respbody), "Invalid username or password.") {
		//	fmt.Println(w.Data, resp.Header)
	}
	//We Read the response body on the line below.
	progressChannel.Completed <- 1

}

// Preciso adicionar no minimo valores vaizos de request
