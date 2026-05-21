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
	_ cbg.CBORMarshaler   = (*OpaqueBytesDirect)(nil)
	_ cbg.CBORUnmarshaler = (*OpaqueBytesDirect)(nil)
	_ cbg.CBORMarshaler   = (*OpaqueNullableDirect)(nil)
	_ cbg.CBORUnmarshaler = (*OpaqueNullableDirect)(nil)
	_ cbg.CBORMarshaler   = (*OpaqueBytesNullable)(nil)
	_ cbg.CBORUnmarshaler = (*OpaqueBytesNullable)(nil)
	_ cbg.CBORMarshaler   = (*OpaqueContainer)(nil)
	_ cbg.CBORUnmarshaler = (*OpaqueContainer)(nil)
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

// TestOpaqueBytesDirectRoundtrip exercises the `bytes` codegen path without a
// parse= hook: the byte string is assigned to the field directly on decode.
func TestOpaqueBytesDirectRoundtrip(t *testing.T) {
	in := NewOpaqueBytesDirect([]byte{0x0a, 0x0b})

	enc := mustMarshal(t, &in)
	want := []byte{0x42, 0x0a, 0x0b} // CBOR byte string of length 2
	if !bytes.Equal(enc, want) {
		t.Fatalf("encoding mismatch: got %x want %x", enc, want)
	}

	var out OpaqueBytesDirect
	if err := out.UnmarshalCBOR(bytes.NewReader(enc)); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}
	if !bytes.Equal(out.Bytes(), in.Bytes()) {
		t.Fatalf("roundtrip mismatch: got %x want %x", out.Bytes(), in.Bytes())
	}
}

// TestOpaqueNullableDirectRoundtrip exercises the nullable text-string path with
// direct field assignment (no parse= hook): a non-empty value encodes as a text
// string, the zero value as CBOR null.
func TestOpaqueNullableDirectRoundtrip(t *testing.T) {
	cases := []struct {
		name   string
		val    string
		golden []byte
	}{
		{"value", "foo", []byte{0x63, 'f', 'o', 'o'}},
		{"empty", "", []byte{0xf6}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			in := NewOpaqueNullableDirect(tc.val)

			enc := mustMarshal(t, &in)
			if !bytes.Equal(enc, tc.golden) {
				t.Fatalf("encoding mismatch: got %x want %x", enc, tc.golden)
			}

			var out OpaqueNullableDirect
			if err := out.UnmarshalCBOR(bytes.NewReader(enc)); err != nil {
				t.Fatalf("unmarshal failed: %v", err)
			}
			if out.String() != tc.val {
				t.Fatalf("roundtrip mismatch: got %q want %q", out.String(), tc.val)
			}
		})
	}
}

// TestOpaqueBytesNullableRoundtrip exercises the combined bytes+nullable path: a
// non-empty value encodes as a byte string, the zero value as CBOR null.
func TestOpaqueBytesNullableRoundtrip(t *testing.T) {
	cases := []struct {
		name   string
		val    []byte
		golden []byte
	}{
		{"value", []byte{0x01, 0x02}, []byte{0x42, 0x01, 0x02}},
		{"empty", nil, []byte{0xf6}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			in := NewOpaqueBytesNullable(tc.val)

			enc := mustMarshal(t, &in)
			if !bytes.Equal(enc, tc.golden) {
				t.Fatalf("encoding mismatch: got %x want %x", enc, tc.golden)
			}

			var out OpaqueBytesNullable
			if err := out.UnmarshalCBOR(bytes.NewReader(enc)); err != nil {
				t.Fatalf("unmarshal failed: %v", err)
			}
			if !bytes.Equal(out.Bytes(), in.Bytes()) {
				t.Fatalf("roundtrip mismatch: got %x want %x", out.Bytes(), in.Bytes())
			}
		})
	}
}

// TestOpaqueRejectsNullWhenNotNullable confirms a CBOR null is rejected when
// decoded into a non-nullable opaque field, rather than silently zeroing it.
func TestOpaqueRejectsNullWhenNotNullable(t *testing.T) {
	null := []byte{0xf6}

	t.Run("bytes", func(t *testing.T) {
		var out OpaqueBytes
		if err := out.UnmarshalCBOR(bytes.NewReader(null)); err == nil {
			t.Fatal("expected decode of null into non-nullable bytes field to fail")
		}
	})

	t.Run("direct", func(t *testing.T) {
		var out OpaqueDirect
		if err := out.UnmarshalCBOR(bytes.NewReader(null)); err == nil {
			t.Fatal("expected decode of null into non-nullable string field to fail")
		}
	})
}

// TestOpaqueBytesMarshalRejectsOverlong confirms the maxlen=8 bound is enforced
// on the encode path, not only on decode. The constructor permits the payload
// (it only rejects empty), so the bound must be caught by MarshalCBOR.
func TestOpaqueBytesMarshalRejectsOverlong(t *testing.T) {
	in, err := NewOpaqueBytes(bytes.Repeat([]byte{0xaa}, 9))
	if err != nil {
		t.Fatalf("NewOpaqueBytes: %v", err)
	}

	var buf bytes.Buffer
	if err := in.MarshalCBOR(&buf); err == nil {
		t.Fatal("expected marshal of overlong payload to fail, got nil error")
	}
}

// TestOpaqueContainerRoundtrip confirms opaque wrappers compose as fields of
// another generated type and round-trip through it.
func TestOpaqueContainerRoundtrip(t *testing.T) {
	name, err := ParseOpaqueString("hello")
	if err != nil {
		t.Fatalf("ParseOpaqueString: %v", err)
	}
	data, err := NewOpaqueBytes([]byte{0x01, 0x02, 0x03})
	if err != nil {
		t.Fatalf("NewOpaqueBytes: %v", err)
	}
	in := OpaqueContainer{Name: name, Data: data}

	enc := mustMarshal(t, &in)

	var out OpaqueContainer
	if err := out.UnmarshalCBOR(bytes.NewReader(enc)); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}
	if out.Name.String() != in.Name.String() {
		t.Fatalf("Name mismatch: got %q want %q", out.Name.String(), in.Name.String())
	}
	if !bytes.Equal(out.Data.Bytes(), in.Data.Bytes()) {
		t.Fatalf("Data mismatch: got %x want %x", out.Data.Bytes(), in.Data.Bytes())
	}
}
