package display

import (
	"gorace/log"
	"strings"
)

/*
━━━━━━━━━━━━━━━━━━━━━━━━━━━━
__     _____   ____   _____
||    ((   )) (( ___ ((   ))
||__|  \\_//   \\_||  \\_//

⸺⸺⸺⸺⸺⸺⸺⸺
[!] CLI
[+] Input
[x] Feedback
⸺⸺⸺⸺⸺⸺⸺⸺
Url: https://REQUEST_CONFIG
Method: POST
Headers: User-Agent: chrome
Cookies: session=vn2yu7908
━━━━━━━━━━━━━━━━━━━━━━━━━━━━

REQUESTS RESPONSES

[##################--------] -> Sent: [i] Complete: [j] Remaining: [k]
*/

func Separator(v int, logChan chan log.Entry) {
	separator := strings.Repeat("⸺", 30)
	logChan <- log.Entry{Text: separator + "\n", Verbosity: v}
}

func Run(progress log.ProgressReader, logChan chan log.Entry) {
	handleAsciiArt()
	go progressMonitor(60, progress, logChan)
}
