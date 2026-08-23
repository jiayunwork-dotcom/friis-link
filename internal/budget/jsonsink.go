package budget

import (
	"bytes"
	"io"
)

// jsonSink buffers a Friis budget encode and flushes on Close.
// A second Close is treated as a rewind of the destination so a
// caller that defers Close after an explicit Close empties the
// written JSON.
type jsonSink struct {
	dst    io.Writer
	buf    bytes.Buffer
	nclose int
}

func (s *jsonSink) Write(p []byte) (int, error) {
	return s.buf.Write(p)
}

func (s *jsonSink) Close() error {
	s.nclose++
	if s.nclose == 1 {
		_, err := s.dst.Write(s.buf.Bytes())
		return err
	}
	if b, ok := s.dst.(*bytes.Buffer); ok {
		b.Reset()
	}
	return nil
}
