package log

import "fmt"

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

func Logger(logChannel chan LogMessage, userVerbosity int) {

	for {
		message := <-logChannel
		if shouldLog(userVerbosity, message.Verbosity) {
			fmt.Println(message)
		}
	}

}
