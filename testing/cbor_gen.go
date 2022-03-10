// Code generated by github.com/whyrusleeping/cbor-gen. DO NOT EDIT.

package testing

import (
	"fmt"
	"io"
	"math"
	"sort"

	cid "github.com/ipfs/go-cid"
	cbg "github.com/whyrusleeping/cbor-gen"
	xerrors "golang.org/x/xerrors"
)

var _ = xerrors.Errorf
var _ = cid.Undef
var _ = math.E
var _ = sort.Sort

var lengthBufSignedArray = []byte{129}

func (t *SignedArray) MarshalCBOR(w io.Writer) error {
	if t == nil {
		_, err := w.Write(cbg.CborNull)
		return err
	}
	if _, err := w.Write(lengthBufSignedArray); err != nil {
		return err
	}

	scratch := make([]byte, 9)

	// t.Signed ([]uint64) (slice)
	if len(t.Signed) > cbg.MaxLength {
		return xerrors.Errorf("Slice value in field t.Signed was too long")
	}

	if err := cbg.WriteMajorTypeHeaderBuf(scratch, w, cbg.MajArray, uint64(len(t.Signed))); err != nil {
		return err
	}
	for _, v := range t.Signed {
		if err := cbg.CborWriteHeader(w, cbg.MajUnsignedInt, uint64(v)); err != nil {
			return err
		}
	}
	return nil
}

func (t *SignedArray) UnmarshalCBOR(r io.Reader) (err error) {
	*t = SignedArray{}

	br := cbg.GetPeeker(r)
	scratch := make([]byte, 8)

	hasReadOnce := false
	defer func() {
		if err == io.EOF && hasReadOnce {
			err = io.ErrUnexpectedEOF
		}
	}()
	maj, extra, err := cbg.CborReadHeaderBuf(br, scratch)
	if err != nil {
		return err
	}
	hasReadOnce = true

	if maj != cbg.MajArray {
		return fmt.Errorf("cbor input should be of type array")
	}

	if extra != 1 {
		return fmt.Errorf("cbor input had wrong number of fields")
	}

	// t.Signed ([]uint64) (slice)

	maj, extra, err = cbg.CborReadHeaderBuf(br, scratch)
	if err != nil {
		return err
	}

	if extra > cbg.MaxLength {
		return fmt.Errorf("t.Signed: array too large (%d)", extra)
	}

	if maj != cbg.MajArray {
		return fmt.Errorf("expected cbor array")
	}

	if extra > 0 {
		t.Signed = make([]uint64, extra)
	}

	for i := 0; i < int(extra); i++ {

		maj, val, err := cbg.CborReadHeaderBuf(br, scratch)
		if err != nil {
			return xerrors.Errorf("failed to read uint64 for t.Signed slice: %w", err)
		}

		if maj != cbg.MajUnsignedInt {
			return xerrors.Errorf("value read for array t.Signed was not a uint, instead got %d", maj)
		}

		t.Signed[i] = uint64(val)
	}

	return nil
}

var lengthBufSimpleTypeOne = []byte{133}

func (t *SimpleTypeOne) MarshalCBOR(w io.Writer) error {
	if t == nil {
		_, err := w.Write(cbg.CborNull)
		return err
	}
	if _, err := w.Write(lengthBufSimpleTypeOne); err != nil {
		return err
	}

	scratch := make([]byte, 9)

	// t.Foo (string) (string)
	if len(t.Foo) > cbg.MaxLength {
		return xerrors.Errorf("Value in field t.Foo was too long")
	}

	if err := cbg.WriteMajorTypeHeaderBuf(scratch, w, cbg.MajTextString, uint64(len(t.Foo))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string(t.Foo)); err != nil {
		return err
	}

	// t.Value (uint64) (uint64)

	if err := cbg.WriteMajorTypeHeaderBuf(scratch, w, cbg.MajUnsignedInt, uint64(t.Value)); err != nil {
		return err
	}

	// t.Binary ([]uint8) (slice)
	if len(t.Binary) > cbg.ByteArrayMaxLen {
		return xerrors.Errorf("Byte array in field t.Binary was too long")
	}

	if err := cbg.WriteMajorTypeHeaderBuf(scratch, w, cbg.MajByteString, uint64(len(t.Binary))); err != nil {
		return err
	}

	if _, err := w.Write(t.Binary[:]); err != nil {
		return err
	}

	// t.Signed (int64) (int64)
	if t.Signed >= 0 {
		if err := cbg.WriteMajorTypeHeaderBuf(scratch, w, cbg.MajUnsignedInt, uint64(t.Signed)); err != nil {
			return err
		}
	} else {
		if err := cbg.WriteMajorTypeHeaderBuf(scratch, w, cbg.MajNegativeInt, uint64(-t.Signed-1)); err != nil {
			return err
		}
	}

	// t.NString (testing.NamedString) (string)
	if len(t.NString) > cbg.MaxLength {
		return xerrors.Errorf("Value in field t.NString was too long")
	}

	if err := cbg.WriteMajorTypeHeaderBuf(scratch, w, cbg.MajTextString, uint64(len(t.NString))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string(t.NString)); err != nil {
		return err
	}
	return nil
}

