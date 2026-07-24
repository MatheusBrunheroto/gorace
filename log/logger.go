package log

import (
	"fmt"
	"regexp"

	"atomicgo.dev/cursor"
)

/*

--verbosity 0 -> Mandatory Logs (Progress Bar, Fatal Errors)
--verbosity 1 -> All Above + Input Processing Feedback + String Matches
--verbosity 2 -> Visual Mode Feedback + (Which requests are starting, Cache Hits, status codes, content lenght)
--verbosity 3 -> full responses for each request that went right it's recommended to use --no-color and output via terminal for this
--verbosity 4 -> Unconditional Full ResponseBody string, it's recommended to use --no-color and output via terminal for this

*/
// VERBOSE 3 NAO TA COMPLETO IGUAL O 2
type Entry struct {
	Text       string
	isResponse bool
	Verbosity  int
}

var ansiRegex = regexp.MustCompile(`\x1b\[[0-9;]*m`)

func removeColor(text *string) {
	*text = ansiRegex.ReplaceAllString(*text, "")
}

// Determines whether an specific message can be displayed or not, based on --verbose input
func shouldLog(userVerbosity int, messageVerbosity int) bool {

	if userVerbosity == 1 && messageVerbosity <= 1 { // Only prints the simplest logs
		return true
	}
	if userVerbosity == 2 && messageVerbosity <= 2 { // Prints logs from verbose level 1 and 2
		return true
	}
	if userVerbosity == 3 && (messageVerbosity%2 != 0) { // Print all logs possible but 2 and 4
		return true
	}
	if userVerbosity == 4 && messageVerbosity != 2 { // Print all logs possible, but cuts 2 so don't send duplicates
		return true
	}
	return false // Don't Log the specified message

}

func Run(logChan chan Entry, userVerbosity *int, noColor *bool) {

	progress := cursor.NewArea()
	var lastProgress string

	for {

		log := <-logChan // Each log contains it's own verbosity, which is compared to the user's requested one

		// Mandatory Logs
		if log.Verbosity == 0 {
			lastProgress = log.Text
			removeColor(&log.Text)
			progress.Update(lastProgress)
			continue
		}

		if shouldLog(*userVerbosity, log.Verbosity) {
			progress.Clear()
			if *noColor {
				removeColor(&log.Text)
			}
			fmt.Println(log.Text)
			progress.Update(lastProgress)
		}

	}

}
