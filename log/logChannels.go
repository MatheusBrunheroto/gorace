package log

type Progress struct {
	Total     chan int
	Sent      chan int
	Completed chan int
	Finished  chan int
	Log       chan string
}

type LogMessage struct {
	Message   string
	Verbosity int
}
