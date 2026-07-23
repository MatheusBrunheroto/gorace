package verbose

import (
	"fmt"
	"gorace/input"
	"gorace/log"
	"io"
	"net/http"
	"regexp"
	"strings"
)

var whitespaceRegex = regexp.MustCompile(`\s+`)

func highlightMatch(text string, match string) string {
	if match == "" || !strings.Contains(text, match) {
		return text
	}
	color := rgb{red: 255, green: 255, blue: 0}
	colored := fmt.Sprintf("\x1b[38;2;%d;%d;%dm%s\x1b[0m", color.red, color.green, color.blue, match)
	return strings.ReplaceAll(text, match, colored)
}

func findMatchLineIndex(lines []string, match string) (string, bool) {
	if match == "" {
		return "", false
	}
	for _, line := range lines {
		if strings.Contains(line, match) {
			return line, true
		}
	}
	return "", false
}

func Worker(w input.Config, hash uint64, hit bool, logChan chan<- log.Entry) {

	preview := hashText(hash)
	var text string

	if hit {
		text = preview + "- (Cache Hit) - " + fmt.Sprintf("\x1b[38;2;%d;%d;%dmStarting... \x1b[0m", 100+int(hash%60), 100+int(hash%60), 100+int(hash%60))
	} else {
		text = preview + "- NEW - " + fmt.Sprintf("\x1b[38;2;%d;%d;%dmStarting... -> %s\x1b[0m", 100+int(hash%60), 100+int(hash%60), 100+int(hash%60), fmt.Sprint(w))
	}

	logChan <- log.Entry{Text: text, Verbosity: 2}

}

func WorkerError(hash uint64, err string, logChan chan<- log.Entry) {

	preview := hashText(hash)

	errorColor := rgb{
		red:   100 + int(hash%60),
		green: 100 + int(hash%60),
		blue:  100 + int(hash%60),
	}
	errorText := fmt.Sprintf("\x1b[38;2;%d;%d;%dm%s\x1b[0m",
		errorColor.red, errorColor.green, errorColor.blue,
		err)

	text := fmt.Sprintf("✘ - %s %s", preview, errorText)
	logChan <- log.Entry{Text: fmt.Sprintf("%s", text), Verbosity: 1}

}

func WorkerResponse(hash uint64, resp *http.Response, match string, logChan chan<- log.Entry) {

	prefix := hashText(hash)

	respbody, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		WorkerError(hash, err.Error(), logChan)
		return
	}
	body := string(respbody)
	body = strings.ReplaceAll(body, "\n", "")
	body = whitespaceRegex.ReplaceAllString(body, " ")

	lines := strings.Split(body, "><")

	line, found := findMatchLineIndex(lines, match)

	if found {

		// --verbosity 1 -> Match warning
		logChan <- log.Entry{Text: fmt.Sprintf("%s- %d | Match Found!", prefix, resp.StatusCode), Verbosity: 1}

		// --verbosity 2 -> Surroudings of highlighted match
		logChan <- log.Entry{Text: fmt.Sprintf("-> %s\n", highlightMatch(line, match)), Verbosity: 2}

		// --verbosity 3 -> Entire response with match highlighted
		logChan <- log.Entry{Text: highlightMatch(body, match), Verbosity: 3}

		return
	}

	// --- Verbosity 4: Entire bodies, if has match it will be highlighted on verbosity 3
	logChan <- log.Entry{Text: fmt.Sprintf("%s- %d -> %s", prefix, resp.StatusCode, body), Verbosity: 4}

}
