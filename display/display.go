package display

import (
	"fmt"
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

type Session struct {
	Draw     chan string
	Ready    chan struct{}
	Finished chan struct{}
	Progress log.ProgressReader
}

func NewSession(progress log.ProgressReader) Session {

	return Session{
		Draw:     make(chan string),
		Ready:    make(chan struct{}),
		Finished: make(chan struct{}),
		Progress: progress,
	}

}

func progressBar(sent int, total int, amount int) string {

	s := float32(sent)
	t := float32(total)
	a := float32(amount)

	percentage := (100 * s) / t
	var bar string = "["

	i := float32(1)
	for ; i <= a; i++ {

		if percentage >= i*(100/a) {
			bar = bar + "#"
		} else {
			bar = bar + "-"
		}

	}
	bar = bar + "]"
	return bar

}

func incrementIfOpen(counter *int, received bool) {
	if received {
		(*counter)++
	}
}
func monitorProgress(barSize int, progress log.ProgressReader, finished chan<- struct{}) {

	total := <-progress.Total // Forces listener to stay off until wordlist reading finishes

	var sent, succeeded, failed int
	var completed, remaining int

	for {

		select {
		case _, isOpen := <-progress.Sent:
			incrementIfOpen(&sent, isOpen)

		case _, isOpen := <-progress.Succeeded:
			incrementIfOpen(&succeeded, isOpen)

		case _, isOpen := <-progress.Failed:
			incrementIfOpen(&failed, isOpen)
		}

		completed = succeeded + failed
		remaining = total - completed

		bar := progressBar(completed, total, barSize)
		fmt.Printf("\r\033[K%s -> Sent: [%d] Complete: [%d] Remaining: [%d]", bar, sent, succeeded, remaining)

		if completed == total {
			finished <- struct{}{}
			return
		}

	}
}

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
