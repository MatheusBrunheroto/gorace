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

func listener(
	barSize int,
	progress log.ProgressReader, finish chan int,
) {
	var sent, completed int

	total := <-progress.Total // Forces listener to stay off until end of wordlist reading
	s := progress.Sent
	c := progress.Completed

	for {
		select {
		case _, ok := <-s:
			if !ok {
				break
			}
			fmt.Println("SENT")
			sent++

		case _, ok := <-c:
			if !ok {
				break
			}
			fmt.Println("COMPLETE")
			completed++
		}

		bar := progressBar(completed, total, barSize)

		remaining := total - completed
		// sobe pra linha da barra

		fmt.Printf("\r\033[K%s -> Sent: [%d] Complete: [%d] Remaining: [%d]", bar, sent, completed, remaining)

		if completed == total {
			finish <- 1
			return
		}
	}
}

func Display(progress log.ProgressReader, finish chan int) error {

	barSize, err := handleAsciiArt()
	if err != nil {
		fmt.Println(err)
	}

	go listener(barSize+2, progress, finish) // +2 for some reason fixes a lot of imprecisions

	return nil
}
