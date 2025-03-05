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

var lengthBufLimitedStruct = []byte{131}

func (t *LimitedStruct) MarshalCBOR(w io.Writer) error {
	if t == nil {
		_, err := w.Write(cbg.CborNull)
		return err
	}

	cw := cbg.NewCborWriter(w)

	if _, err := cw.Write(lengthBufLimitedStruct); err != nil {
		return err
	}

	// t.Arr ([]uint64) (slice)
	if len(t.Arr) > 10 {
		return xerrors.Errorf("Slice value in field t.Arr was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajArray, uint64(len(t.Arr))); err != nil {
		return err
	}
	for _, v := range t.Arr {

		if err := cw.WriteMajorTypeHeader(cbg.MajUnsignedInt, uint64(v)); err != nil {
			return err
		}

	}

	// t.Byts ([]uint8) (slice)
	if len(t.Byts) > 9 {
		return xerrors.Errorf("Byte array in field t.Byts was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajByteString, uint64(len(t.Byts))); err != nil {
		return err
	}

	if _, err := cw.Write(t.Byts); err != nil {
		return err
	}

	// t.Str (string) (string)
	if len(t.Str) > 8 {
		return xerrors.Errorf("Value in field t.Str was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajTextString, uint64(len(t.Str))); err != nil {
		return err
	}
	if _, err := cw.WriteString(string(t.Str)); err != nil {
		return err
	}
	return nil
}

func (t *LimitedStruct) UnmarshalCBOR(r io.Reader) (err error) {
	*t = LimitedStruct{}

	cr := cbg.NewCborReader(r)

	maj, extra, err := cr.ReadHeader()
	if err != nil {
		return err
	}
	defer func() {
		if err == io.EOF {
			err = io.ErrUnexpectedEOF
		}
	}()

	if maj != cbg.MajArray {
		return fmt.Errorf("cbor input should be of type array")
	}

	if extra != 3 {
		return fmt.Errorf("cbor input had wrong number of fields")
	}

	// t.Arr ([]uint64) (slice)

	{

		maj, extra, err := cr.ReadHeader()
		if err != nil {
			return err
		}

		if extra > 10 {
			return fmt.Errorf("t.Arr: array too large (%d)", extra)
		}

		if maj != cbg.MajArray {
			return fmt.Errorf("expected cbor array")
		}

		if extra > 0 {
			t.Arr = make([]uint64, extra)
		}

		for i := 0; i < int(extra); i++ {
			{
				var maj byte
				var extra uint64
				var err error
				_ = maj
				_ = extra
				_ = err

				{

					maj, extra, err := cr.ReadHeader()
					if err != nil {
						return err
					}
					if maj != cbg.MajUnsignedInt {
						return fmt.Errorf("wrong type for uint64 field")
					}
					t.Arr[i] = uint64(extra)

				}

			}
		}
	}
	// t.Byts ([]uint8) (slice)

	{

		maj, extra, err := cr.ReadHeader()
		if err != nil {
			return err
		}

		if extra > 9 {
			return fmt.Errorf("t.Byts: byte array too large (%d)", extra)
		}
		if maj != cbg.MajByteString {
			return fmt.Errorf("expected byte array")
		}

		if extra > 0 {
			t.Byts = make([]uint8, extra)
		}

		if _, err := io.ReadFull(cr, t.Byts); err != nil {
			return err
		}

	}
	// t.Str (string) (string)

	{
		sval, err := cbg.ReadStringWithMax(cr, 8)
		if err != nil {
			return err
		}

		t.Str = string(sval)
	}
	return nil
}
