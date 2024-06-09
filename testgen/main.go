package main

import (
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
