package typegen

import (
	"bytes"
	"testing"
)

func TestDeferredMaxLengthSingle(t *testing.T) {
	var header bytes.Buffer
	if err := WriteMajorTypeHeader(&header, MajByteString, ByteArrayMaxLen+1); err != nil {
		t.Fatal("failed to write header")
	}

	var deferred Deferred
	err := deferred.UnmarshalCBOR(&header)
	if err != maxLengthError {
		t.Fatal("deferred: allowed more than the maximum allocation supported")
	}
}

func TestDeferredMaxLengthRecursive(t *testing.T) {
	var header bytes.Buffer
	for i := 0; i < MaxLength+1; i++ {
		if err := WriteMajorTypeHeader(&header, MajTag, 0); err != nil {
			t.Fatal("failed to write header")
		}
	}

	var deferred Deferred
	err := deferred.UnmarshalCBOR(&header)
	if err != maxLengthError {
		t.Fatal("deferred: allowed more than the maximum number of elements")
	}
}
