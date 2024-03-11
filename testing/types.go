package testing

import (
	"math/rand"
	"reflect"

	"github.com/ipfs/go-cid"

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
	Strings []string
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
	StringPtr                        *string
	BoolPtr                          *bool
}

type SimpleStructV1 struct {
	OldStr         string
	OldBytes       []byte
	OldNum         uint64
	OldPtr         *cid.Cid
	OldMap         map[string]SimpleTypeOne
	OldArray       []SimpleTypeOne
	OldStruct      SimpleTypeOne
	OldCidArray    []cid.Cid
	OldCidPtrArray []*cid.Cid
}

type SimpleStructV2 struct {
	OldStr string
	NewStr string

	OldBytes []byte
	NewBytes []byte

	OldNum uint64
	NewNum uint64

	OldPtr *cid.Cid
	NewPtr *cid.Cid

	OldMap map[string]SimpleTypeOne
	NewMap map[string]SimpleTypeOne

	OldArray []SimpleTypeOne
	NewArray []SimpleTypeOne

	OldStruct SimpleTypeOne
	NewStruct SimpleTypeOne
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

type ThingWithSomeTime struct {
	When    cbg.CborTime
	Stuff   int64
	CatName string
}

// Do not add fields to this type.
type NeedScratchForMap struct {
	Thing bool
}

type RenamedFields struct {
	Foo int64  `cborgen:"foo"`
	Bar string `cborgen:"beep"`
}

type BigField struct {
	LargeBytes []byte `cborgen:"maxlen=10000000"`
}

type TestEmpty struct {
	Foo  *string `cborgen:"omitempty"`
	Beep string  `cborgen:"omitempty"`
	Cat  int64
}

type TestConstField struct {
	Cats  string `cborgen:"const=dogsdrool"`
	Thing int64
}

type TestCanonicalFieldOrder struct {
	Foo   int64  `cborgen:"foo"`
	Bar   string `cborgen:"beep"`
	Drond int64
	Zp    string `cborgen:"ap"`
}

type MapStringString struct {
	Snorkleblump map[string]string
}

type TestSliceNilPreserve struct {
	Cat      string
	Stuff    []uint64
	Not      []uint64 `cborgen:"preservenil"`
	Other    []byte
	NotOther []byte `cborgen:"preservenil"`
	Beep     int64
}

type IntAlias int64

type IntArray struct {
	Ints []int64 `cborgen:"transparent"`
}
type IntAliasArray struct {
	Ints []IntAlias `cborgen:"transparent"`
}

type IntArrayNewType []int64

type IntArrayAliasNewType []IntAlias

type TupleIntArray struct {
	Int1 int64
	Int2 int64
	Int3 int64
}

type TupleIntArrayOptionals struct {
	Int1 *int64
	Int2 int64
	Int3 uint64
	Int4 *uint64
}

type MapTransparentType map[string]string

type LimitedStruct struct {
	Arr  []uint64
	Byts []byte
	Str  string
}

type LongString struct {
	Val string
}

func (ls LongString) Generate(rand *rand.Rand, size int) reflect.Value {
	ols := new(LongString)
	s := make([]byte, 9999)
	rand.Read(s)
	ols.Val = string(s)
	return reflect.ValueOf(ols).Elem()
}
