package request

import (
	"gorace/input"
	"gorace/log"
	"sync"
)

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

/**/
func runWorkers(websites []input.Website, round bool, sequential bool, progress log.ProgressWriter) {
	var outerWg sync.WaitGroup

	var loops int
	if round {
		loops = maxThreads(websites)
	} else {
		loops = len(websites)
	}

	for i := 0; i < loops; i++ {
		start := make(chan struct{})

		var innerWg sync.WaitGroup
		currentWg := &outerWg
		if sequential {
			currentWg = &innerWg
		}

		if round {
			for _, w := range websites {
				currentWg.Go(func() { worker(start, w, progress) })
			}
		} else {
			w := websites[i]
			for t := 0; t < w.Threads; t++ {
				currentWg.Go(func() { worker(start, w, progress) })
			}
		}

		close(start)

		if sequential {
			innerWg.Wait()
		}
	}

	if !sequential {
		outerWg.Wait()
	}
}

func InitWorkers(websites []input.Website, mode string, progress log.ProgressWriter) {

	// This constants are intended to make it easier to see the init parameters below
	const ROUND, NORMAL bool = true, false
	const SEQUENTIAL, CASCADE bool = true, false

	switch mode {

	// After N threads of an URL requests were sent to worker, waits for them to finish before starting next URL requests
	case "sequential":
		runWorkers(websites, NORMAL, SEQUENTIAL, progress)
	// Same as sequential, but doesn't wait for its requests to finish before starting the next URL requests
	case "cascade":
		runWorkers(websites, NORMAL, CASCADE, progress)

	// Sequential's behaviour, but cycles through the URLs requests for N times, N = largest amount of threads informed
	case "round-sequential":
		runWorkers(websites, ROUND, SEQUENTIAL, progress)
	// Cascade's behaviour, but cycles through the URLs requests for N times, N = largest amount of threads informed
	case "round-cascade":
		runWorkers(websites, ROUND, CASCADE, progress)

	// This is the default mode "flood", group all the requests and fire them at the exact same moment
	default:
		start := make(chan struct{})
		var wg sync.WaitGroup
		for _, w := range websites {
			for range w.Threads {
				wg.Go(func() { worker(start, w, progress) })
			}
		}
		close(start)
		wg.Wait()
	}

}
