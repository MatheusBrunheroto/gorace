package log

type Progress struct {
	Started   chan struct{}
	Total     chan int
	Sent      chan int
	Succeeded chan int
	Failed    chan int
	Finished  chan struct{}
}
type ProgressReader struct {
	Started   chan struct{}
	Total     <-chan int
	Sent      <-chan int
	Succeeded <-chan int
	Failed    <-chan int
	Finished  chan struct{}
}
type ProgressWriter struct {
	Started   chan struct{}
	Total     chan<- int
	Sent      chan<- int
	Succeeded chan<- int
	Failed    chan<- int
	Finished  chan struct{}
}

func (p Progress) Reader() ProgressReader {
	return ProgressReader{
		Started:   p.Started,
		Total:     p.Total,
		Sent:      p.Sent,
		Succeeded: p.Succeeded,
		Failed:    p.Failed,
		Finished:  p.Finished,
	}
}
func (p Progress) Writer() ProgressWriter {
	return ProgressWriter{
		Started:   p.Started,
		Total:     p.Total,
		Sent:      p.Sent,
		Succeeded: p.Succeeded,
		Failed:    p.Failed,
		Finished:  p.Finished,
	}
}
