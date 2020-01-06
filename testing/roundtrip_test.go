package testing

import (
	"bytes"
	"math/rand"
	"reflect"
	"testing"
	"testing/quick"

	cbg "github.com/whyrusleeping/cbor-gen"
)

func TestSimpleTypeOne(t *testing.T) {
	testTypeRoundtrips(t, reflect.TypeOf(SimpleTypeOne{}))
}

func TestSimpleTypeTwo(t *testing.T) {
	testTypeRoundtrips(t, reflect.TypeOf(SimpleTypeTwo{}))
}

func TestSimpleTypeTree(t *testing.T) {
	testTypeRoundtrips(t, reflect.TypeOf(SimpleTypeTree{}))
}

func testTypeRoundtrips(t *testing.T, typ reflect.Type) {
	r := rand.New(rand.NewSource(56887))
	for i := 0; i < 1000; i++ {
		val, ok := quick.Value(typ, r)
		if !ok {
			t.Fatal("failed to generate test value")
		}

		obj := val.Addr().Interface().(cbg.CBORMarshaler)

		buf := new(bytes.Buffer)
		if err := obj.MarshalCBOR(buf); err != nil {
			t.Fatal("i guess its fine to fail marshaling")
		}

		enc := buf.Bytes()

		nobj := reflect.New(typ).Interface().(cbg.CBORUnmarshaler)
		if err := nobj.UnmarshalCBOR(bytes.NewReader(enc)); err != nil {
			t.Logf("got bad bytes: %x", enc)
			t.Fatal("failed to round trip object: ", err)
		}

		if !reflect.DeepEqual(obj, nobj) {
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
}
