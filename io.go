package typegen

import (
	"fmt"
	"io"

	cid "github.com/ipfs/go-cid"
)

var (
	_ io.Reader      = (*CborReader)(nil)
	_ io.ByteScanner = (*CborReader)(nil)
)

type CborReader struct {
	r    BytePeeker
	hbuf [maxHeaderSize]byte
}

func NewCborReader(r io.Reader) *CborReader {
	if r, ok := r.(*CborReader); ok {
		return r
	}

	return &CborReader{
		r: GetPeeker(r),
	}
}

func (cr *CborReader) Read(p []byte) (n int, err error) {
	return cr.r.Read(p)
}

func (cr *CborReader) ReadByte() (byte, error) {
	return cr.r.ReadByte()
}

func (cr *CborReader) UnreadByte() error {
	return cr.r.UnreadByte()
}

func (cr *CborReader) ReadHeader() (byte, uint64, error) {
	return CborReadHeaderBuf(cr.r, cr.hbuf[:])
}

func (cr *CborReader) SetReader(r io.Reader) {
	cr.r = GetPeeker(r)
}

func (cr *CborReader) ReadCid(scratchBuf []byte) (cid.Cid, error) {
	maj, extra, err := cr.ReadHeader()
	if err != nil {
		if err == io.EOF {
			err = io.ErrUnexpectedEOF
		}
		return cid.Undef, err
	}

	if maj != MajTag {
		return cid.Undef, fmt.Errorf("expected cbor type 'tag' in input")
	}

	if extra != 42 {
		return cid.Undef, fmt.Errorf("expected tag 42")
	}

	if extra > 512 {
		return cid.Undef, fmt.Errorf("header size too big for a cid")
	}

	_, err = io.ReadAtLeast(cr, scratchBuf[:extra], int(extra))
	if err != nil {
		if err == io.EOF {
			err = io.ErrUnexpectedEOF
		}
		return cid.Undef, err
	}

	return bufToCid(scratchBuf[:extra])
}

var (
	_ io.Writer       = (*CborWriter)(nil)
	_ io.StringWriter = (*CborWriter)(nil)
)

type CborWriter struct {
	w    io.Writer
	hbuf [maxHeaderSize]byte

	sw io.StringWriter
}

func NewCborWriter(w io.Writer) *CborWriter {
	if w, ok := w.(*CborWriter); ok {
		return w
	}

	cw := &CborWriter{
		w: w,
	}

	if sw, ok := w.(io.StringWriter); ok {
		cw.sw = sw
	}

	return cw
}

func (cw *CborWriter) SetWriter(w io.Writer) {
	cw.w = w
	if sw, ok := w.(io.StringWriter); ok {
		cw.sw = sw
	} else {
		cw.sw = nil
	}
}

func (cw *CborWriter) Write(p []byte) (n int, err error) {
	return cw.w.Write(p)
}

func (cw *CborWriter) WriteMajorTypeHeader(t byte, l uint64) error {
	return WriteMajorTypeHeaderBuf(cw.hbuf[:], cw.w, t, l)
}

func (cw *CborWriter) CborWriteHeader(t byte, l uint64) error {
	return WriteMajorTypeHeaderBuf(cw.hbuf[:], cw.w, t, l)
}

func (cw *CborWriter) WriteString(s string) (int, error) {
	if cw.sw != nil {
		return cw.sw.WriteString(s)
	}
	return cw.w.Write([]byte(s))
}
