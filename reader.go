package bufpool

import (
	"bufio"
	"io"
	"sync"
)

// Readers is a pool of bufio.Reader.
type Readers struct {
	pool *sync.Pool
}

// NewReaders creates a new pool.
// If bufSize>0, then new readers will be allocated with size of bufSize.
// If bufSize<=0, then new readers will be created with default size defined by bufio package.
func NewReaders(bufSize int) *Readers {
	newReader := func() any {
		return bufio.NewReader(nil)
	}
	if bufSize > 0 {
		newReader = func() any {
			return bufio.NewReaderSize(nil, bufSize)
		}
	}
	return &Readers{
		pool: &sync.Pool{
			New: newReader,
		},
	}
}

// New gets a new reader from the pool or allocates a new one.
func (readers *Readers) New(re io.Reader) *bufio.Reader {
	r := readers.pool.Get().(*bufio.Reader)
	r.Reset(re)
	return r
}

// Put resets and puts reader to the pool.
func (readers *Readers) Put(re *bufio.Reader) {
	re.Reset(nil)
	readers.pool.Put(re)
}

// Copy streams bytes (as io.Copy does) from src to dst using buffer from the pool.
// It can call src.WriteTo or dst.ReadFrom if provided.
func (readers *Readers) Copy(dst io.Writer, src io.Reader) (int64, error) {
	// don't allocate a new buffer if we can copy src to dst directly
	nativeCopied, n, err := nativeCopy(dst, src)
	if nativeCopied {
		return n, err
	}

	re := readers.New(src)
	defer readers.Put(re)

	return re.WriteTo(dst)
}
