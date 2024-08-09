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
		types.CborByteArray{},
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
	dummy1 struct{}
	dummy2 struct{}
)

// dummy generic types that cbor-gen will replace, should be able to be handle both pointer and
// value types
var (
	_ cbg.CBORGeneric[dummy1] = dummy1{}
	_ cbg.CBORGeneric[dummy2] = dummy2{}
)

func (d dummy1) ToCBOR(io.Writer) error             { return nil }
func (d dummy1) FromCBOR(io.Reader) (dummy1, error) { return d, nil }

func (d dummy2) ToCBOR(io.Writer) error             { return nil }
func (d dummy2) FromCBOR(io.Reader) (dummy2, error) { return d, nil }
