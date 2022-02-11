package bufpool

import "io"

func nativeCopy(dst io.Writer, src io.Reader) (performed bool, _ int64, _ error) {
	if src, isWriterTo := src.(io.WriterTo); isWriterTo {
		n, err := src.WriteTo(dst)
		return true, n, err
	}
	if dst, isReaderFrom := dst.(io.ReaderFrom); isReaderFrom {
		n, err := dst.ReadFrom(src)
		return true, n, err
	}
	return false, 0, nil
}
