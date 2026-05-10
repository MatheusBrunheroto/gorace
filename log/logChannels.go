package log

type LogMessage struct {
	Message   string
	Verbosity int
}

type Progress struct {
	Total     chan int
	Sent      chan int
	Completed chan int
}
type ProgressReader struct {
	Total     <-chan int
	Sent      <-chan int
	Completed <-chan int
}
type ProgressWriter struct {
	Total     chan<- int
	Sent      chan<- int
	Completed chan<- int
}

func (p Progress) Reader() ProgressReader {
	return ProgressReader{
		Total:     p.Total,
		Sent:      p.Sent,
		Completed: p.Completed,
	}
}
func (p Progress) Writer() ProgressWriter {
	return ProgressWriter{
		Total:     p.Total,
		Sent:      p.Sent,
		Completed: p.Completed,
	}
}
