package typegen

import (
	"bufio"
	"bytes"
	"io"
	"testing"
)

func TestPeeker(t *testing.T) {
	buf := bytes.NewBuffer([]byte{0, 1, 2, 3})
	p := peeker{reader: buf}
	n, err := p.Read(nil)
	if err != nil {
		t.Fatal(err)
	}
	if n != 0 {
		t.Fatal(err)
	}

	err = p.UnreadByte()
	if err != bufio.ErrInvalidUnreadByte {
		t.Fatal(err)
	}

	// read 2 bytes
	var out [2]byte
	n, err = p.Read(out[:])
	if err != nil {
		t.Fatal(err)
	}
	if n != 2 {
		t.Fatalf("expected 2 bytes, got %d", n)
	}
	if !bytes.Equal(out[:], []byte{0, 1}) {
		t.Fatalf("unexpected output")
	}

	// unread that last byte and read it again.
	err = p.UnreadByte()
	if err != nil {
		t.Fatal(err)
	}
	b, err := p.ReadByte()
	if err != nil {
		t.Fatal(err)
	}
	if b != 1 {
		t.Fatal("expected 1")
	}

	// unread that last byte then read 2
	err = p.UnreadByte()
	if err != nil {
		t.Fatal(err)
	}
	n, err = p.Read(out[:])
	if err != nil {
		t.Fatal(err)
	}
	if n != 2 {
		t.Fatalf("expected 2 bytes, got %d", n)
	}
	if !bytes.Equal(out[:], []byte{1, 2}) {
		t.Fatalf("unexpected output")
	}

	// read another byte
	b, err = p.ReadByte()
	if err != nil {
		t.Fatal(err)
	}
	if b != 3 {
		t.Fatal("expected 1")
	}

	// Should read eof at end.
	n, err = p.Read(out[:])
	if err != io.EOF {
		t.Fatal(err)
	}
	if n != 0 {
		t.Fatal("should have been at end")
	}
	// should unread eof
	err = p.UnreadByte()
	if err != nil {
		t.Fatal(err)
	}

	_, err = p.Read(nil)
	if err != nil {
		t.Fatal(err)
	}

	b, err = p.ReadByte()
	if err != nil {
		t.Fatal(err)
	}
	if b != 3 {
		t.Fatal("expected 1")
	}
}
