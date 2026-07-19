package request

import (
	"gorace/input"
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
func runWorkers(configs []input.Config, modes [2]bool, match string, ch WorkerChans) {

	isRound := modes[0]
	isSequential := modes[1]

	iterations := len(configs)
	if isRound {
		iterations = maxThreads(configs)
	}

	var outerWg sync.WaitGroup

	for i := 0; i < iterations; i++ {
		start := make(chan struct{})

		var innerWg sync.WaitGroup
		currentWg := &outerWg
		if isSequential {
			currentWg = &innerWg
		}

		// If isRound, will send len(configs)
		if isRound {
			for _, w := range configs {
				currentWg.Go(func() { worker(start, w, match, ch) })
			}

		} else {
			w := configs[i]
			for t := 0; t < w.Threads; t++ {
				currentWg.Go(func() { worker(start, w, match, ch) })
			}
		}
		close(start)

		// If isSequential every iteration will wait for response
		if isSequential {
			innerWg.Wait()
		}
	}

	// If cascade, wait for response after sending every single possibility inside the for loop
	if !isSequential {
		outerWg.Wait()
	}

}

func InitWorkers(configs []input.Config, global input.GlobalFlags, ch WorkerChans) { // Intended behavior is below the function

	// Easier to see the init parameters
	const ROUND, NORMAL bool = true, false
	const SEQUENTIAL, CASCADE bool = true, false

	switch global.Mode {

	// After N threads of an URL requests were sent to worker, waits for them to finish before starting next URL requests
	case "sequential":
		ch.Progress.Total <- totalThreads(configs)
		runWorkers(configs, [2]bool{NORMAL, SEQUENTIAL}, global.Match, ch)
	// Same as sequential, but doesn't wait for its requests to finish before starting the next URL requests
	case "cascade":
		ch.Progress.Total <- totalThreads(configs)
		runWorkers(configs, [2]bool{NORMAL, CASCADE}, global.Match, ch)

	// Sequential's behaviour, but cycles through the URLs requests for N times, N = largest amount of threads informed
	case "round-sequential":
		ch.Progress.Total <- maxThreads(configs) * len(configs)
		runWorkers(configs, [2]bool{ROUND, SEQUENTIAL}, global.Match, ch)
	// Cascade's behaviour, but cycles through the URLs requests for N times, N = largest amount of threads informed
	case "round-cascade":
		ch.Progress.Total <- maxThreads(configs) * len(configs)
		runWorkers(configs, [2]bool{ROUND, CASCADE}, global.Match, ch)

	// This is the default mode "flood", group all the requests and fire them at the exact same moment
	default:
		start := make(chan struct{})
		ch.Progress.Total <- totalThreads(configs)
		var wg sync.WaitGroup
		for _, w := range configs {
			for range w.Threads {
				wg.Go(func() { worker(start, w, global.Match, ch) })
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
