package main

import (
	cbg "github.com/whyrusleeping/cbor-gen"
	types "github.com/whyrusleeping/cbor-gen/testing"
)

func main() {
	if err := cbg.WriteTupleEncodersToFile("testing/cbor_gen.go", "types",
		types.SimpleTypeOne{},
		types.SimpleTypeTwo{},
	); err != nil {
		panic(err)
	}
}
