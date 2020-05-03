package testing

import (
	"io/ioutil"
	"math/rand"
	"reflect"
	"testing"
	"testing/quick"
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