func (t *SimpleTypeOne) UnmarshalCBOR(r io.Reader) (err error) {
	*t = SimpleTypeOne{}

	br := cbg.GetPeeker(r)
	scratch := make([]byte, 8)

	hasReadOnce := false
	defer func() {
		if err == io.EOF && hasReadOnce {
			err = io.ErrUnexpectedEOF
		}
	}()
	maj, extra, err := cbg.CborReadHeaderBuf(br, scratch)
	if err != nil {
		return err
	}
	hasReadOnce = true

	if maj != cbg.MajArray {
		return fmt.Errorf("cbor input should be of type array")
	}

	if extra != 5 {
		return fmt.Errorf("cbor input had wrong number of fields")
	}

	// t.Foo (string) (string)

	{
		sval, err := cbg.ReadStringBuf(br, scratch)
		if err != nil {
			return err
		}

		t.Foo = string(sval)
	}
	// t.Value (uint64) (uint64)

	{

		maj, extra, err = cbg.CborReadHeaderBuf(br, scratch)
		if err != nil {
			return err
		}
		if maj != cbg.MajUnsignedInt {
			return fmt.Errorf("wrong type for uint64 field")
		}
		t.Value = uint64(extra)

	}
	// t.Binary ([]uint8) (slice)

	maj, extra, err = cbg.CborReadHeaderBuf(br, scratch)
	if err != nil {
		return err
	}

	if extra > cbg.ByteArrayMaxLen {
		return fmt.Errorf("t.Binary: byte array too large (%d)", extra)
	}
	if maj != cbg.MajByteString {
		return fmt.Errorf("expected byte array")
	}

	if extra > 0 {
		t.Binary = make([]uint8, extra)
	}

	if _, err := io.ReadFull(br, t.Binary[:]); err != nil {
		return err
	}
	// t.Signed (int64) (int64)
	{
		maj, extra, err := cbg.CborReadHeaderBuf(br, scratch)
		var extraI int64
		if err != nil {
			return err
		}
		switch maj {
		case cbg.MajUnsignedInt:
			extraI = int64(extra)
			if extraI < 0 {
				return fmt.Errorf("int64 positive overflow")
			}
		case cbg.MajNegativeInt:
			extraI = int64(extra)
			if extraI < 0 {
				return fmt.Errorf("int64 negative oveflow")
			}
			extraI = -1 - extraI
		default:
			return fmt.Errorf("wrong type for int64 field: %d", maj)
		}

		t.Signed = int64(extraI)
	}
	// t.NString (testing.NamedString) (string)

	{
		sval, err := cbg.ReadStringBuf(br, scratch)
		if err != nil {
			return err
		}

		t.NString = NamedString(sval)
	}
	return nil
}

var lengthBufSimpleTypeTwo = []byte{137}

