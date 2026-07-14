package request

import (
	"gorace/input"
	"gorace/log"
	"sync"
)

// Returns the largest thread amount inside "configs"
func maxThreads(configs []input.Config) int {
	n := 1
	for _, w := range configs {
		if w.Threads > n {
			n = w.Threads
		}
	}
	return n
}

// Returns the sum of all thread amounts inside "configs"
func totalThreads(configs []input.Config) int {
	n := 0
	for _, w := range configs {
		n += w.Threads
	}
	return n
}

/*
Create N channels:
  - N = largest_thread_amount for "round" modes, so it loops through all configs at each N iteration
  - N = config_amount for "normal" modes, so each config loop runs their given amount of threads
*/

/**/
func runWorkers(configs []input.Config, round bool, sequential bool, iterations int, ch WorkerChans, logChan chan<- log.Entry) {

	var outerWg sync.WaitGroup

	for i := 0; i < iterations; i++ {
		start := make(chan struct{})

		var innerWg sync.WaitGroup
		currentWg := &outerWg
		if sequential {
			currentWg = &innerWg
		}

		// If round, will send len(configs)
		if round {
			for _, w := range configs {
				currentWg.Go(func() { worker(start, w, ch, logChan) })
			}

		} else {
			w := configs[i]
			for t := 0; t < w.Threads; t++ {
				currentWg.Go(func() { worker(start, w, ch, logChan) })
			}
		}
		close(start)

		// If sequential every iteration will wait for response
		if sequential {
			innerWg.Wait()
		}
	}

	// If cascade, wait for response after sending every single possibility inside the for loop
	if !sequential {
		outerWg.Wait()
	}

}

func InitWorkers(configs []input.Config, mode string, ch WorkerChans, logChan chan<- log.Entry) { // Intended behavior is below the function

	// Easier to see the init parameters
	const ROUND, NORMAL bool = true, false
	const SEQUENTIAL, CASCADE bool = true, false

	switch mode {

	// After N threads of an URL requests were sent to worker, waits for them to finish before starting next URL requests
	case "sequential":
		ch.Progress.Total <- totalThreads(configs)
		runWorkers(configs, NORMAL, SEQUENTIAL, len(configs), ch, logChan)
	// Same as sequential, but doesn't wait for its requests to finish before starting the next URL requests
	case "cascade":
		ch.Progress.Total <- totalThreads(configs)
		runWorkers(configs, NORMAL, CASCADE, len(configs), ch, logChan)

	// Sequential's behaviour, but cycles through the URLs requests for N times, N = largest amount of threads informed
	case "round-sequential":
		ch.Progress.Total <- maxThreads(configs) * len(configs)
		runWorkers(configs, ROUND, SEQUENTIAL, maxThreads(configs), ch, logChan)
	// Cascade's behaviour, but cycles through the URLs requests for N times, N = largest amount of threads informed
	case "round-cascade":
		ch.Progress.Total <- maxThreads(configs) * len(configs)
		runWorkers(configs, ROUND, CASCADE, maxThreads(configs), ch, logChan)

	// This is the default mode "flood", group all the requests and fire them at the exact same moment
	default:
		start := make(chan struct{})
		ch.Progress.Total <- totalThreads(configs)
		var wg sync.WaitGroup
		for _, w := range configs {
			for range w.Threads {
				wg.Go(func() { worker(start, w, ch, logChan) })
			}
		}
		close(start)
		wg.Wait()
	}

}

/* INTENDED BEHAVIOR WHEN TESTING/DEBUGGING:

Assuming 2 URLs, URL_1 and URL_2, both containing a bool variable .sent and .completed,
when simulating an attack, with the command "gorace -u 'URL_1' --threads 2 -u 'URL_2' --threads 3 --mode MODE",
the expected output for each mode should be:

--mode sequential
  - URL_1.sent (2x)
  - URL_1.completed (2x)
  - URL_2.sent (3x)
  - URL_2.completed (3x)

--mode cascade
  - URL_1.sent (2x)
  - URL_2.sent (3x)			-- The .completed can arrive at any moment, it doesn't influenciate in the order

--mode round-sequential (will repeat 3x since maxThreads informed = 3)
  - URL_1.sent
  - URL_1.completed
  - URL_2.sent
  - URL_2.completed
  - (+2x)

--mode round-cascade (will repeat 3x since maxThreads informed = 3)
  - URL_1.sent
  - URL_2.sent
  - (+2x)					-- The .completed can arrive at any moment, it doesn't influenciate in the order

--mode flood (it's the default): Everything is sent at the same time, good for single endpoint
*/
