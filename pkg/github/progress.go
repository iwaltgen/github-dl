package github

import (
	"math"

	"github.com/reactivex/rxgo/v2"
)

// Progress is progress info of work.
type Progress interface {
	Percentage() float64
}

// DownloadProgress is the current downloaded size of a file.
type DownloadProgress struct {
	Total    int64
	Received int64
}

// Percentage is downloaded size over total size.
func (p *DownloadProgress) Percentage() float64 {
	return math.Round(float64(p.Received)/float64(p.Total)*10000) / 100
}

// WriteCounter counts the number of bytes written to it. It implements to the io.Writer interface
// and we can pass this into io.TeeReader() which will report progress on each write cycle.
type WriteCounter struct {
	Total   int64
	Written int64
	ch      chan<- rxgo.Item
}

// NewWriteCounter creates WriteCounter
func NewWriteCounter(ch chan<- rxgo.Item, total int64) *WriteCounter {
	return &WriteCounter{
		ch:    ch,
		Total: total,
	}
}

func (w *WriteCounter) Write(p []byte) (int, error) {
	n := len(p)
	w.Written += int64(n)
	w.ch <- rxgo.Of(&DownloadProgress{
		Total:    w.Total,
		Received: w.Written,
	})
	return n, nil
}