func (t *SimpleTypeTwo) MarshalCBOR(w io.Writer) error {
	if t == nil {
		_, err := w.Write(cbg.CborNull)
		return err
	}
	if _, err := w.Write(lengthBufSimpleTypeTwo); err != nil {
		return err
	}

	scratch := make([]byte, 9)

	// t.Stuff (testing.SimpleTypeTwo) (struct)
	if err := t.Stuff.MarshalCBOR(w); err != nil {
		return err
	}

	// t.Others ([]uint64) (slice)
	if len(t.Others) > cbg.MaxLength {
		return xerrors.Errorf("Slice value in field t.Others was too long")
	}

	if err := cbg.WriteMajorTypeHeaderBuf(scratch, w, cbg.MajArray, uint64(len(t.Others))); err != nil {
		return err
	}
	for _, v := range t.Others {
		if err := cbg.CborWriteHeader(w, cbg.MajUnsignedInt, uint64(v)); err != nil {
			return err
		}
	}

	// t.SignedOthers ([]int64) (slice)
	if len(t.SignedOthers) > cbg.MaxLength {
		return xerrors.Errorf("Slice value in field t.SignedOthers was too long")
	}

	if err := cbg.WriteMajorTypeHeaderBuf(scratch, w, cbg.MajArray, uint64(len(t.SignedOthers))); err != nil {
		return err
	}
	for _, v := range t.SignedOthers {
		if v >= 0 {
			if err := cbg.WriteMajorTypeHeaderBuf(scratch, w, cbg.MajUnsignedInt, uint64(v)); err != nil {
				return err
			}
		} else {
			if err := cbg.WriteMajorTypeHeaderBuf(scratch, w, cbg.MajNegativeInt, uint64(-v-1)); err != nil {
				return err
			}
		}
	}

	// t.Test ([][]uint8) (slice)
	if len(t.Test) > cbg.MaxLength {
		return xerrors.Errorf("Slice value in field t.Test was too long")
	}

	if err := cbg.WriteMajorTypeHeaderBuf(scratch, w, cbg.MajArray, uint64(len(t.Test))); err != nil {
		return err
	}
	for _, v := range t.Test {
		if len(v) > cbg.ByteArrayMaxLen {
			return xerrors.Errorf("Byte array in field v was too long")
		}

		if err := cbg.WriteMajorTypeHeaderBuf(scratch, w, cbg.MajByteString, uint64(len(v))); err != nil {
			return err
		}

		if _, err := w.Write(v[:]); err != nil {
			return err
		}
	}

	// t.Dog (string) (string)
	if len(t.Dog) > cbg.MaxLength {
		return xerrors.Errorf("Value in field t.Dog was too long")
	}

	if err := cbg.WriteMajorTypeHeaderBuf(scratch, w, cbg.MajTextString, uint64(len(t.Dog))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string(t.Dog)); err != nil {
		return err
	}

	// t.Numbers ([]testing.NamedNumber) (slice)
	if len(t.Numbers) > cbg.MaxLength {
		return xerrors.Errorf("Slice value in field t.Numbers was too long")
	}

	if err := cbg.WriteMajorTypeHeaderBuf(scratch, w, cbg.MajArray, uint64(len(t.Numbers))); err != nil {
		return err
	}
	for _, v := range t.Numbers {
		if err := cbg.CborWriteHeader(w, cbg.MajUnsignedInt, uint64(v)); err != nil {
			return err
		}
	}

	// t.Pizza (uint64) (uint64)

	if t.Pizza == nil {
		if _, err := w.Write(cbg.CborNull); err != nil {
			return err
		}
	} else {
		if err := cbg.WriteMajorTypeHeaderBuf(scratch, w, cbg.MajUnsignedInt, uint64(*t.Pizza)); err != nil {
			return err
		}
	}

	// t.PointyPizza (testing.NamedNumber) (uint64)

	if t.PointyPizza == nil {
		if _, err := w.Write(cbg.CborNull); err != nil {
			return err
		}
	} else {
		if err := cbg.WriteMajorTypeHeaderBuf(scratch, w, cbg.MajUnsignedInt, uint64(*t.PointyPizza)); err != nil {
			return err
		}
	}

	// t.Arrrrrghay ([3]testing.SimpleTypeOne) (array)
	if len(t.Arrrrrghay) > cbg.MaxLength {
		return xerrors.Errorf("Slice value in field t.Arrrrrghay was too long")
	}

	if err := cbg.WriteMajorTypeHeaderBuf(scratch, w, cbg.MajArray, uint64(len(t.Arrrrrghay))); err != nil {
		return err
	}
	for _, v := range t.Arrrrrghay {
		if err := v.MarshalCBOR(w); err != nil {
			return err
		}
	}
	return nil
}

