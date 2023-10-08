package expect

import (
	"context"
	"errors"
	"io"
)

type ReaderLease struct {
	reader io.Reader
	byte   chan byte
}

func NewReaderLease(reader io.Reader) *ReaderLease {
	readerLease := &ReaderLease{
		reader: reader,
		byte:   make(chan byte),
	}
	go func() {
		for {
			b := make([]byte, 1)
			_, err := readerLease.reader.Read(b)
			if err != nil {
				return
			}
			readerLease.byte <- b[0]
		}
	}()
	return readerLease
}

func (r *ReaderLease) NewCtxReader(ctx context.Context) *CtxReader {
	return &CtxReader{
		ctx:  ctx,
		byte: r.byte,
	}
}

type CtxReader struct {
	ctx  context.Context
	byte <-chan byte
}

func (c *CtxReader) Read(p []byte) (n int, err error) {
	select {
	case <-c.ctx.Done():
		return 0, io.EOF
	case b := <-c.byte:
		if len(p) < 1 {
			return 0, errors.New("can't read into 0 length type slice")
		}
		p[0] = b
		return 1, nil
	}
}
