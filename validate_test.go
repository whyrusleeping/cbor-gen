package typegen

import (
	"bytes"
	"testing"
)

func TestValidateShort(t *testing.T) {
	var buf bytes.Buffer
	if err := WriteMajorTypeHeader(&buf, MajByteString, 100); err != nil {
		t.Fatal("failed to write header")
	}

	if err := ValidateCBOR(buf.Bytes()); err == nil {
		t.Fatal("expected an error checking truncated cbor")
	}
}

func TestValidateDouble(t *testing.T) {
	var buf bytes.Buffer
	if err := WriteBool(&buf, false); err != nil {
		t.Fatal(err)
	}
	if err := WriteBool(&buf, false); err != nil {
		t.Fatal(err)
	}

	if err := ValidateCBOR(buf.Bytes()); err == nil {
		t.Fatal("expected an error checking cbor with two objects")
	}
}

func TestValidate(t *testing.T) {
	var buf bytes.Buffer
	if err := WriteBool(&buf, false); err != nil {
		t.Fatal(err)
	}

	if err := ValidateCBOR(buf.Bytes()); err != nil {
		t.Fatal(err)
	}
}
