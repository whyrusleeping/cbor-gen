package testing

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"math/rand"
	"reflect"
	"testing"
	"testing/quick"
	"time"

	"github.com/ipfs/go-cid"

	"github.com/google/go-cmp/cmp"
	cbg "github.com/whyrusleeping/cbor-gen"
)

var alwaysEqual = cmp.Comparer(func(_, _ interface{}) bool { return true })

// This option handles slices and maps of any type.
var alwaysEqualOpt = cmp.FilterValues(func(x, y interface{}) bool {
	vx, vy := reflect.ValueOf(x), reflect.ValueOf(y)
	return (vx.IsValid() && vy.IsValid() && vx.Type() == vy.Type()) &&
		(vx.Kind() == reflect.Slice || vx.Kind() == reflect.Map) &&
		(vx.Len() == 0 && vy.Len() == 0)
}, alwaysEqual)

func TestSimpleSigned(t *testing.T) {
	testTypeRoundtrips(t, reflect.TypeOf(SignedArray{}))
}

func TestSimpleTypeOne(t *testing.T) {
	testTypeRoundtrips(t, reflect.TypeOf(SimpleTypeOne{}))
}

func TestSimpleTypeTwo(t *testing.T) {
	testTypeRoundtrips(t, reflect.TypeOf(SimpleTypeTwo{}))
}

func TestSimpleTypeTree(t *testing.T) {
	testTypeRoundtrips(t, reflect.TypeOf(SimpleTypeTree{}))
}

func TestNeedScratchForMap(t *testing.T) {
	testTypeRoundtrips(t, reflect.TypeOf(NeedScratchForMap{}))
}

func testValueRoundtrip(t *testing.T, obj cbg.CBORMarshaler, nobj cbg.CBORUnmarshaler) {

	buf := new(bytes.Buffer)
	if err := obj.MarshalCBOR(buf); err != nil {
		t.Fatal("i guess its fine to fail marshaling")
	}

	enc := buf.Bytes()

	if err := nobj.UnmarshalCBOR(bytes.NewReader(enc)); err != nil {
		t.Logf("got bad bytes: %x", enc)
		t.Fatal("failed to round trip object: ", err)
	}

	if !cmp.Equal(obj, nobj, alwaysEqualOpt) {
		t.Logf("%#v != %#v", obj, nobj)
		t.Log("not equal after round trip!")
	}

	nbuf := new(bytes.Buffer)
	if err := nobj.(cbg.CBORMarshaler).MarshalCBOR(nbuf); err != nil {
		t.Fatal("failed to remarshal object: ", err)
	}

	if !bytes.Equal(nbuf.Bytes(), enc) {
		t.Fatalf("objects encodings different: %x != %x", nbuf.Bytes(), enc)
	}

}

func testTypeRoundtrips(t *testing.T, typ reflect.Type) {
	r := rand.New(rand.NewSource(56887))
	for i := 0; i < 1000; i++ {
		val, ok := quick.Value(typ, r)
		if !ok {
			t.Fatal("failed to generate test value")
		}

		obj := val.Addr().Interface().(cbg.CBORMarshaler)
		nobj := reflect.New(typ).Interface().(cbg.CBORUnmarshaler)
		testValueRoundtrip(t, obj, nobj)
	}
}

func TestDeferredContainer(t *testing.T) {
	zero := &DeferredContainer{}
	recepticle := &DeferredContainer{}
	testValueRoundtrip(t, zero, recepticle)
}

func TestNilValueDeferredUnmarshaling(t *testing.T) {
	var zero DeferredContainer
	zero.Deferred = &cbg.Deferred{Raw: []byte{0xf6}}

	buf := new(bytes.Buffer)
	if err := zero.MarshalCBOR(buf); err != nil {
		t.Fatal(err)
	}

	var n DeferredContainer
	if err := n.UnmarshalCBOR(buf); err != nil {
		t.Fatal(err)
	}

	if n.Deferred == nil {
		t.Fatal("shouldnt be nil!")
	}
}

