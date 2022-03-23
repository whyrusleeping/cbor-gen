package testing

import (
	"bytes"
	"io"
	"io/ioutil"
	"math/rand"
	"reflect"
	"testing"
	"testing/quick"

	"github.com/ipfs/go-cid"
	cbg "github.com/whyrusleeping/cbor-gen"
)

func BenchmarkMarshaling(b *testing.B) {
	r := rand.New(rand.NewSource(56887))
	val, ok := quick.Value(reflect.TypeOf(SimpleTypeTwo{}), r)
	if !ok {
		b.Fatal("failed to construct type")
	}

	tt := val.Interface().(SimpleTypeTwo)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		if err := tt.MarshalCBOR(ioutil.Discard); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkUnmarshaling(b *testing.B) {
	r := rand.New(rand.NewSource(123456))
	val, ok := quick.Value(reflect.TypeOf(SimpleTypeTwo{}), r)
	if !ok {
		b.Fatal("failed to construct type")
	}

	tt := val.Interface().(SimpleTypeTwo)

	buf := new(bytes.Buffer)
	if err := tt.MarshalCBOR(buf); err != nil {
		b.Fatal(err)
	}

	reader := bytes.NewReader(buf.Bytes())

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		reader.Seek(0, io.SeekStart)
		var tt SimpleTypeTwo
		if err := tt.UnmarshalCBOR(reader); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkLinkScan(b *testing.B) {
	r := rand.New(rand.NewSource(123456))
	val, ok := quick.Value(reflect.TypeOf(SimpleTypeTwo{}), r)
	if !ok {
		b.Fatal("failed to construct type")
	}

	tt := val.Interface().(SimpleTypeTwo)

	buf := new(bytes.Buffer)
	if err := tt.MarshalCBOR(buf); err != nil {
		b.Fatal(err)
	}

	reader := bytes.NewReader(buf.Bytes())

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		reader.Seek(0, io.SeekStart)
		if err := cbg.ScanForLinks(reader, func(cid.Cid) {}); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkDeferred(b *testing.B) {
	r := rand.New(rand.NewSource(123456))
	val, ok := quick.Value(reflect.TypeOf(SimpleTypeTwo{}), r)
	if !ok {
		b.Fatal("failed to construct type")
	}

	tt := val.Interface().(SimpleTypeTwo)

	buf := new(bytes.Buffer)
	if err := tt.MarshalCBOR(buf); err != nil {
		b.Fatal(err)
	}

	var (
		deferred cbg.Deferred
		reader   = bytes.NewReader(buf.Bytes())
	)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		reader.Seek(0, io.SeekStart)
		if err := deferred.UnmarshalCBOR(reader); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkMapMarshaling(b *testing.B) {
	r := rand.New(rand.NewSource(56887))
	val, ok := quick.Value(reflect.TypeOf(SimpleTypeTree{}), r)
	if !ok {
		b.Fatal("failed to construct type")
	}

	tt := val.Interface().(SimpleTypeTree)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		if err := tt.MarshalCBOR(ioutil.Discard); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkMapUnmarshaling(b *testing.B) {
	r := rand.New(rand.NewSource(123456))
	val, ok := quick.Value(reflect.TypeOf(SimpleTypeTree{}), r)
	if !ok {
		b.Fatal("failed to construct type")
	}

	tt := val.Interface().(SimpleTypeTree)

	buf := new(bytes.Buffer)
	if err := tt.MarshalCBOR(buf); err != nil {
		b.Fatal(err)
	}

	reader := bytes.NewReader(buf.Bytes())

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		reader.Seek(0, io.SeekStart)
		var tt SimpleTypeTree
		if err := tt.UnmarshalCBOR(reader); err != nil {
			b.Fatal(err)
		}
	}
}
