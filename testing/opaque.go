package testing

import (
	"fmt"
	"strings"
)

// The types in this file exercise the "opaque wrapper" pattern: a struct whose
// sole, unexported field holds a validated value and which encodes
// transparently as that field alone. This mirrors real-world types such as
// filecoin-project/go-address.Address, and ucantone's did.DID and
// command.Command, which previously had to hand-roll their CBOR codecs.

// OpaqueString is encoded as a transparent CBOR text string. It carries a
// validating constructor (parse=) and treats its zero value as CBOR null
// (nullable), mirroring did.DID / command.Command.
type OpaqueString struct {
	str string `cborgen:"transparent,parse=ParseOpaqueString,nullable"`
}

// ParseOpaqueString validates s and returns an OpaqueString. The empty string
// is permitted and represents the undefined value (encoded as CBOR null).
func ParseOpaqueString(s string) (OpaqueString, error) {
	if strings.ContainsRune(s, ' ') {
		return OpaqueString{}, fmt.Errorf("opaque string must not contain spaces: %q", s)
	}
	return OpaqueString{str: s}, nil
}

func (o OpaqueString) String() string { return o.str }

// OpaqueBytes is encoded as a transparent CBOR byte string with a validating
// constructor, mirroring go-address.Address. Its zero value is not nullable;
// an empty payload is rejected on decode by the constructor.
type OpaqueBytes struct {
	raw string `cborgen:"transparent,bytes,parse=NewOpaqueBytes,maxlen=8"`
}

// NewOpaqueBytes validates b and returns an OpaqueBytes. An empty payload is
// rejected.
func NewOpaqueBytes(b []byte) (OpaqueBytes, error) {
	if len(b) == 0 {
		return OpaqueBytes{}, fmt.Errorf("opaque bytes must not be empty")
	}
	return OpaqueBytes{raw: string(b)}, nil
}

func (o OpaqueBytes) Bytes() []byte { return []byte(o.raw) }

// OpaqueDirect is a transparent wrapper around an unexported string with no
// validating constructor; decoding assigns the field directly. It verifies
// that the unexported-field path works without a parse= hook.
type OpaqueDirect struct {
	val string `cborgen:"transparent"`
}

func NewOpaqueDirect(s string) OpaqueDirect { return OpaqueDirect{val: s} }

func (o OpaqueDirect) String() string { return o.val }