func TestFixedArrays(t *testing.T) {
	zero := &FixedArrays{}
	recepticle := &FixedArrays{}
	testValueRoundtrip(t, zero, recepticle)
}

func TestTimeIsh(t *testing.T) {
	val := &ThingWithSomeTime{
		When:    cbg.CborTime(time.Now()),
		Stuff:   1234,
		CatName: "hank",
	}

	buf := new(bytes.Buffer)
	if err := val.MarshalCBOR(buf); err != nil {
		t.Fatal(err)
	}

	out := ThingWithSomeTime{}
	if err := out.UnmarshalCBOR(buf); err != nil {
		t.Fatal(err)
	}

	if out.When.Time().UnixNano() != val.When.Time().UnixNano() {
		t.Fatal("time didnt round trip properly", out.When.Time(), val.When.Time())
	}

	if out.Stuff != val.Stuff {
		t.Fatal("no")
	}

	if out.CatName != val.CatName {
		t.Fatal("no")
	}

	b, err := json.Marshal(val)
	if err != nil {
		t.Fatal(err)
	}

	var out2 ThingWithSomeTime
	if err := json.Unmarshal(b, &out2); err != nil {
		t.Fatal(err)
	}

	if out2.When != out.When {
		t.Fatal(err)
	}

}

func TestLessToMoreFieldsRoundTrip(t *testing.T) {
	dummyCid, _ := cid.Parse("bafkqaaa")
	simpleTypeOne := SimpleTypeOne{
		Foo:     "foo",
		Value:   1,
		Binary:  []byte("bin"),
		Signed:  -1,
		NString: "namedstr",
	}
	obj := &SimpleStructV1{
		OldStr:    "hello",
		OldBytes:  []byte("bytes"),
		OldNum:    10,
		OldPtr:    &dummyCid,
		OldMap:    map[string]SimpleTypeOne{"first": simpleTypeOne},
		OldArray:  []SimpleTypeOne{simpleTypeOne},
		OldStruct: simpleTypeOne,
	}

	buf := new(bytes.Buffer)
	if err := obj.MarshalCBOR(buf); err != nil {
		t.Fatal("failed marshaling", err)
	}

	enc := buf.Bytes()

	nobj := SimpleStructV2{}
	if err := nobj.UnmarshalCBOR(bytes.NewReader(enc)); err != nil {
		t.Logf("got bad bytes: %x", enc)
		t.Fatal("failed to round trip object: ", err)
	}

	if obj.OldStr != nobj.OldStr {
		t.Fatal("mismatch ", obj.OldStr, " != ", nobj.OldStr)
	}
	if nobj.NewStr != "" {
		t.Fatal("expected field to be zero value")
	}

	if obj.OldNum != nobj.OldNum {
		t.Fatal("mismatch ", obj.OldNum, " != ", nobj.OldNum)
	}
	if nobj.NewNum != 0 {
		t.Fatal("expected field to be zero value")
	}

	if !bytes.Equal(obj.OldBytes, nobj.OldBytes) {
		t.Fatal("mismatch ", obj.OldBytes, " != ", nobj.OldBytes)
	}
	if nobj.NewBytes != nil {
		t.Fatal("expected field to be zero value")
	}

	if *obj.OldPtr != *nobj.OldPtr {
		t.Fatal("mismatch ", obj.OldPtr, " != ", nobj.OldPtr)
	}
	if nobj.NewPtr != nil {
		t.Fatal("expected field to be zero value")
	}

	if !cmp.Equal(obj.OldMap, nobj.OldMap) {
		t.Fatal("mismatch map marshal / unmarshal")
	}
	if len(nobj.NewMap) != 0 {
		t.Fatal("expected field to be zero value")
	}

	if !cmp.Equal(obj.OldArray, nobj.OldArray) {
		t.Fatal("mismatch array marshal / unmarshal")
	}
	if len(nobj.NewArray) != 0 {
		t.Fatal("expected field to be zero value")
	}

	if !cmp.Equal(obj.OldStruct, nobj.OldStruct) {
		t.Fatal("mismatch struct marshal / unmarshal")
	}
	if !cmp.Equal(nobj.NewStruct, SimpleTypeOne{}) {
		t.Fatal("expected field to be zero value")
	}
}

