package verbose

import (
	"fmt"
	"gorace/input"
	"gorace/log"
	"strconv"
)

// arrumar
func Worker(w input.Config, hash uint64, hit bool, logChan chan<- log.Entry) {

	hashColor := hashToVividColor(hash)

	hashPreview := fmt.Sprintf("\x1b[38;2;%d;%d;%dm[%s...] \x1b[0m",
		hashColor.red, hashColor.green, hashColor.blue,
		strconv.FormatUint(hash, 6)[:6])

	text := hashPreview
	var info string

	if hit {
		info = "- (Cache Hit) - "
		text += info + fmt.Sprintf("\x1b[38;2;%d;%d;%dmStarting... \x1b[0m", 100+int(hash%60), 100+int(hash%60), 100+int(hash%60))
	} else {
		info = "- NEW - "
		text += info + fmt.Sprintf("\x1b[38;2;%d;%d;%dmStarting... -> %s\x1b[0m", 100+int(hash%60), 100+int(hash%60), 100+int(hash%60), fmt.Sprint(w))
	}

	logChan <- log.Entry{Text: text, Verbosity: 3}

}

func WorkerError(hash uint64, err string, logChan chan<- log.Entry) {

	hashColor := hashToVividColor(hash)
	hashPreview := fmt.Sprintf("\x1b[38;2;%d;%d;%dm[%s...]\x1b[0m",
		hashColor.red, hashColor.green, hashColor.blue,
		strconv.FormatUint(hash, 6)[:6])

	errorColor := rgb{
		red:   100 + int(hash%60), // 60-199
		green: 100 + int(hash%60),
		blue:  100 + int(hash%60),
	}
	errorText := fmt.Sprintf("\x1b[38;2;%d;%d;%dm%s\x1b[0m",
		errorColor.red, errorColor.green, errorColor.blue,
		err)

	text := fmt.Sprintf("✘ - %s %s", hashPreview, errorText)
	logChan <- log.Entry{Text: fmt.Sprintf("%s", text), Verbosity: 1}

}
