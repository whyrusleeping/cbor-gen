package testing

import (
	"bytes"
	"encoding/json"
	"math/rand"
	"reflect"
	"testing"
	"testing/quick"
	"time"

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

func TestMapWithRenames(t *testing.T) {
	testTypeRoundtrips(t, reflect.TypeOf(MapWithRenames{}))
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
