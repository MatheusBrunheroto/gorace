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
}
type ProgressReader struct {
	Total     <-chan int
	Sent      <-chan int
	Succeeded <-chan int
	Failed    <-chan int
}
type ProgressWriter struct {
	Total     chan<- int
	Sent      chan<- int
	Succeeded chan<- int
	Failed    chan<- int
}

func (p Progress) Reader() ProgressReader {
	return ProgressReader{
		Total:     p.Total,
		Sent:      p.Sent,
		Succeeded: p.Succeeded,
		Failed:    p.Failed,
	}
}
func (p Progress) Writer() ProgressWriter {
	return ProgressWriter{
		Total:     p.Total,
		Sent:      p.Sent,
		Succeeded: p.Succeeded,
		Failed:    p.Failed,
	}
}
