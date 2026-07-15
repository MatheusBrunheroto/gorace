package display

import (
	"gorace/log"
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

func Run(progress log.ProgressReader, logChan chan log.Entry) {

	// 1. Logo
	handleAsciiArt()
	go progressMonitor(60, progress, logChan)

}
