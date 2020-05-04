package testing

import (
	"bytes"
	"encoding/hex"
	"io"
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

func BenchmarkUnmarshaling(b *testing.B) {
	hx := "8989f68080807859f099a586f0908093f1af9fb6f3a0ad82e8aaa0efbfbdf1b88688f29d8aaeecacabf0a4be94f19295a0f19b9081f0b6bf8ff3ad83a6f09e9ca7f2be8a8bf187a8a0f280a8b3f4899a9bf181afb1f0bca2b0f1b5ab9bf2b0a3ac80f6f68384787af2b1a082f3b99d98f1adb9b6f3b9868df29fbbb0f3858791f3b5b39df2b68e92f2a9bb9af282b4b2f18fba9cf294b8b2f3a0a1a9f0aab1aaf28cb994f09796aef195bc90f488be81e0a59af2928183f1a0a4abf393bbbae39d8df28fb287f0bf8fa6f0b79a89f188babcf395b0b1f29ebab7f2b0a091f29db48f1b527ef13ee4f5321a403b2130370299eeb937847848f287b8b9f1ad9e90f1b1b9bbf18ebc91f3908583f0be9ab3f2aca8abf0a8acadf380a7abf293a8aaf1b2b6a6f3b89587f3809fadf3a39f97f3a8b48cf3b299bff19cab9df28399a01b374797708d2015d3401b6bfb7066c509754c8478a1f0ab8f96f287988df297aea7f3afa699f3859788f2a2b2b8f2b681a6f29a95a4f382978cf396b183e2acbdf39cbdb5f0b99b94f1a2baaaf1ba89b0f3a8a7bbf397bdabf3af8c83f1b38ebef0beb1a0f3939f83f0b9ad90f1acb597f0b49eb0f29ab3a3f480808ef39b878ae5989ff0a7b789f48981b6f281aba6f2a9ad88f09fb395f0aa95adf0a1a59ff38a8d97f397b7b0eebfa9f2a5ab87f2afa7b8f0b992b81bab566703ac0b139c401b463b0320db277de1841b73dd7cd1861ff4561bc0256739761d28dd1b39c9019ac37c08721b2f08fbf368bf7f94813b075c40eb7f66e0488078b4f288898bf486bc90e9b8bdf180ab8ff39db1a7f0b1afb8f3ab9fa6f1b19182f189bdaff3bf80a5f1a18fb4f39c99a9f0ba839af2adb88fe4bd9df39bba8bf28ebf9ef2b3a783f0b6b395f197be84f3a1998af1b0898bf3b0b08ef1b49b94f094b59df19dbfa6f2aa8494f48ba0b2e28181f1a08999f2b3ac81eaa689f1bb80bbf2ae918bf0a19397f1a19d9cf3b095b5f1b4baa2f0b7ad92f3ab8c8ef38fab92f489b499f18d9899f0b5bcb5f2a3b6a5f2a1acb1831be5bdbd1384238b4b1b8a95991fbf9ca8d11baf61be2ac6477c7d1ba1d9dac0cecd182d1b4175138c8c7fbb4e8384785ef0a7adaef39b8a8ff1b79bacc693f1948ebef0938383f48aa6abf09ab684f1ba8c89f188b091f18ab2b5f1ac8484f2b7b089f18b97bdf1838aacf397ad98f0b9a8aff394a2a3f39eb6bff09ab8bef39189bef18f89aaf3aca982f29381901b7f0bf8763d569f3b403b75c40e5c6163108084786af1879fa8f2b4af9ef3a8b3b6f3b0be93f0aba9a8f0a1a698f3b6a7a7e6adbef1a8849bf28087a3f3b89f82f38caab6f0b7b09ff1bf938ff0a0b1aff2b79691f0a5b29bf4858896f484a5abf393bbbbf3a2b8bdf29393a6eba180f1a1b3b0f29da098f1b09ca7f3bda2901b6f088c64a0854512401b564b5898ca46ac958467f29c9188eea5941bb7b58825b1edf1ee403b6dcbd95c52f6ca33"

	d, err := hex.DecodeString(hx)
	if err != nil {
		b.Fatal(err)
	}

	buf := bytes.NewReader(d)
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		buf.Seek(0, io.SeekStart)
		var tt SimpleTypeTwo
		if err := tt.UnmarshalCBOR(buf); err != nil {
			b.Fatal(err)
		}
	}

}
