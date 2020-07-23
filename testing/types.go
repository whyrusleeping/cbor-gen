package testing

import (
	cbg "github.com/whyrusleeping/cbor-gen"
)

const Thingc = 3

type NamedNumber uint64
type NamedString uint64

type SignedArray struct {
	Signed []uint64
}

type SimpleTypeOne struct {
	Foo     string
	Value   uint64
	Binary  []byte
	Signed  int64
	NString NamedString
}

type SimpleTypeTwo struct {
	Stuff        *SimpleTypeTwo
	Others       []uint64
	SignedOthers []int64
	Test         [][]byte
	Dog          string
	Numbers      []NamedNumber
	Pizza        *uint64
	PointyPizza  *NamedNumber
	Arrrrrghay   [Thingc]SimpleTypeOne
}

type SimpleTypeTree struct {
	Stuff                            *SimpleTypeTree
	Stufff                           *SimpleTypeTwo
	Others                           []uint64
	Test                             [][]byte
	Dog                              string
	SixtyThreeBitIntegerWithASignBit int64
	NotPizza                         *uint64
}

type DeferredContainer struct {
	Stuff    *SimpleTypeOne
	Deferred *cbg.Deferred
	Value    uint64
}

type FixedArrays struct {
	Bytes  [20]byte
	Uint8  [20]uint8
	Uint64 [20]uint64
}
