package display

import (
	"fmt"
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

func Run(session Session) error {

	barSize, err := handleAsciiArt()
	if err != nil {
		fmt.Println(err)
	}

	session.Ready <- struct{}{}

	fmt.Printf("%s\n\n", strings.Repeat("⸺", barSize/2))

	monitorFinished := make(chan struct{})
	go monitorProgress(barSize, session.Progress, monitorFinished) // +2 for some reason fixes a lot of imprecisions

	for {
		select {

		case <-monitorFinished:
			session.Finished <- struct{}{}
			return nil
		}

	}

	// return nil
}