func (t *SimpleTypeTwo) UnmarshalCBOR(r io.Reader) (err error) {
	*t = SimpleTypeTwo{}

	br := cbg.GetPeeker(r)
	scratch := make([]byte, 8)

	hasReadOnce := false
	defer func() {
		if err == io.EOF && hasReadOnce {
			err = io.ErrUnexpectedEOF
		}
	}()
	maj, extra, err := cbg.CborReadHeaderBuf(br, scratch)
	if err != nil {
		return err
	}
	hasReadOnce = true

	if maj != cbg.MajArray {
		return fmt.Errorf("cbor input should be of type array")
	}

	if extra != 9 {
		return fmt.Errorf("cbor input had wrong number of fields")
	}

	// t.Stuff (testing.SimpleTypeTwo) (struct)

	{

		b, err := br.ReadByte()
		if err != nil {
			return err
		}
		if b != cbg.CborNull[0] {
			if err := br.UnreadByte(); err != nil {
				return err
			}
			t.Stuff = new(SimpleTypeTwo)
			if err := t.Stuff.UnmarshalCBOR(br); err != nil {
				return xerrors.Errorf("unmarshaling t.Stuff pointer: %w", err)
			}
		}

	}
	// t.Others ([]uint64) (slice)

	maj, extra, err = cbg.CborReadHeaderBuf(br, scratch)
	if err != nil {
		return err
	}

	if extra > cbg.MaxLength {
		return fmt.Errorf("t.Others: array too large (%d)", extra)
	}

	if maj != cbg.MajArray {
		return fmt.Errorf("expected cbor array")
	}

	if extra > 0 {
		t.Others = make([]uint64, extra)
	}

	for i := 0; i < int(extra); i++ {

		maj, val, err := cbg.CborReadHeaderBuf(br, scratch)
		if err != nil {
			return xerrors.Errorf("failed to read uint64 for t.Others slice: %w", err)
		}

		if maj != cbg.MajUnsignedInt {
			return xerrors.Errorf("value read for array t.Others was not a uint, instead got %d", maj)
		}

		t.Others[i] = uint64(val)
	}

	// t.SignedOthers ([]int64) (slice)

	maj, extra, err = cbg.CborReadHeaderBuf(br, scratch)
	if err != nil {
		return err
	}

	if extra > cbg.MaxLength {
		return fmt.Errorf("t.SignedOthers: array too large (%d)", extra)
	}

	if maj != cbg.MajArray {
		return fmt.Errorf("expected cbor array")
	}

	if extra > 0 {
		t.SignedOthers = make([]int64, extra)
	}

	for i := 0; i < int(extra); i++ {
		{
			maj, extra, err := cbg.CborReadHeaderBuf(br, scratch)
			var extraI int64
			if err != nil {
				return err
			}
			switch maj {
			case cbg.MajUnsignedInt:
				extraI = int64(extra)
				if extraI < 0 {
					return fmt.Errorf("int64 positive overflow")
				}
			case cbg.MajNegativeInt:
				extraI = int64(extra)
				if extraI < 0 {
					return fmt.Errorf("int64 negative oveflow")
				}
				extraI = -1 - extraI
			default:
				return fmt.Errorf("wrong type for int64 field: %d", maj)
			}

			t.SignedOthers[i] = int64(extraI)
		}
	}

	// t.Test ([][]uint8) (slice)

	maj, extra, err = cbg.CborReadHeaderBuf(br, scratch)
	if err != nil {
		return err
	}

	if extra > cbg.MaxLength {
		return fmt.Errorf("t.Test: array too large (%d)", extra)
	}

	if maj != cbg.MajArray {
		return fmt.Errorf("expected cbor array")
	}

	if extra > 0 {
		t.Test = make([][]uint8, extra)
	}

	for i := 0; i < int(extra); i++ {
		{
			var maj byte
			var extra uint64
			var err error

			maj, extra, err = cbg.CborReadHeaderBuf(br, scratch)
			if err != nil {
				return err
			}

			if extra > cbg.ByteArrayMaxLen {
				return fmt.Errorf("t.Test[i]: byte array too large (%d)", extra)
			}
			if maj != cbg.MajByteString {
				return fmt.Errorf("expected byte array")
			}

			if extra > 0 {
				t.Test[i] = make([]uint8, extra)
			}

			if _, err := io.ReadFull(br, t.Test[i][:]); err != nil {
				return err
			}
		}
	}

	// t.Dog (string) (string)

	{
		sval, err := cbg.ReadStringBuf(br, scratch)
		if err != nil {
			return err
		}

		t.Dog = string(sval)
	}
	// t.Numbers ([]testing.NamedNumber) (slice)

	maj, extra, err = cbg.CborReadHeaderBuf(br, scratch)
	if err != nil {
		return err
	}

	if extra > cbg.MaxLength {
		return fmt.Errorf("t.Numbers: array too large (%d)", extra)
	}

	if maj != cbg.MajArray {
		return fmt.Errorf("expected cbor array")
	}

	if extra > 0 {
		t.Numbers = make([]NamedNumber, extra)
	}

	for i := 0; i < int(extra); i++ {

		maj, val, err := cbg.CborReadHeaderBuf(br, scratch)
		if err != nil {
			return xerrors.Errorf("failed to read uint64 for t.Numbers slice: %w", err)
		}

		if maj != cbg.MajUnsignedInt {
			return xerrors.Errorf("value read for array t.Numbers was not a uint, instead got %d", maj)
		}

		t.Numbers[i] = NamedNumber(val)
	}

	// t.Pizza (uint64) (uint64)

	{

		b, err := br.ReadByte()
		if err != nil {
			return err
		}
		if b != cbg.CborNull[0] {
			if err := br.UnreadByte(); err != nil {
				return err
			}
			maj, extra, err = cbg.CborReadHeaderBuf(br, scratch)
			if err != nil {
				return err
			}
			if maj != cbg.MajUnsignedInt {
				return fmt.Errorf("wrong type for uint64 field")
			}
			typed := uint64(extra)
			t.Pizza = &typed
		}

	}
	// t.PointyPizza (testing.NamedNumber) (uint64)

	{

		b, err := br.ReadByte()
		if err != nil {
			return err
		}
		if b != cbg.CborNull[0] {
			if err := br.UnreadByte(); err != nil {
				return err
			}
			maj, extra, err = cbg.CborReadHeaderBuf(br, scratch)
			if err != nil {
				return err
			}
			if maj != cbg.MajUnsignedInt {
				return fmt.Errorf("wrong type for uint64 field")
			}
			typed := NamedNumber(extra)
			t.PointyPizza = &typed
		}

	}
	// t.Arrrrrghay ([3]testing.SimpleTypeOne) (array)

	maj, extra, err = cbg.CborReadHeaderBuf(br, scratch)
	if err != nil {
		return err
	}

	if extra > cbg.MaxLength {
		return fmt.Errorf("t.Arrrrrghay: array too large (%d)", extra)
	}

	if maj != cbg.MajArray {
		return fmt.Errorf("expected cbor array")
	}

	if extra != 3 {
		return fmt.Errorf("expected array to have 3 elements")
	}

	t.Arrrrrghay = [3]SimpleTypeOne{}

	for i := 0; i < int(extra); i++ {

		var v SimpleTypeOne
		if err := v.UnmarshalCBOR(br); err != nil {
			return err
		}

		t.Arrrrrghay[i] = v
	}

	return nil
}

