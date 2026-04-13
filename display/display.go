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

func Progress(wordlist <-chan int) error {

	size, err := handleAsciiArt()
	if err != nil {
		fmt.Println(err)
	}
	//fmt.Println("\n")
	// sent total size
	//fmt.Println(size)

	for {

		v, ok := <-wordlist
		if !ok {
			break
		}
		bar := progressBar(v, 99, size)
		fmt.Printf("\r->%s [%d/100]", bar, v)
		//fmt.Println(v)

	}
	fmt.Println("\n")
	fmt.Println("\n")
	fmt.Println("\n")
	// var total_workers pra usar ali em baixo

	for {

		v, ok := <-wordlist
		if !ok {
			break
		}

		fmt.Printf("\rProgress -> Sent: [%d] Complete: [%d]\\nn", v, v)
	}

	return nil
}
