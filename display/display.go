package display

import (
	"fmt"
	"gorace/log"
)

// PEGA LINHA DO MAIS DE BAIXO
func progressBar(sent int, total int, amount int) string {

	// total * sent_percentage = 100 * sent
	s := float32(sent)
	t := float32(total)
	a := float32(amount)

	percentage := (100 * s) / t
	//fmt.Println(percentage)
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

func listener(progressChannel log.Progress, barSize int) {

	total := <-progressChannel.Total
	sentChannel := progressChannel.Sent
	completedChannel := progressChannel.Completed

	var sent, completed int

	for {

		select {
		case _, ok := <-sentChannel:
			if !ok {
				return
			}
			sent++

		case _, ok := <-completedChannel:
			if !ok {
				return
			}
			completed++
		}

		bar := progressBar(completed, total, barSize)

		remaining := total - completed
		// sobe pra linha da barra
		fmt.Print("\033[A")

		fmt.Printf("\r\033[K%s -> Sent: [%d] Complete: [%d] Remaining: %d", bar, sent, completed, remaining)
	}

}

func Display(progressChannel log.Progress) error {

	barSize, err := handleAsciiArt()
	if err != nil {
		fmt.Println(err)
	}

	go listener(progressChannel, barSize+2)

	return nil
}