var lengthBufDeferredContainer = []byte{131}

func (t *DeferredContainer) MarshalCBOR(w io.Writer) error {
	if t == nil {
		_, err := w.Write(cbg.CborNull)
		return err
	}
	if _, err := w.Write(lengthBufDeferredContainer); err != nil {
		return err
	}

	scratch := make([]byte, 9)

	// t.Stuff (testing.SimpleTypeOne) (struct)
	if err := t.Stuff.MarshalCBOR(w); err != nil {
		return err
	}

	// t.Deferred (typegen.Deferred) (struct)
	if err := t.Deferred.MarshalCBOR(w); err != nil {
		return err
	}

	// t.Value (uint64) (uint64)

	if err := cbg.WriteMajorTypeHeaderBuf(scratch, w, cbg.MajUnsignedInt, uint64(t.Value)); err != nil {
		return err
	}

	return nil
}

func (t *DeferredContainer) UnmarshalCBOR(r io.Reader) (err error) {
	*t = DeferredContainer{}

	br := cbg.GetPeeker(r)
	scratch := make([]byte, 8)

	hasReadOnce := false
	defer func() {
		if err == io.EOF && hasReadOnce {
			err = io.ErrUnexpectedEOF
		}
	}()
	maj, extra, err := cbg.CborReadHeaderBuf(br, scratch)
	if err != nil {
		return err
	}
	hasReadOnce = true

	if maj != cbg.MajArray {
		return fmt.Errorf("cbor input should be of type array")
	}

	if extra != 3 {
		return fmt.Errorf("cbor input had wrong number of fields")
	}

	// t.Stuff (testing.SimpleTypeOne) (struct)

	{

		b, err := br.ReadByte()
		if err != nil {
			return err
		}
		if b != cbg.CborNull[0] {
			if err := br.UnreadByte(); err != nil {
				return err
			}
			t.Stuff = new(SimpleTypeOne)
			if err := t.Stuff.UnmarshalCBOR(br); err != nil {
				return xerrors.Errorf("unmarshaling t.Stuff pointer: %w", err)
			}
		}

	}
	// t.Deferred (typegen.Deferred) (struct)

	{

		t.Deferred = new(cbg.Deferred)

		if err := t.Deferred.UnmarshalCBOR(br); err != nil {
			return xerrors.Errorf("failed to read deferred field: %w", err)
		}
	}
	// t.Value (uint64) (uint64)

	{

		maj, extra, err = cbg.CborReadHeaderBuf(br, scratch)
		if err != nil {
			return err
		}
		if maj != cbg.MajUnsignedInt {
			return fmt.Errorf("wrong type for uint64 field")
		}
		t.Value = uint64(extra)

	}
	return nil
}

