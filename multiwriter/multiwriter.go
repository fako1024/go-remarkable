package multiwriter

import (
	"io"
	"sync"

	"github.com/fako1024/go-remarkable/compressor"
)

type multiWriter struct {
	done chan error
	io.Writer
}

type MultiWriter struct {
	writers []multiWriter
	c       *compressor.Compressor
	wg      sync.WaitGroup
	sync.Mutex
}

// New creates a writer that duplicates its writes to all the
// provided writers, similar to the Unix tee(1) command.
func New(writers ...io.Writer) *MultiWriter {

	mw := make([]multiWriter, 0, len(writers))
	for _, w := range writers {
		mw = append(mw, multiWriter{
			done:   make(chan error),
			Writer: w,
		})
	}

	return &MultiWriter{
		writers: mw,
		c:       compressor.New(),
	}
}

func (t *MultiWriter) remove(w io.Writer, err error) {
	for i := len(t.writers) - 1; i >= 0; i-- {
		if t.writers[i].Writer == w {
			t.writers = append(t.writers[:i], t.writers[i+1:]...)
			return
		}
	}
}

func (t *MultiWriter) Remove(w io.Writer, err error) {
	t.Lock()
	defer t.Unlock()

	t.remove(w, err)
}

func (t *MultiWriter) Append(w io.Writer) chan error {
	t.Lock()
	defer t.Unlock()

	done := make(chan error)
	t.writers = append(t.writers, multiWriter{
		done:   make(chan error),
		Writer: w,
	})

	return done
}

func (t *MultiWriter) N() int {
	return len(t.writers)
}

func (t *MultiWriter) Write(p []byte) (n int, err error) {
	t.Lock()
	defer t.Unlock()

	if len(t.writers) == 0 {
		return 0, io.ErrClosedPipe
	}

	l, data, err := t.c.Compress(p)
	if err != nil {
		return 0, err
	}

	t.wg.Add(len(t.writers))

	for _, w := range t.writers {
		go func(w multiWriter) {

			defer t.wg.Done()

			n, err = w.Write(data[:l])
			if err != nil {
				t.remove(w.Writer, err)
				return
			}
			if n != l {
				err = io.ErrShortWrite
				t.remove(w.Writer, err)
				return
			}
		}(w)
	}

	t.wg.Wait()

	return len(p), nil
}
