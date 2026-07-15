package log

import (
	"fmt"

	"atomicgo.dev/cursor"
)

// Determines whether an specific message can be displayed or not, based on --verbose input
func shouldLog(userVerbosity int, messageVerbosity int) bool {

	if userVerbosity == 1 { // Only prints the simplest logs
		return true
	}
	if userVerbosity == 2 && messageVerbosity <= 2 { // Prints logs from verbose level 1 and 2
		return true
	}
	if userVerbosity == 3 && messageVerbosity <= 3 { // Print all logs possible
		return true
	}
	return false // userVerbosity == 0

}

func Run(logChan chan Entry, userVerbosity *int) {

	progress := cursor.NewArea()
	var lastProgress string

	for {

		log := <-logChan // Each log contains it's own verbosity, which is compared to the user's requested one

		if log.Verbosity == 0 {
			lastProgress = log.Text
			progress.Update(lastProgress)
			continue
		}

		if shouldLog(*userVerbosity, log.Verbosity) {
			progress.Clear()
			fmt.Println(log.Text)
			progress.Update(lastProgress)
		}

	}

}
