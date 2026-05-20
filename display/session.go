package display

import "gorace/log"

type Session struct {
	Draw     chan string
	Ready    chan struct{}
	Finished chan struct{}
	Progress log.ProgressReader
}

func NewSession(progress log.ProgressReader) Session {
	return Session{
		Draw:     make(chan string),
		Ready:    make(chan struct{}),
		Finished: make(chan struct{}),
		Progress: progress,
	}
}
