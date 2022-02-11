package bufpool

import (
	"bytes"
	"sync"
)

// Buffers is a pool of *bytes.Buffers.
type Buffers struct {
	pool *sync.Pool
}

// NewBuffers creates a new buffers pool.
// If bufSize>0, then it will be used as a bootstrap size for new buffers.
// If bufSize<=0, then new buffers will be empty.
func NewBuffers(bufSize int) *Buffers {
	newBuf := func() any {
		return &bytes.Buffer{}
	}
	if bufSize > 0 {
		newBuf = func() any {
			buf := &bytes.Buffer{}
			buf.Grow(bufSize)
			return buf
		}
	}
	return &Buffers{
		pool: &sync.Pool{
			New: newBuf,
		},
	}
}

// New returns a new bytes buffer from the pool or allocates a new one.
func (buffers *Buffers) New() *bytes.Buffer {
	return buffers.pool.Get().(*bytes.Buffer)
}

// NewString creates a new buffer from provided string.
// It's an analog of the bytes.NewBufferString.
func (buffers *Buffers) NewString(str string) *bytes.Buffer {
	buf := buffers.New()
	_, _ = buf.WriteString(str)
	return buf
}

// NewBytes creates a new buffer from provided byte slice.
// Is's an analog of bytes.NewBuffer.
func (buffers *Buffers) NewBytes(bb []byte) *bytes.Buffer {
	buf := buffers.New()
	_, _ = buf.Write(bb)
	return buf
}

// Put resets buffer and returns it to the pool.
func (buffers *Buffers) Put(buf *bytes.Buffer) {
	buf.Reset()
	buffers.pool.Put(buf)
}
