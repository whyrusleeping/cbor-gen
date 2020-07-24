package testing

import (
	cbg "github.com/whyrusleeping/cbor-gen"
)

const Thingc = 3

type NamedNumber uint64
type NamedString string

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
	Stuff        *SimpleTypeTwo `cbg:"nullable"`
	Others       []uint64
	SignedOthers []int64
	Test         [][]byte
	Dog          string
	Numbers      []NamedNumber
	Pizza        *uint64      `cbg:"nullable"`
	PointyPizza  *NamedNumber `cbg:"nullable"`
	Arrrrrghay   [Thingc]SimpleTypeOne
}

type SimpleTypeTree struct {
	Stuff                            *SimpleTypeTree `cbg:"nullable"`
	Stufff                           *SimpleTypeTwo  `cbg:"nullable"`
	Others                           []uint64
	Test                             [][]byte
	Dog                              string
	SixtyThreeBitIntegerWithASignBit int64
	NotPizza                         *uint64 `cbg:"nullable"`
}

type DeferredContainer struct {
	Stuff    *SimpleTypeOne `cbg:"nullable"`
	Deferred *cbg.Deferred  `cbg:"nullable"`
	Value    uint64
}

type NotNull struct {
	Always *SimpleTypeOne
	Value  *uint64
	Slice  []*NotNull
	Map    map[string]*NotNull
}

type YesNull struct {
	Always *SimpleTypeOne      `cbg:"nullable"`
	Value  *uint64             `cbg:"nullable"`
	Slice  []*YesNull          `cbg:"nullable"`
	Map    map[string]*YesNull `cbg:"nullable"`
}

type FixedArrays struct {
	Bytes  [20]byte
	Uint8  [20]uint8
	Uint64 [20]uint64
}
