// Code generated by github.com/whyrusleeping/cbor-gen. DO NOT EDIT.

package testing

import (
	"fmt"
	"io"
	"math"
	"sort"

	cid "github.com/ipfs/go-cid"
	cbg "github.com/whyrusleeping/cbor-gen"
)

var _ = fmt.Errorf
var _ = cid.Undef
var _ = math.E
var _ = sort.Sort

func (t *SimpleTypeTree) MarshalCBOR(w io.Writer) error {
	if t == nil {
		_, err := w.Write(cbg.CborNull)
		return err
	}

	cw := cbg.NewCborWriter(w)

	if _, err := cw.Write([]byte{167}); err != nil {
		return err
	}

	// t.Stuff (testing.SimpleTypeTree) (struct)
	if len("Stuff") > cbg.MaxLength {
		return fmt.Errorf("Value in field \"Stuff\" was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajTextString, uint64(len("Stuff"))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string("Stuff")); err != nil {
		return err
	}

	if err := t.Stuff.MarshalCBOR(cw); err != nil {
		return err
	}

	// t.Stufff (testing.SimpleTypeTwo) (struct)
	if len("Stufff") > cbg.MaxLength {
		return fmt.Errorf("Value in field \"Stufff\" was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajTextString, uint64(len("Stufff"))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string("Stufff")); err != nil {
		return err
	}

	if err := t.Stufff.MarshalCBOR(cw); err != nil {
		return err
	}

	// t.Others ([]uint64) (slice)
	if len("Others") > cbg.MaxLength {
		return fmt.Errorf("Value in field \"Others\" was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajTextString, uint64(len("Others"))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string("Others")); err != nil {
		return err
	}

	if len(t.Others) > cbg.MaxLength {
		return fmt.Errorf("Slice value in field t.Others was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajArray, uint64(len(t.Others))); err != nil {
		return err
	}
	for _, v := range t.Others {
		if err := cw.CborWriteHeader(cbg.MajUnsignedInt, uint64(v)); err != nil {
			return err
		}
	}

	// t.Test ([][]uint8) (slice)
	if len("Test") > cbg.MaxLength {
		return fmt.Errorf("Value in field \"Test\" was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajTextString, uint64(len("Test"))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string("Test")); err != nil {
		return err
	}

	if len(t.Test) > cbg.MaxLength {
		return fmt.Errorf("Slice value in field t.Test was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajArray, uint64(len(t.Test))); err != nil {
		return err
	}
	for _, v := range t.Test {
		if len(v) > cbg.ByteArrayMaxLen {
			return fmt.Errorf("Byte array in field v was too long")
		}

		if err := cw.WriteMajorTypeHeader(cbg.MajByteString, uint64(len(v))); err != nil {
			return err
		}

		if _, err := cw.Write(v[:]); err != nil {
			return err
		}
	}

	// t.Dog (string) (string)
	if len("Dog") > cbg.MaxLength {
		return fmt.Errorf("Value in field \"Dog\" was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajTextString, uint64(len("Dog"))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string("Dog")); err != nil {
		return err
	}

	if len(t.Dog) > cbg.MaxLength {
		return fmt.Errorf("Value in field t.Dog was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajTextString, uint64(len(t.Dog))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string(t.Dog)); err != nil {
		return err
	}

	// t.SixtyThreeBitIntegerWithASignBit (int64) (int64)
	if len("SixtyThreeBitIntegerWithASignBit") > cbg.MaxLength {
		return fmt.Errorf("Value in field \"SixtyThreeBitIntegerWithASignBit\" was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajTextString, uint64(len("SixtyThreeBitIntegerWithASignBit"))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string("SixtyThreeBitIntegerWithASignBit")); err != nil {
		return err
	}

	if t.SixtyThreeBitIntegerWithASignBit >= 0 {
		if err := cw.WriteMajorTypeHeader(cbg.MajUnsignedInt, uint64(t.SixtyThreeBitIntegerWithASignBit)); err != nil {
			return err
		}
	} else {
		if err := cw.WriteMajorTypeHeader(cbg.MajNegativeInt, uint64(-t.SixtyThreeBitIntegerWithASignBit-1)); err != nil {
			return err
		}
	}

	// t.NotPizza (uint64) (uint64)
	if len("NotPizza") > cbg.MaxLength {
		return fmt.Errorf("Value in field \"NotPizza\" was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajTextString, uint64(len("NotPizza"))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string("NotPizza")); err != nil {
		return err
	}

	if t.NotPizza == nil {
		if _, err := cw.Write(cbg.CborNull); err != nil {
			return err
		}
	} else {
		if err := cw.WriteMajorTypeHeader(cbg.MajUnsignedInt, uint64(*t.NotPizza)); err != nil {
			return err
		}
	}

	return nil
}

func (t *SimpleTypeTree) UnmarshalCBOR(r io.Reader) (err error) {
	*t = SimpleTypeTree{}

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

	if maj != cbg.MajMap {
		return fmt.Errorf("cbor input should be of type map")
	}

	if extra > cbg.MaxLength {
		return fmt.Errorf("SimpleTypeTree: map struct too large (%d)", extra)
	}

	var name string
	n := extra

	for i := uint64(0); i < n; i++ {

		{
			sval, err := cbg.ReadString(cr)
			if err != nil {
				return err
			}

			name = string(sval)
		}

		switch name {
		// t.Stuff (testing.SimpleTypeTree) (struct)
		case "Stuff":

			{

				b, err := cr.ReadByte()
				if err != nil {
					return err
				}
				if b != cbg.CborNull[0] {
					if err := cr.UnreadByte(); err != nil {
						return err
					}
					t.Stuff = new(SimpleTypeTree)
					if err := t.Stuff.UnmarshalCBOR(cr); err != nil {
						return fmt.Errorf("unmarshaling t.Stuff pointer: %w", err)
					}
				}

			}
			// t.Stufff (testing.SimpleTypeTwo) (struct)
		case "Stufff":

			{

				b, err := cr.ReadByte()
				if err != nil {
					return err
				}
				if b != cbg.CborNull[0] {
					if err := cr.UnreadByte(); err != nil {
						return err
					}
					t.Stufff = new(SimpleTypeTwo)
					if err := t.Stufff.UnmarshalCBOR(cr); err != nil {
						return fmt.Errorf("unmarshaling t.Stufff pointer: %w", err)
					}
				}

			}
			// t.Others ([]uint64) (slice)
		case "Others":

			maj, extra, err = cr.ReadHeader()
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

				maj, val, err := cr.ReadHeader()
				if err != nil {
					return fmt.Errorf("failed to read uint64 for t.Others slice: %w", err)
				}

				if maj != cbg.MajUnsignedInt {
					return fmt.Errorf("value read for array t.Others was not a uint, instead got %d", maj)
				}

				t.Others[i] = uint64(val)
			}

			// t.Test ([][]uint8) (slice)
		case "Test":

			maj, extra, err = cr.ReadHeader()
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

					maj, extra, err = cr.ReadHeader()
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

					if _, err := io.ReadFull(cr, t.Test[i][:]); err != nil {
						return err
					}
				}
			}

			// t.Dog (string) (string)
		case "Dog":

			{
				sval, err := cbg.ReadString(cr)
				if err != nil {
					return err
				}

				t.Dog = string(sval)
			}
			// t.SixtyThreeBitIntegerWithASignBit (int64) (int64)
		case "SixtyThreeBitIntegerWithASignBit":
			{
				maj, extra, err := cr.ReadHeader()
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

				t.SixtyThreeBitIntegerWithASignBit = int64(extraI)
			}
			// t.NotPizza (uint64) (uint64)
		case "NotPizza":

			{

				b, err := cr.ReadByte()
				if err != nil {
					return err
				}
				if b != cbg.CborNull[0] {
					if err := cr.UnreadByte(); err != nil {
						return err
					}
					maj, extra, err = cr.ReadHeader()
					if err != nil {
						return err
					}
					if maj != cbg.MajUnsignedInt {
						return fmt.Errorf("wrong type for uint64 field")
					}
					typed := uint64(extra)
					t.NotPizza = &typed
				}

			}

		default:
			// Field doesn't exist on this type, so ignore it
			cbg.ScanForLinks(r, func(cid.Cid) {})
		}
	}

	return nil
}
func (t *NeedScratchForMap) MarshalCBOR(w io.Writer) error {
	if t == nil {
		_, err := w.Write(cbg.CborNull)
		return err
	}

	cw := cbg.NewCborWriter(w)

	if _, err := cw.Write([]byte{161}); err != nil {
		return err
	}

	// t.Thing (bool) (bool)
	if len("Thing") > cbg.MaxLength {
		return fmt.Errorf("Value in field \"Thing\" was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajTextString, uint64(len("Thing"))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string("Thing")); err != nil {
		return err
	}

	if err := cbg.WriteBool(w, t.Thing); err != nil {
		return err
	}
	return nil
}

func (t *NeedScratchForMap) UnmarshalCBOR(r io.Reader) (err error) {
	*t = NeedScratchForMap{}

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

	if maj != cbg.MajMap {
		return fmt.Errorf("cbor input should be of type map")
	}

	if extra > cbg.MaxLength {
		return fmt.Errorf("NeedScratchForMap: map struct too large (%d)", extra)
	}

	var name string
	n := extra

	for i := uint64(0); i < n; i++ {

		{
			sval, err := cbg.ReadString(cr)
			if err != nil {
				return err
			}

			name = string(sval)
		}

		switch name {
		// t.Thing (bool) (bool)
		case "Thing":

			maj, extra, err = cr.ReadHeader()
			if err != nil {
				return err
			}
			if maj != cbg.MajOther {
				return fmt.Errorf("booleans must be major type 7")
			}
			switch extra {
			case 20:
				t.Thing = false
			case 21:
				t.Thing = true
			default:
				return fmt.Errorf("booleans are either major type 7, value 20 or 21 (got %d)", extra)
			}

		default:
			// Field doesn't exist on this type, so ignore it
			cbg.ScanForLinks(r, func(cid.Cid) {})
		}
	}

	return nil
}
func (t *SimpleStructV1) MarshalCBOR(w io.Writer) error {
	if t == nil {
		_, err := w.Write(cbg.CborNull)
		return err
	}

	cw := cbg.NewCborWriter(w)

	if _, err := cw.Write([]byte{167}); err != nil {
		return err
	}

	// t.OldStr (string) (string)
	if len("OldStr") > cbg.MaxLength {
		return fmt.Errorf("Value in field \"OldStr\" was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajTextString, uint64(len("OldStr"))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string("OldStr")); err != nil {
		return err
	}

	if len(t.OldStr) > cbg.MaxLength {
		return fmt.Errorf("Value in field t.OldStr was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajTextString, uint64(len(t.OldStr))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string(t.OldStr)); err != nil {
		return err
	}

	// t.OldBytes ([]uint8) (slice)
	if len("OldBytes") > cbg.MaxLength {
		return fmt.Errorf("Value in field \"OldBytes\" was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajTextString, uint64(len("OldBytes"))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string("OldBytes")); err != nil {
		return err
	}

	if len(t.OldBytes) > cbg.ByteArrayMaxLen {
		return fmt.Errorf("Byte array in field t.OldBytes was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajByteString, uint64(len(t.OldBytes))); err != nil {
		return err
	}

	if _, err := cw.Write(t.OldBytes[:]); err != nil {
		return err
	}

	// t.OldNum (uint64) (uint64)
	if len("OldNum") > cbg.MaxLength {
		return fmt.Errorf("Value in field \"OldNum\" was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajTextString, uint64(len("OldNum"))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string("OldNum")); err != nil {
		return err
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajUnsignedInt, uint64(t.OldNum)); err != nil {
		return err
	}

	// t.OldPtr (cid.Cid) (struct)
	if len("OldPtr") > cbg.MaxLength {
		return fmt.Errorf("Value in field \"OldPtr\" was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajTextString, uint64(len("OldPtr"))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string("OldPtr")); err != nil {
		return err
	}

	if t.OldPtr == nil {
		if _, err := cw.Write(cbg.CborNull); err != nil {
			return err
		}
	} else {
		if err := cbg.WriteCid(cw, *t.OldPtr); err != nil {
			return fmt.Errorf("failed to write cid field t.OldPtr: %w", err)
		}
	}

	// t.OldMap (map[string]testing.SimpleTypeOne) (map)
	if len("OldMap") > cbg.MaxLength {
		return fmt.Errorf("Value in field \"OldMap\" was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajTextString, uint64(len("OldMap"))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string("OldMap")); err != nil {
		return err
	}

	{
		if len(t.OldMap) > 4096 {
			return fmt.Errorf("cannot marshal t.OldMap map too large")
		}

		if err := cw.WriteMajorTypeHeader(cbg.MajMap, uint64(len(t.OldMap))); err != nil {
			return err
		}

		keys := make([]string, 0, len(t.OldMap))
		for k := range t.OldMap {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			v := t.OldMap[k]

			if len(k) > cbg.MaxLength {
				return fmt.Errorf("Value in field k was too long")
			}

			if err := cw.WriteMajorTypeHeader(cbg.MajTextString, uint64(len(k))); err != nil {
				return err
			}
			if _, err := io.WriteString(w, string(k)); err != nil {
				return err
			}

			if err := v.MarshalCBOR(cw); err != nil {
				return err
			}

		}
	}

	// t.OldArray ([]testing.SimpleTypeOne) (slice)
	if len("OldArray") > cbg.MaxLength {
		return fmt.Errorf("Value in field \"OldArray\" was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajTextString, uint64(len("OldArray"))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string("OldArray")); err != nil {
		return err
	}

	if len(t.OldArray) > cbg.MaxLength {
		return fmt.Errorf("Slice value in field t.OldArray was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajArray, uint64(len(t.OldArray))); err != nil {
		return err
	}
	for _, v := range t.OldArray {
		if err := v.MarshalCBOR(cw); err != nil {
			return err
		}
	}

	// t.OldStruct (testing.SimpleTypeOne) (struct)
	if len("OldStruct") > cbg.MaxLength {
		return fmt.Errorf("Value in field \"OldStruct\" was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajTextString, uint64(len("OldStruct"))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string("OldStruct")); err != nil {
		return err
	}

	if err := t.OldStruct.MarshalCBOR(cw); err != nil {
		return err
	}
	return nil
}

func (t *SimpleStructV1) UnmarshalCBOR(r io.Reader) (err error) {
	*t = SimpleStructV1{}

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

	if maj != cbg.MajMap {
		return fmt.Errorf("cbor input should be of type map")
	}

	if extra > cbg.MaxLength {
		return fmt.Errorf("SimpleStructV1: map struct too large (%d)", extra)
	}

	var name string
	n := extra

	for i := uint64(0); i < n; i++ {

		{
			sval, err := cbg.ReadString(cr)
			if err != nil {
				return err
			}

			name = string(sval)
		}

		switch name {
		// t.OldStr (string) (string)
		case "OldStr":

			{
				sval, err := cbg.ReadString(cr)
				if err != nil {
					return err
				}

				t.OldStr = string(sval)
			}
			// t.OldBytes ([]uint8) (slice)
		case "OldBytes":

			maj, extra, err = cr.ReadHeader()
			if err != nil {
				return err
			}

			if extra > cbg.ByteArrayMaxLen {
				return fmt.Errorf("t.OldBytes: byte array too large (%d)", extra)
			}
			if maj != cbg.MajByteString {
				return fmt.Errorf("expected byte array")
			}

			if extra > 0 {
				t.OldBytes = make([]uint8, extra)
			}

			if _, err := io.ReadFull(cr, t.OldBytes[:]); err != nil {
				return err
			}
			// t.OldNum (uint64) (uint64)
		case "OldNum":

			{

				maj, extra, err = cr.ReadHeader()
				if err != nil {
					return err
				}
				if maj != cbg.MajUnsignedInt {
					return fmt.Errorf("wrong type for uint64 field")
				}
				t.OldNum = uint64(extra)

			}
			// t.OldPtr (cid.Cid) (struct)
		case "OldPtr":

			{

				b, err := cr.ReadByte()
				if err != nil {
					return err
				}
				if b != cbg.CborNull[0] {
					if err := cr.UnreadByte(); err != nil {
						return err
					}

					c, err := cbg.ReadCid(cr)
					if err != nil {
						return fmt.Errorf("failed to read cid field t.OldPtr: %w", err)
					}

					t.OldPtr = &c
				}

			}
			// t.OldMap (map[string]testing.SimpleTypeOne) (map)
		case "OldMap":

			maj, extra, err = cr.ReadHeader()
			if err != nil {
				return err
			}
			if maj != cbg.MajMap {
				return fmt.Errorf("expected a map (major type 5)")
			}
			if extra > 4096 {
				return fmt.Errorf("t.OldMap: map too large")
			}

			t.OldMap = make(map[string]SimpleTypeOne, extra)

			for i, l := 0, int(extra); i < l; i++ {

				var k string

				{
					sval, err := cbg.ReadString(cr)
					if err != nil {
						return err
					}

					k = string(sval)
				}

				var v SimpleTypeOne

				{

					if err := v.UnmarshalCBOR(cr); err != nil {
						return fmt.Errorf("unmarshaling v: %w", err)
					}

				}

				t.OldMap[k] = v

			}
			// t.OldArray ([]testing.SimpleTypeOne) (slice)
		case "OldArray":

			maj, extra, err = cr.ReadHeader()
			if err != nil {
				return err
			}

			if extra > cbg.MaxLength {
				return fmt.Errorf("t.OldArray: array too large (%d)", extra)
			}

			if maj != cbg.MajArray {
				return fmt.Errorf("expected cbor array")
			}

			if extra > 0 {
				t.OldArray = make([]SimpleTypeOne, extra)
			}

			for i := 0; i < int(extra); i++ {

				var v SimpleTypeOne
				if err := v.UnmarshalCBOR(cr); err != nil {
					return err
				}

				t.OldArray[i] = v
			}

			// t.OldStruct (testing.SimpleTypeOne) (struct)
		case "OldStruct":

			{

				if err := t.OldStruct.UnmarshalCBOR(cr); err != nil {
					return fmt.Errorf("unmarshaling t.OldStruct: %w", err)
				}

			}

		default:
			// Field doesn't exist on this type, so ignore it
			cbg.ScanForLinks(r, func(cid.Cid) {})
		}
	}

	return nil
}
func (t *SimpleStructV2) MarshalCBOR(w io.Writer) error {
	if t == nil {
		_, err := w.Write(cbg.CborNull)
		return err
	}

	cw := cbg.NewCborWriter(w)

	if _, err := cw.Write([]byte{174}); err != nil {
		return err
	}

	// t.OldStr (string) (string)
	if len("OldStr") > cbg.MaxLength {
		return fmt.Errorf("Value in field \"OldStr\" was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajTextString, uint64(len("OldStr"))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string("OldStr")); err != nil {
		return err
	}

	if len(t.OldStr) > cbg.MaxLength {
		return fmt.Errorf("Value in field t.OldStr was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajTextString, uint64(len(t.OldStr))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string(t.OldStr)); err != nil {
		return err
	}

	// t.NewStr (string) (string)
	if len("NewStr") > cbg.MaxLength {
		return fmt.Errorf("Value in field \"NewStr\" was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajTextString, uint64(len("NewStr"))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string("NewStr")); err != nil {
		return err
	}

	if len(t.NewStr) > cbg.MaxLength {
		return fmt.Errorf("Value in field t.NewStr was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajTextString, uint64(len(t.NewStr))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string(t.NewStr)); err != nil {
		return err
	}

	// t.OldBytes ([]uint8) (slice)
	if len("OldBytes") > cbg.MaxLength {
		return fmt.Errorf("Value in field \"OldBytes\" was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajTextString, uint64(len("OldBytes"))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string("OldBytes")); err != nil {
		return err
	}

	if len(t.OldBytes) > cbg.ByteArrayMaxLen {
		return fmt.Errorf("Byte array in field t.OldBytes was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajByteString, uint64(len(t.OldBytes))); err != nil {
		return err
	}

	if _, err := cw.Write(t.OldBytes[:]); err != nil {
		return err
	}

	// t.NewBytes ([]uint8) (slice)
	if len("NewBytes") > cbg.MaxLength {
		return fmt.Errorf("Value in field \"NewBytes\" was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajTextString, uint64(len("NewBytes"))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string("NewBytes")); err != nil {
		return err
	}

	if len(t.NewBytes) > cbg.ByteArrayMaxLen {
		return fmt.Errorf("Byte array in field t.NewBytes was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajByteString, uint64(len(t.NewBytes))); err != nil {
		return err
	}

	if _, err := cw.Write(t.NewBytes[:]); err != nil {
		return err
	}

	// t.OldNum (uint64) (uint64)
	if len("OldNum") > cbg.MaxLength {
		return fmt.Errorf("Value in field \"OldNum\" was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajTextString, uint64(len("OldNum"))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string("OldNum")); err != nil {
		return err
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajUnsignedInt, uint64(t.OldNum)); err != nil {
		return err
	}

	// t.NewNum (uint64) (uint64)
	if len("NewNum") > cbg.MaxLength {
		return fmt.Errorf("Value in field \"NewNum\" was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajTextString, uint64(len("NewNum"))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string("NewNum")); err != nil {
		return err
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajUnsignedInt, uint64(t.NewNum)); err != nil {
		return err
	}

	// t.OldPtr (cid.Cid) (struct)
	if len("OldPtr") > cbg.MaxLength {
		return fmt.Errorf("Value in field \"OldPtr\" was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajTextString, uint64(len("OldPtr"))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string("OldPtr")); err != nil {
		return err
	}

	if t.OldPtr == nil {
		if _, err := cw.Write(cbg.CborNull); err != nil {
			return err
		}
	} else {
		if err := cbg.WriteCid(cw, *t.OldPtr); err != nil {
			return fmt.Errorf("failed to write cid field t.OldPtr: %w", err)
		}
	}

	// t.NewPtr (cid.Cid) (struct)
	if len("NewPtr") > cbg.MaxLength {
		return fmt.Errorf("Value in field \"NewPtr\" was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajTextString, uint64(len("NewPtr"))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string("NewPtr")); err != nil {
		return err
	}

	if t.NewPtr == nil {
		if _, err := cw.Write(cbg.CborNull); err != nil {
			return err
		}
	} else {
		if err := cbg.WriteCid(cw, *t.NewPtr); err != nil {
			return fmt.Errorf("failed to write cid field t.NewPtr: %w", err)
		}
	}

	// t.OldMap (map[string]testing.SimpleTypeOne) (map)
	if len("OldMap") > cbg.MaxLength {
		return fmt.Errorf("Value in field \"OldMap\" was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajTextString, uint64(len("OldMap"))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string("OldMap")); err != nil {
		return err
	}

	{
		if len(t.OldMap) > 4096 {
			return fmt.Errorf("cannot marshal t.OldMap map too large")
		}

		if err := cw.WriteMajorTypeHeader(cbg.MajMap, uint64(len(t.OldMap))); err != nil {
			return err
		}

		keys := make([]string, 0, len(t.OldMap))
		for k := range t.OldMap {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			v := t.OldMap[k]

			if len(k) > cbg.MaxLength {
				return fmt.Errorf("Value in field k was too long")
			}

			if err := cw.WriteMajorTypeHeader(cbg.MajTextString, uint64(len(k))); err != nil {
				return err
			}
			if _, err := io.WriteString(w, string(k)); err != nil {
				return err
			}

			if err := v.MarshalCBOR(cw); err != nil {
				return err
			}

		}
	}

	// t.NewMap (map[string]testing.SimpleTypeOne) (map)
	if len("NewMap") > cbg.MaxLength {
		return fmt.Errorf("Value in field \"NewMap\" was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajTextString, uint64(len("NewMap"))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string("NewMap")); err != nil {
		return err
	}

	{
		if len(t.NewMap) > 4096 {
			return fmt.Errorf("cannot marshal t.NewMap map too large")
		}

		if err := cw.WriteMajorTypeHeader(cbg.MajMap, uint64(len(t.NewMap))); err != nil {
			return err
		}

		keys := make([]string, 0, len(t.NewMap))
		for k := range t.NewMap {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			v := t.NewMap[k]

			if len(k) > cbg.MaxLength {
				return fmt.Errorf("Value in field k was too long")
			}

			if err := cw.WriteMajorTypeHeader(cbg.MajTextString, uint64(len(k))); err != nil {
				return err
			}
			if _, err := io.WriteString(w, string(k)); err != nil {
				return err
			}

			if err := v.MarshalCBOR(cw); err != nil {
				return err
			}

		}
	}

	// t.OldArray ([]testing.SimpleTypeOne) (slice)
	if len("OldArray") > cbg.MaxLength {
		return fmt.Errorf("Value in field \"OldArray\" was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajTextString, uint64(len("OldArray"))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string("OldArray")); err != nil {
		return err
	}

	if len(t.OldArray) > cbg.MaxLength {
		return fmt.Errorf("Slice value in field t.OldArray was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajArray, uint64(len(t.OldArray))); err != nil {
		return err
	}
	for _, v := range t.OldArray {
		if err := v.MarshalCBOR(cw); err != nil {
			return err
		}
	}

	// t.NewArray ([]testing.SimpleTypeOne) (slice)
	if len("NewArray") > cbg.MaxLength {
		return fmt.Errorf("Value in field \"NewArray\" was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajTextString, uint64(len("NewArray"))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string("NewArray")); err != nil {
		return err
	}

	if len(t.NewArray) > cbg.MaxLength {
		return fmt.Errorf("Slice value in field t.NewArray was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajArray, uint64(len(t.NewArray))); err != nil {
		return err
	}
	for _, v := range t.NewArray {
		if err := v.MarshalCBOR(cw); err != nil {
			return err
		}
	}

	// t.OldStruct (testing.SimpleTypeOne) (struct)
	if len("OldStruct") > cbg.MaxLength {
		return fmt.Errorf("Value in field \"OldStruct\" was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajTextString, uint64(len("OldStruct"))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string("OldStruct")); err != nil {
		return err
	}

	if err := t.OldStruct.MarshalCBOR(cw); err != nil {
		return err
	}

	// t.NewStruct (testing.SimpleTypeOne) (struct)
	if len("NewStruct") > cbg.MaxLength {
		return fmt.Errorf("Value in field \"NewStruct\" was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajTextString, uint64(len("NewStruct"))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string("NewStruct")); err != nil {
		return err
	}

	if err := t.NewStruct.MarshalCBOR(cw); err != nil {
		return err
	}
	return nil
}

func (t *SimpleStructV2) UnmarshalCBOR(r io.Reader) (err error) {
	*t = SimpleStructV2{}

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

	if maj != cbg.MajMap {
		return fmt.Errorf("cbor input should be of type map")
	}

	if extra > cbg.MaxLength {
		return fmt.Errorf("SimpleStructV2: map struct too large (%d)", extra)
	}

	var name string
	n := extra

	for i := uint64(0); i < n; i++ {

		{
			sval, err := cbg.ReadString(cr)
			if err != nil {
				return err
			}

			name = string(sval)
		}

		switch name {
		// t.OldStr (string) (string)
		case "OldStr":

			{
				sval, err := cbg.ReadString(cr)
				if err != nil {
					return err
				}

				t.OldStr = string(sval)
			}
			// t.NewStr (string) (string)
		case "NewStr":

			{
				sval, err := cbg.ReadString(cr)
				if err != nil {
					return err
				}

				t.NewStr = string(sval)
			}
			// t.OldBytes ([]uint8) (slice)
		case "OldBytes":

			maj, extra, err = cr.ReadHeader()
			if err != nil {
				return err
			}

			if extra > cbg.ByteArrayMaxLen {
				return fmt.Errorf("t.OldBytes: byte array too large (%d)", extra)
			}
			if maj != cbg.MajByteString {
				return fmt.Errorf("expected byte array")
			}

			if extra > 0 {
				t.OldBytes = make([]uint8, extra)
			}

			if _, err := io.ReadFull(cr, t.OldBytes[:]); err != nil {
				return err
			}
			// t.NewBytes ([]uint8) (slice)
		case "NewBytes":

			maj, extra, err = cr.ReadHeader()
			if err != nil {
				return err
			}

			if extra > cbg.ByteArrayMaxLen {
				return fmt.Errorf("t.NewBytes: byte array too large (%d)", extra)
			}
			if maj != cbg.MajByteString {
				return fmt.Errorf("expected byte array")
			}

			if extra > 0 {
				t.NewBytes = make([]uint8, extra)
			}

			if _, err := io.ReadFull(cr, t.NewBytes[:]); err != nil {
				return err
			}
			// t.OldNum (uint64) (uint64)
		case "OldNum":

			{

				maj, extra, err = cr.ReadHeader()
				if err != nil {
					return err
				}
				if maj != cbg.MajUnsignedInt {
					return fmt.Errorf("wrong type for uint64 field")
				}
				t.OldNum = uint64(extra)

			}
			// t.NewNum (uint64) (uint64)
		case "NewNum":

			{

				maj, extra, err = cr.ReadHeader()
				if err != nil {
					return err
				}
				if maj != cbg.MajUnsignedInt {
					return fmt.Errorf("wrong type for uint64 field")
				}
				t.NewNum = uint64(extra)

			}
			// t.OldPtr (cid.Cid) (struct)
		case "OldPtr":

			{

				b, err := cr.ReadByte()
				if err != nil {
					return err
				}
				if b != cbg.CborNull[0] {
					if err := cr.UnreadByte(); err != nil {
						return err
					}

					c, err := cbg.ReadCid(cr)
					if err != nil {
						return fmt.Errorf("failed to read cid field t.OldPtr: %w", err)
					}

					t.OldPtr = &c
				}

			}
			// t.NewPtr (cid.Cid) (struct)
		case "NewPtr":

			{

				b, err := cr.ReadByte()
				if err != nil {
					return err
				}
				if b != cbg.CborNull[0] {
					if err := cr.UnreadByte(); err != nil {
						return err
					}

					c, err := cbg.ReadCid(cr)
					if err != nil {
						return fmt.Errorf("failed to read cid field t.NewPtr: %w", err)
					}

					t.NewPtr = &c
				}

			}
			// t.OldMap (map[string]testing.SimpleTypeOne) (map)
		case "OldMap":

			maj, extra, err = cr.ReadHeader()
			if err != nil {
				return err
			}
			if maj != cbg.MajMap {
				return fmt.Errorf("expected a map (major type 5)")
			}
			if extra > 4096 {
				return fmt.Errorf("t.OldMap: map too large")
			}

			t.OldMap = make(map[string]SimpleTypeOne, extra)

			for i, l := 0, int(extra); i < l; i++ {

				var k string

				{
					sval, err := cbg.ReadString(cr)
					if err != nil {
						return err
					}

					k = string(sval)
				}

				var v SimpleTypeOne

				{

					if err := v.UnmarshalCBOR(cr); err != nil {
						return fmt.Errorf("unmarshaling v: %w", err)
					}

				}

				t.OldMap[k] = v

			}
			// t.NewMap (map[string]testing.SimpleTypeOne) (map)
		case "NewMap":

			maj, extra, err = cr.ReadHeader()
			if err != nil {
				return err
			}
			if maj != cbg.MajMap {
				return fmt.Errorf("expected a map (major type 5)")
			}
			if extra > 4096 {
				return fmt.Errorf("t.NewMap: map too large")
			}

			t.NewMap = make(map[string]SimpleTypeOne, extra)

			for i, l := 0, int(extra); i < l; i++ {

				var k string

				{
					sval, err := cbg.ReadString(cr)
					if err != nil {
						return err
					}

					k = string(sval)
				}

				var v SimpleTypeOne

				{

					if err := v.UnmarshalCBOR(cr); err != nil {
						return fmt.Errorf("unmarshaling v: %w", err)
					}

				}

				t.NewMap[k] = v

			}
			// t.OldArray ([]testing.SimpleTypeOne) (slice)
		case "OldArray":

			maj, extra, err = cr.ReadHeader()
			if err != nil {
				return err
			}

			if extra > cbg.MaxLength {
				return fmt.Errorf("t.OldArray: array too large (%d)", extra)
			}

			if maj != cbg.MajArray {
				return fmt.Errorf("expected cbor array")
			}

			if extra > 0 {
				t.OldArray = make([]SimpleTypeOne, extra)
			}

			for i := 0; i < int(extra); i++ {

				var v SimpleTypeOne
				if err := v.UnmarshalCBOR(cr); err != nil {
					return err
				}

				t.OldArray[i] = v
			}

			// t.NewArray ([]testing.SimpleTypeOne) (slice)
		case "NewArray":

			maj, extra, err = cr.ReadHeader()
			if err != nil {
				return err
			}

			if extra > cbg.MaxLength {
				return fmt.Errorf("t.NewArray: array too large (%d)", extra)
			}

			if maj != cbg.MajArray {
				return fmt.Errorf("expected cbor array")
			}

			if extra > 0 {
				t.NewArray = make([]SimpleTypeOne, extra)
			}

			for i := 0; i < int(extra); i++ {

				var v SimpleTypeOne
				if err := v.UnmarshalCBOR(cr); err != nil {
					return err
				}

				t.NewArray[i] = v
			}

			// t.OldStruct (testing.SimpleTypeOne) (struct)
		case "OldStruct":

			{

				if err := t.OldStruct.UnmarshalCBOR(cr); err != nil {
					return fmt.Errorf("unmarshaling t.OldStruct: %w", err)
				}

			}
			// t.NewStruct (testing.SimpleTypeOne) (struct)
		case "NewStruct":

			{

				if err := t.NewStruct.UnmarshalCBOR(cr); err != nil {
					return fmt.Errorf("unmarshaling t.NewStruct: %w", err)
				}

			}

		default:
			// Field doesn't exist on this type, so ignore it
			cbg.ScanForLinks(r, func(cid.Cid) {})
		}
	}

	return nil
}
func (t *RenamedFields) MarshalCBOR(w io.Writer) error {
	if t == nil {
		_, err := w.Write(cbg.CborNull)
		return err
	}

	cw := cbg.NewCborWriter(w)

	if _, err := cw.Write([]byte{162}); err != nil {
		return err
	}

	// t.Foo (int64) (int64)
	if len("foo") > cbg.MaxLength {
		return fmt.Errorf("Value in field \"foo\" was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajTextString, uint64(len("foo"))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string("foo")); err != nil {
		return err
	}

	if t.Foo >= 0 {
		if err := cw.WriteMajorTypeHeader(cbg.MajUnsignedInt, uint64(t.Foo)); err != nil {
			return err
		}
	} else {
		if err := cw.WriteMajorTypeHeader(cbg.MajNegativeInt, uint64(-t.Foo-1)); err != nil {
			return err
		}
	}

	// t.Bar (string) (string)
	if len("beep") > cbg.MaxLength {
		return fmt.Errorf("Value in field \"beep\" was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajTextString, uint64(len("beep"))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string("beep")); err != nil {
		return err
	}

	if len(t.Bar) > cbg.MaxLength {
		return fmt.Errorf("Value in field t.Bar was too long")
	}

	if err := cw.WriteMajorTypeHeader(cbg.MajTextString, uint64(len(t.Bar))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string(t.Bar)); err != nil {
		return err
	}
	return nil
}

func (t *RenamedFields) UnmarshalCBOR(r io.Reader) (err error) {
	*t = RenamedFields{}

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

	if maj != cbg.MajMap {
		return fmt.Errorf("cbor input should be of type map")
	}

	if extra > cbg.MaxLength {
		return fmt.Errorf("RenamedFields: map struct too large (%d)", extra)
	}

	var name string
	n := extra

	for i := uint64(0); i < n; i++ {

		{
			sval, err := cbg.ReadString(cr)
			if err != nil {
				return err
			}

			name = string(sval)
		}

		switch name {
		// t.Foo (int64) (int64)
		case "foo":
			{
				maj, extra, err := cr.ReadHeader()
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

				t.Foo = int64(extraI)
			}
			// t.Bar (string) (string)
		case "beep":

			{
				sval, err := cbg.ReadString(cr)
				if err != nil {
					return err
				}

				t.Bar = string(sval)
			}

		default:
			// Field doesn't exist on this type, so ignore it
			cbg.ScanForLinks(r, func(cid.Cid) {})
		}
	}

	return nil
}
