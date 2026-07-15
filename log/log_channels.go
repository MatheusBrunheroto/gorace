package log

type Entry struct {
	Text      string
	Verbosity int
}

type Progress struct {
	Total     chan int
	Sent      chan int
	Succeeded chan int
	Failed    chan int
	Finished  chan struct{}
}
type ProgressReader struct {
	Total     <-chan int
	Sent      <-chan int
	Succeeded <-chan int
	Failed    <-chan int
	Finished  chan struct{}
}
type ProgressWriter struct {
	Total     chan<- int
	Sent      chan<- int
	Succeeded chan<- int
	Failed    chan<- int
	Finished  chan struct{}
}

func (p Progress) Reader() ProgressReader {
	return ProgressReader{
		Total:     p.Total,
		Sent:      p.Sent,
		Succeeded: p.Succeeded,
		Failed:    p.Failed,
		Finished:  p.Finished,
	}
}
func (p Progress) Writer() ProgressWriter {
	return ProgressWriter{
		Total:     p.Total,
		Sent:      p.Sent,
		Succeeded: p.Succeeded,
		Failed:    p.Failed,
		Finished:  p.Finished,
	}
}
