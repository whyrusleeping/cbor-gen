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
		types.NotNull{},
		types.YesNull{},
	); err != nil {
		panic(err)
	}

	if err := cbg.WriteMapEncodersToFile("testing/cbor_map_gen.go", "testing",
		types.SimpleTypeTree{},
	); err != nil {
		panic(err)
	}
}
