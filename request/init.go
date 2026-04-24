package request

import (
	"gorace/input"
	"sync"
)

type modular struct {
}

// Returns the largest thread amount inside "websites"
func maxThreads(websites []input.Website) int {
	n := 1
	for _, w := range websites {
		if w.Threads > n {
			n = w.Threads
		}
	}
	return n
}

/*
Create N channels:
  - N = largest_thread_amount for "round" modes, so it loops through all websites at each N iteration
  - N = website_amount for "normal" modes, so each website loop runs their given amount of threads
*/
func initChan(n int) []chan struct{} {
	c := make([]chan struct{}, n)
	for i := range c {
		c[i] = make(chan struct{})
	}
	return c
}

/**/
func normalWorker(progressChannel [2]chan int, start chan struct{}, wg *sync.WaitGroup, websites ...input.Website) {
	website := websites[0]
	for i := 0; i < website.Threads; i++ {
		wg.Go(func() { worker(progressChannel, start, website) })
	}
}
func roundWorker(progressChannel [2]chan int, start chan struct{}, wg *sync.WaitGroup, websites ...input.Website) {
	for _, w := range websites {
		wg.Go(func() { worker(progressChannel, start, w) })
	}
}

func aux(progressChannel [2]chan int, start chan struct{}, isSequential bool,
	init func([2]chan int, chan struct{}, *sync.WaitGroup, ...input.Website),
	outerWg *sync.WaitGroup, websites ...input.Website) {

	if isSequential {
		var innerWg sync.WaitGroup                          // Starts inside worker, waits after function
		init(progressChannel, start, &innerWg, websites...) // Sequential use inner
		close(start)
		innerWg.Wait()
	} else {
		init(progressChannel, start, outerWg, websites...) //
		close(start)
	}

}

func generalInit(progressChannel [2]chan int, websites []input.Website, modes [2]bool) {

	isRound := modes[0]
	isSequential := modes[1]

	var outerWg sync.WaitGroup // CASCADE

	if isRound {
		startChannels := initChan(maxThreads(websites))
		for i := range startChannels {
			aux(progressChannel, startChannels[i], isSequential, roundWorker, &outerWg, websites...)
		}

	} else {
		startChannels := initChan(len(websites))
		for i := range websites {
			aux(progressChannel, startChannels[i], isSequential, normalWorker, &outerWg, websites[i])
		}

	}

	if isSequential {
		return
	} else {
		outerWg.Wait()
	}

}

func InitWorker(progressChannel [2]chan int, websites []input.Website, mode string) {

	// This constants are intended to make it easier to see the init parameters below
	const ROUND, NORMAL bool = true, false
	const SEQUENTIAL, CASCADE bool = true, false

	switch mode {

	// After N threads of an URL requests were sent to worker, waits for them to finish before starting next URL requests
	case "sequential":
		generalInit(progressChannel, websites, [2]bool{NORMAL, SEQUENTIAL})
	// Same as sequential, but doesn't wait for its requests to finish before starting the next URL requests
	case "cascade":
		generalInit(progressChannel, websites, [2]bool{NORMAL, CASCADE})

	// Sequential's behaviour, but cycles through the URLs requests for N times, N = largest amount of threads informed
	case "round-sequential":
		generalInit(progressChannel, websites, [2]bool{ROUND, SEQUENTIAL})
	// Cascade's behaviour, but cycles through the URLs requests for N times, N = largest amount of threads informed
	case "round-cascade":
		generalInit(progressChannel, websites, [2]bool{ROUND, CASCADE})

	// This is the default mode "flood", group all the requests and fire them at the exact same moment
	default:
		start := make(chan struct{})
		var wg sync.WaitGroup
		for _, w := range websites {
			for range w.Threads {
				wg.Go(func() { worker(start, w) })
			}
		}
		close(start)
		wg.Wait()
	} // WG GO INTERNO PRA GERAR COISA, ESPERAR POR DENTRO, ALGO DO TIPO, E UM POR FORA SEI LA

}
