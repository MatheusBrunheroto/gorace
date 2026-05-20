package display

import (
	"fmt"
	"gorace/log"
)

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