var lengthBufFixedArrays = []byte{131}

func (t *FixedArrays) MarshalCBOR(w io.Writer) error {
	if t == nil {
		_, err := w.Write(cbg.CborNull)
		return err
	}
	if _, err := w.Write(lengthBufFixedArrays); err != nil {
		return err
	}

	scratch := make([]byte, 9)

	// t.Bytes ([20]uint8) (array)
	if len(t.Bytes) > cbg.ByteArrayMaxLen {
		return xerrors.Errorf("Byte array in field t.Bytes was too long")
	}

	if err := cbg.WriteMajorTypeHeaderBuf(scratch, w, cbg.MajByteString, uint64(len(t.Bytes))); err != nil {
		return err
	}

	if _, err := w.Write(t.Bytes[:]); err != nil {
		return err
	}

	// t.Uint8 ([20]uint8) (array)
	if len(t.Uint8) > cbg.ByteArrayMaxLen {
		return xerrors.Errorf("Byte array in field t.Uint8 was too long")
	}

	if err := cbg.WriteMajorTypeHeaderBuf(scratch, w, cbg.MajByteString, uint64(len(t.Uint8))); err != nil {
		return err
	}

	if _, err := w.Write(t.Uint8[:]); err != nil {
		return err
	}

	// t.Uint64 ([20]uint64) (array)
	if len(t.Uint64) > cbg.MaxLength {
		return xerrors.Errorf("Slice value in field t.Uint64 was too long")
	}

	if err := cbg.WriteMajorTypeHeaderBuf(scratch, w, cbg.MajArray, uint64(len(t.Uint64))); err != nil {
		return err
	}
	for _, v := range t.Uint64 {
		if err := cbg.CborWriteHeader(w, cbg.MajUnsignedInt, uint64(v)); err != nil {
			return err
		}
	}
	return nil
}

func (t *FixedArrays) UnmarshalCBOR(r io.Reader) (err error) {
	*t = FixedArrays{}

	br := cbg.GetPeeker(r)
	scratch := make([]byte, 8)

	hasReadOnce := false
	defer func() {
		if err == io.EOF && hasReadOnce {
			err = io.ErrUnexpectedEOF
		}
	}()
	maj, extra, err := cbg.CborReadHeaderBuf(br, scratch)
	if err != nil {
		return err
	}
	hasReadOnce = true

	if maj != cbg.MajArray {
		return fmt.Errorf("cbor input should be of type array")
	}

	if extra != 3 {
		return fmt.Errorf("cbor input had wrong number of fields")
	}

	// t.Bytes ([20]uint8) (array)

	maj, extra, err = cbg.CborReadHeaderBuf(br, scratch)
	if err != nil {
		return err
	}

	if extra > cbg.ByteArrayMaxLen {
		return fmt.Errorf("t.Bytes: byte array too large (%d)", extra)
	}
	if maj != cbg.MajByteString {
		return fmt.Errorf("expected byte array")
	}

	if extra != 20 {
		return fmt.Errorf("expected array to have 20 elements")
	}

	t.Bytes = [20]uint8{}

	if _, err := io.ReadFull(br, t.Bytes[:]); err != nil {
		return err
	}
	// t.Uint8 ([20]uint8) (array)

	maj, extra, err = cbg.CborReadHeaderBuf(br, scratch)
	if err != nil {
		return err
	}

	if extra > cbg.ByteArrayMaxLen {
		return fmt.Errorf("t.Uint8: byte array too large (%d)", extra)
	}
	if maj != cbg.MajByteString {
		return fmt.Errorf("expected byte array")
	}

	if extra != 20 {
		return fmt.Errorf("expected array to have 20 elements")
	}

	t.Uint8 = [20]uint8{}

	if _, err := io.ReadFull(br, t.Uint8[:]); err != nil {
		return err
	}
	// t.Uint64 ([20]uint64) (array)

	maj, extra, err = cbg.CborReadHeaderBuf(br, scratch)
	if err != nil {
		return err
	}

	if extra > cbg.MaxLength {
		return fmt.Errorf("t.Uint64: array too large (%d)", extra)
	}

	if maj != cbg.MajArray {
		return fmt.Errorf("expected cbor array")
	}

	if extra != 20 {
		return fmt.Errorf("expected array to have 20 elements")
	}

	t.Uint64 = [20]uint64{}

	for i := 0; i < int(extra); i++ {

		maj, val, err := cbg.CborReadHeaderBuf(br, scratch)
		if err != nil {
			return xerrors.Errorf("failed to read uint64 for t.Uint64 slice: %w", err)
		}

		if maj != cbg.MajUnsignedInt {
			return xerrors.Errorf("value read for array t.Uint64 was not a uint, instead got %d", maj)
		}

		t.Uint64[i] = uint64(val)
	}

	return nil
}

