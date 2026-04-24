package display

import (
	"fmt"
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

func listener(progressChannel [3]chan int, barSize int) {

	total := <-progressChannel[0]
	sentChannel := progressChannel[1]
	completedChannel := progressChannel[2]

	for {
		sent, s := <-sentChannel
		completed, c := <-completedChannel
		if !s || !c {
			break
		}
		bar := progressBar(completed, total, barSize)
		remaining := total - completed
		fmt.Printf("\r%s -> Sent: [%d] Complete: [%d] Remaining:\\nn", bar, sent, completed, remaining)
	}

}

func Display(progressChannel [3]chan int) error {

	barSize, err := handleAsciiArt()
	if err != nil {
		fmt.Println(err)
	}

	go listener(progressChannel, barSize)
	//fmt.Println("\n")
	// sent total size
	//fmt.Println(size)

	return nil
}
