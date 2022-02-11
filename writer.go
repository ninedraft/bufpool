package bufpool

import (
	"bufio"
	"io"
	"sync"
)

// Writers is bufio.Writer pool.
type Writers struct {
	pool *sync.Pool
}

// NewWriters creates a new pool of bufio.Writer.
// If bufSize>0, then new writers will be allocated with size bufSize.
// If bufSize<=0, then new writers will be created with a default size defined by bufio.
func NewWriters(bufSize int) *Writers {
	newReader := func() any {
		return bufio.NewWriter(nil)
	}
	if bufSize > 0 {
		newReader = func() any {
			return bufio.NewWriterSize(nil, bufSize)
		}
	}
	return &Writers{
		pool: &sync.Pool{
			New: newReader,
		},
	}
}

// New gets a new writer from the pool or allocates a new one.
func (writers *Writers) New(wr io.Writer) *bufio.Writer {
	w := writers.pool.Get().(*bufio.Writer)
	w.Reset(wr)
	return w
}

// Put resets writer and puts it in the pool.
func (writers *Writers) Put(re *bufio.Writer) {
	re.Reset(nil)
	writers.pool.Put(re)
}

// Copy streams bytes from src to dst (as io.Copy does) using bufio.Writer from the pool.
// It can use src.WriteTo or dst.ReadFrom if provided.
func (writers *Writers) Copy(dst io.Writer, src io.Reader) (int64, error) {
	// don't allocate a new buffer if we can copy src to dst directly
	nativeCopied, n, err := nativeCopy(dst, src)
	if nativeCopied {
		return n, err
	}

	wr := writers.New(dst)
	defer writers.Put(wr)

	n, err = wr.ReadFrom(src)
	if err != nil {
		return n, err
	}
	if err := wr.Flush(); err != nil {
		return n, err
	}
	return n, nil
}