var lengthBufThingWithSomeTime = []byte{131}

func (t *ThingWithSomeTime) MarshalCBOR(w io.Writer) error {
	if t == nil {
		_, err := w.Write(cbg.CborNull)
		return err
	}
	if _, err := w.Write(lengthBufThingWithSomeTime); err != nil {
		return err
	}

	scratch := make([]byte, 9)

	// t.When (typegen.CborTime) (struct)
	if err := t.When.MarshalCBOR(w); err != nil {
		return err
	}

	// t.Stuff (int64) (int64)
	if t.Stuff >= 0 {
		if err := cbg.WriteMajorTypeHeaderBuf(scratch, w, cbg.MajUnsignedInt, uint64(t.Stuff)); err != nil {
			return err
		}
	} else {
		if err := cbg.WriteMajorTypeHeaderBuf(scratch, w, cbg.MajNegativeInt, uint64(-t.Stuff-1)); err != nil {
			return err
		}
	}

	// t.CatName (string) (string)
	if len(t.CatName) > cbg.MaxLength {
		return xerrors.Errorf("Value in field t.CatName was too long")
	}

	if err := cbg.WriteMajorTypeHeaderBuf(scratch, w, cbg.MajTextString, uint64(len(t.CatName))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string(t.CatName)); err != nil {
		return err
	}
	return nil
}

func (t *ThingWithSomeTime) UnmarshalCBOR(r io.Reader) (err error) {
	*t = ThingWithSomeTime{}

	br := cbg.GetPeeker(r)
	scratch := make([]byte, 8)

	hasReadOnce := false
	defer func() {
		if err == io.EOF && hasReadOnce {
			err = io.ErrUnexpectedEOF
		}
	}()
	maj, extra, err := cbg.CborReadHeaderBuf(br, scratch)
	if err != nil {
		return err
	}
	hasReadOnce = true

	if maj != cbg.MajArray {
		return fmt.Errorf("cbor input should be of type array")
	}

	if extra != 3 {
		return fmt.Errorf("cbor input had wrong number of fields")
	}

	// t.When (typegen.CborTime) (struct)

	{

		if err := t.When.UnmarshalCBOR(br); err != nil {
			return xerrors.Errorf("unmarshaling t.When: %w", err)
		}

	}
	// t.Stuff (int64) (int64)
	{
		maj, extra, err := cbg.CborReadHeaderBuf(br, scratch)
		var extraI int64
		if err != nil {
			return err
		}
		switch maj {
		case cbg.MajUnsignedInt:
			extraI = int64(extra)
			if extraI < 0 {
				return fmt.Errorf("int64 positive overflow")
			}
		case cbg.MajNegativeInt:
			extraI = int64(extra)
			if extraI < 0 {
				return fmt.Errorf("int64 negative oveflow")
			}
			extraI = -1 - extraI
		default:
			return fmt.Errorf("wrong type for int64 field: %d", maj)
		}

		t.Stuff = int64(extraI)
	}
	// t.CatName (string) (string)

	{
		sval, err := cbg.ReadStringBuf(br, scratch)
		if err != nil {
			return err
		}

		t.CatName = string(sval)
	}
	return nil
}
