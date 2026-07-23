package verbose

import (
	"fmt"
	"gorace/input"
	"gorace/log"
	"io"
	"net/http"
	"strings"
)

//

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

const contextRadius = 1

func WorkerResponse(hash uint64, resp *http.Response, match string, logChan chan<- log.Entry) {

	prefix := hashText(hash)

	respbody, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		WorkerError(hash, err.Error(), logChan)
		return
	}
	body := string(respbody)
	lines := strings.Split(body, "\n")

	idx, found := findMatchLineIndex(lines, match)

	if found {
		// --- Verbosity 2: janela de 50 linhas ao redor do match ---
		start := idx - contextRadius
		if start < 0 {
			start = 0
		}
		end := idx + contextRadius + 1
		if end > len(lines) {
			end = len(lines)
		}
		context := highlightMatch(strings.Join(lines[start:end], "\n"), match)
		contextText := fmt.Sprintf("%s- %d -\n%s", prefix, resp.StatusCode, context)
		logChan <- log.Entry{Text: contextText, Verbosity: 2}

		// --- Verbosity 3: só a linha do match, só o match grifado ---
		matchText := fmt.Sprintf("%s- %d - %s", prefix, resp.StatusCode, highlightMatch(lines[idx], match))
		logChan <- log.Entry{Text: matchText, Verbosity: 3}
	}

	// --- Verbosity 4: body inteiro, SEMPRE, grifado só se houver match ---
	fullBody := body
	if found {
		fullBody = highlightMatch(body, match)
	}
	fullText := fmt.Sprintf("%s- %d -\n%s", prefix, resp.StatusCode, fullBody)
	logChan <- log.Entry{Text: fullText, Verbosity: 4}
}

func findMatchLineIndex(lines []string, match string) (int, bool) {
	if match == "" {
		return -1, false
	}
	for i, line := range lines {
		if strings.Contains(line, match) {
			return i, true
		}
	}
	return -1, false
}

func highlightMatch(text string, match string) string {
	if match == "" || !strings.Contains(text, match) {
		return text
	}
	color := rgb{red: 255, green: 255, blue: 0}
	colored := fmt.Sprintf("\x1b[38;2;%d;%d;%dm%s\x1b[0m", color.red, color.green, color.blue, match)
	return strings.ReplaceAll(text, match, colored)
}
