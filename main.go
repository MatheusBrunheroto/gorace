package main

import (
	"fmt"
	"gorace/input"
	"gorace/request"
)

func main() {

	// If opçao de wordlist, nao passar pelo gettarget, sim pelo readwordlist

	var websites = []input.Website{}
	input.GetTargetInfo(&websites)
	fmt.Println(websites)
	request.RequestThread(websites)

}
