package main

import (
	"io"

	cbg "github.com/whyrusleeping/cbor-gen"
	types "github.com/whyrusleeping/cbor-gen/testing"
)

func main() {
	if err := cbg.WriteTupleEncodersToFile("testing/cbor_gen.go", "testing",
		types.SignedArray{},
		types.SimpleTypeOne{},
		types.SimpleTypeTwo{},
		types.DeferredContainer{},
		types.FixedArrays{},
		types.ThingWithSomeTime{},
		types.BigField{},
		types.IntArray{},
		types.IntAliasArray{},
		types.TupleIntArray{},
		types.TupleIntArrayOptionals{},
		types.IntArrayNewType{},
		types.IntArrayAliasNewType{},
		types.MapTransparentType{},
		types.BigIntContainer{},
		types.GenericStruct[dummy1, dummy2]{},
		types.SubGenericStruct[dummy1, dummy2]{},
	); err != nil {
		panic(err)
	}

	if err := cbg.WriteMapEncodersToFile("testing/cbor_map_gen.go", "testing",
		types.SimpleTypeTree{},
		types.NeedScratchForMap{},
		types.SimpleStructV1{},
		types.SimpleStructV2{},
		types.RenamedFields{},
		types.TestEmpty{},
		types.TestConstField{},
		types.TestCanonicalFieldOrder{},
		types.MapStringString{},
		types.TestSliceNilPreserve{},
		types.StringPtrSlices{},
	); err != nil {
		panic(err)
	}

	err := cbg.Gen{
		MaxArrayLength:  10,
		MaxByteLength:   9,
		MaxStringLength: 8,
	}.WriteTupleEncodersToFile("testing/cbor_options_gen.go", "testing",
		types.LimitedStruct{},
	)
	if err != nil {
		panic(err)
	}

	err = cbg.Gen{
		MaxArrayLength:  10,
		MaxByteLength:   9,
		MaxStringLength: 10000,
	}.WriteTupleEncodersToFile("testing/cbor_options_gen2.go", "testing",
		types.LongString{},
	)
	if err != nil {
		panic(err)
	}
}

// dummy generic types that cbor-gen will replace
type (
	dummy1 int64
	dummy2 int64
)

func (d dummy1) New() dummy1                   { return dummy1(0) }
func (d dummy1) Equals(dummy1) bool            { return false }
func (d dummy1) MarshalCBOR(io.Writer) error   { return nil }
func (d dummy1) UnmarshalCBOR(io.Reader) error { return nil }

func (d dummy2) New() dummy2                   { return dummy2(0) }
func (d dummy2) Equals(dummy2) bool            { return false }
func (d dummy2) MarshalCBOR(io.Writer) error   { return nil }
func (d dummy2) UnmarshalCBOR(io.Reader) error { return nil }