func TestMoreToLessFieldsRoundTrip(t *testing.T) {
	dummyCid1, _ := cid.Parse("bafkqaaa")
	dummyCid2, _ := cid.Parse("bafkqaab")
	simpleType1 := SimpleTypeOne{
		Foo:     "foo",
		Value:   1,
		Binary:  []byte("bin"),
		Signed:  -1,
		NString: "namedstr",
	}
	simpleType2 := SimpleTypeOne{
		Foo:     "bar",
		Value:   2,
		Binary:  []byte("bin2"),
		Signed:  -2,
		NString: "namedstr2",
	}
	obj := &SimpleStructV2{
		OldStr:    "oldstr",
		NewStr:    "newstr",
		OldBytes:  []byte("oldbytes"),
		NewBytes:  []byte("newbytes"),
		OldNum:    10,
		NewNum:    11,
		OldPtr:    &dummyCid1,
		NewPtr:    &dummyCid2,
		OldMap:    map[string]SimpleTypeOne{"foo": simpleType1},
		NewMap:    map[string]SimpleTypeOne{"bar": simpleType2},
		OldArray:  []SimpleTypeOne{simpleType1},
		NewArray:  []SimpleTypeOne{simpleType1, simpleType2},
		OldStruct: simpleType1,
		NewStruct: simpleType2,
	}

	buf := new(bytes.Buffer)
	if err := obj.MarshalCBOR(buf); err != nil {
		t.Fatal("failed marshaling", err)
	}

	enc := buf.Bytes()

	nobj := SimpleStructV1{}
	if err := nobj.UnmarshalCBOR(bytes.NewReader(enc)); err != nil {
		t.Logf("got bad bytes: %x", enc)
		t.Fatal("failed to round trip object: ", err)
	}

	if obj.OldStr != nobj.OldStr {
		t.Fatal("mismatch", obj.OldStr, " != ", nobj.OldStr)
	}
	if obj.OldNum != nobj.OldNum {
		t.Fatal("mismatch ", obj.OldNum, " != ", nobj.OldNum)
	}
	if !bytes.Equal(obj.OldBytes, nobj.OldBytes) {
		t.Fatal("mismatch ", obj.OldBytes, " != ", nobj.OldBytes)
	}
	if *obj.OldPtr != *nobj.OldPtr {
		t.Fatal("mismatch ", obj.OldPtr, " != ", nobj.OldPtr)
	}
	if !cmp.Equal(obj.OldMap, nobj.OldMap) {
		t.Fatal("mismatch map marshal / unmarshal")
	}
	if !cmp.Equal(obj.OldArray, nobj.OldArray) {
		t.Fatal("mismatch array marshal / unmarshal")
	}
	if !cmp.Equal(obj.OldStruct, nobj.OldStruct) {
		t.Fatal("mismatch struct marshal / unmarshal")
	}
}

func TestErrUnexpectedEOF(t *testing.T) {
	err := quick.Check(func(val SimpleTypeTwo, endIdx uint) bool {
		return t.Run("quickcheck", func(t *testing.T) {
			buf := new(bytes.Buffer)
			if err := val.MarshalCBOR(buf); err != nil {
				t.Error(err)
			}

			enc := buf.Bytes()
			originalLen := len(enc)
			endIdx = endIdx % uint(len(enc))
			enc = enc[:endIdx]

			nobj := SimpleTypeTwo{}
			err := nobj.UnmarshalCBOR(bytes.NewReader(enc))
			t.Logf("endIdx=%v, originalLen=%v", endIdx, originalLen)
			if int(endIdx) == originalLen && err != nil {
				t.Fatal("failed to round trip object: ", err)
			} else if endIdx == 0 && !errors.Is(err, io.EOF) {
				t.Fatal("expected EOF got", err)
			} else if endIdx != 0 && err == io.EOF {
				t.Fatal("did not expect EOF but got it")
			}
		})

	}, &quick.Config{MaxCount: 1000})

	if err != nil {
		t.Error(err)
	}

}
