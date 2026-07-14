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

func Run(progress log.ProgressReader, finished chan struct{}, logChan chan log.Entry) {

	// 1. Logo
	handleAsciiArt()
	go log.Run(logChan, 1) // ler a bosta
	go progressMonitor(30, progress, finished, logChan)

}
