package testing

import (
	"bytes"
	"testing"

	cbg "github.com/whyrusleeping/cbor-gen"
)

// Compile-time proof that the opaque wrapper types implement the CBOR
// interfaces and therefore compose as fields of other generated types.
var (
	_ cbg.CBORMarshaler   = (*OpaqueString)(nil)
	_ cbg.CBORUnmarshaler = (*OpaqueString)(nil)
	_ cbg.CBORMarshaler   = (*OpaqueBytes)(nil)
	_ cbg.CBORUnmarshaler = (*OpaqueBytes)(nil)
	_ cbg.CBORMarshaler   = (*OpaqueDirect)(nil)
	_ cbg.CBORUnmarshaler = (*OpaqueDirect)(nil)
)

func mustMarshal(t *testing.T, m cbg.CBORMarshaler) []byte {
	t.Helper()
	var buf bytes.Buffer
	if err := m.MarshalCBOR(&buf); err != nil {
		t.Fatalf("marshal failed: %v", err)
	}
	return buf.Bytes()
}

func TestOpaqueStringRoundtrip(t *testing.T) {
	cases := []struct {
		name   string
		val    string
		golden []byte
	}{
		// non-empty value encodes as a CBOR text string
		{"value", "foo", []byte{0x63, 'f', 'o', 'o'}},
		// the zero value encodes as CBOR null (nullable)
		{"empty", "", []byte{0xf6}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			in, err := ParseOpaqueString(tc.val)
			if err != nil {
				t.Fatalf("ParseOpaqueString(%q): %v", tc.val, err)
			}

			enc := mustMarshal(t, &in)
			if !bytes.Equal(enc, tc.golden) {
				t.Fatalf("encoding mismatch: got %x want %x", enc, tc.golden)
			}

			var out OpaqueString
			if err := out.UnmarshalCBOR(bytes.NewReader(enc)); err != nil {
				t.Fatalf("unmarshal failed: %v", err)
			}
			if out.String() != tc.val {
				t.Fatalf("roundtrip mismatch: got %q want %q", out.String(), tc.val)
			}
		})
	}
}

// TestOpaqueStringValidatesOnDecode confirms the parse= hook rejects invalid
// input on the decode path: a text string carrying a disallowed space.
func TestOpaqueStringValidatesOnDecode(t *testing.T) {
	// CBOR text string "a b" — valid CBOR, but rejected by ParseOpaqueString.
	enc := []byte{0x63, 'a', ' ', 'b'}

	var out OpaqueString
	if err := out.UnmarshalCBOR(bytes.NewReader(enc)); err == nil {
		t.Fatal("expected decode to fail validation, got nil error")
	}
}

func TestOpaqueBytesRoundtrip(t *testing.T) {
	in, err := NewOpaqueBytes([]byte{0x01, 0x02, 0x03})
	if err != nil {
		t.Fatalf("NewOpaqueBytes: %v", err)
	}

	enc := mustMarshal(t, &in)
	// CBOR byte string of length 3.
	want := []byte{0x43, 0x01, 0x02, 0x03}
	if !bytes.Equal(enc, want) {
		t.Fatalf("encoding mismatch: got %x want %x", enc, want)
	}

	var out OpaqueBytes
	if err := out.UnmarshalCBOR(bytes.NewReader(enc)); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}
	if !bytes.Equal(out.Bytes(), in.Bytes()) {
		t.Fatalf("roundtrip mismatch: got %x want %x", out.Bytes(), in.Bytes())
	}
}

// TestOpaqueBytesValidatesOnDecode confirms the constructor rejects an empty
// payload on decode.
func TestOpaqueBytesValidatesOnDecode(t *testing.T) {
	enc := []byte{0x40} // empty CBOR byte string

	var out OpaqueBytes
	if err := out.UnmarshalCBOR(bytes.NewReader(enc)); err == nil {
		t.Fatal("expected decode of empty payload to fail, got nil error")
	}
}

// TestOpaqueBytesRejectsOverlong confirms the maxlen=8 bound is enforced on
// decode.
func TestOpaqueBytesRejectsOverlong(t *testing.T) {
	// CBOR byte string of length 9 (0x40 | 9 = 0x49) followed by 9 bytes.
	enc := append([]byte{0x49}, bytes.Repeat([]byte{0xaa}, 9)...)

	var out OpaqueBytes
	if err := out.UnmarshalCBOR(bytes.NewReader(enc)); err == nil {
		t.Fatal("expected decode of overlong payload to fail, got nil error")
	}
}

func TestOpaqueDirectRoundtrip(t *testing.T) {
	in := NewOpaqueDirect("bar")

	enc := mustMarshal(t, &in)
	want := []byte{0x63, 'b', 'a', 'r'}
	if !bytes.Equal(enc, want) {
		t.Fatalf("encoding mismatch: got %x want %x", enc, want)
	}

	var out OpaqueDirect
	if err := out.UnmarshalCBOR(bytes.NewReader(enc)); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}
	if out.String() != in.String() {
		t.Fatalf("roundtrip mismatch: got %q want %q", out.String(), in.String())
	}
}
