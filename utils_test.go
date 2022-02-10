package typegen

import (
	"bufio"
	"bytes"
	"encoding/hex"
	"io"
	"testing"

	cid "github.com/ipfs/go-cid"
)

func TestLinkScan(t *testing.T) {
	inp := "82442847c0498ba16131818242000484d82a57000155001266696c2f312f73746f72616765706f776572d82a5827000171a0e40220740d4196aaaee66d8e9b828bc6f9271662096e36782de248e3b3ed28443dbc810040a16131828242000584d82a5818000155001366696c2f312f73746f726167656d61726b6574d82a5827000171a0e402209e59ceb041921650967e8c77d36269c10049140d28d5015165cc8eb897a2555300408242006684d82a52000155000d66696c2f312f6163636f756e74d82a5827000171a0e40220f9556f0d5a735ff53cc327e85a46c7c094028b9da894de73caa88f162429c29d004b000a968163f0a57b400000a16131818242000084d82a51000155000c66696c2f312f73797374656dd82a5827000171a0e4022045b0cfc220ceec5b7c1c62c4d4193d38e4eba48e8815729ce75f9c0ab0e4c1c00040a16131818242006384d82a52000155000d66696c2f312f6163636f756e74d82a5827000171a0e4022045b0cfc220ceec5b7c1c62c4d4193d38e4eba48e8815729ce75f9c0ab0e4c1c00040a16131818242000184d82a4f000155000a66696c2f312f696e6974d82a5827000171a0e4022050f3c45d0e78f04688c6e8cbdc45f71a5cbcec731519ffdcdd92765fc5ba0da30040a16131818242000684d82a581b000155001666696c2f312f76657269666965647265676973747279d82a5827000171a0e40220fd3fee39acd88c8808110d9741149a79939f52c798b6539048e4835aa4d34fd50040a16131818242006584d82a52000155000d66696c2f312f6163636f756e74d82a5827000171a0e402200293716d8503737644624c102d9ba1514d599044a2e6f2038330124bc6f54361004c002116545850052128000000a16131818242000384d82a4f000155000a66696c2f312f63726f6ed82a5827000171a0e4022065d1dad76492ccd5d010197dc26bd5fb07c0cf85ccdcaff21084c0d47bfd17590040a16131818242005084d82a52000155000d66696c2f312f6163636f756e74d82a5827000171a0e402202e12f4eedaac06c2040df4923656f7f2d6c5991b1874a32e9f52b0e48e61d8410040a16131818242006484d82a52000155000d66696c2f312f6163636f756e74d82a5827000171a0e40220bde06af1782cb302e0973658dcb44b2cbd72891598b7a7b02381b563f1dc9c57004c002116545850052128000000a16131818242000284d82a51000155000c66696c2f312f726577617264d82a5827000171a0e4022083c127ddb0ba85f585b06365346eabe5ef98861d85770670dee402dc3810131e004d0004860d8812f0b38878000000"
	inpb, err := hex.DecodeString(inp)
	if err != nil {
		t.Fatal(err)
	}

	var cids []cid.Cid
	if err := ScanForLinks(bytes.NewReader(inpb), func(c cid.Cid) {
		cids = append(cids, c)
	}); err != nil {
		t.Fatal(err)
	}
	t.Log(cids)
}

func TestDeferredMaxLengthSingle(t *testing.T) {
	var header bytes.Buffer
	if err := WriteMajorTypeHeader(&header, MajByteString, ByteArrayMaxLen+1); err != nil {
		t.Fatal("failed to write header")
	}

	var deferred Deferred
	err := deferred.UnmarshalCBOR(&header)
	if err != maxLengthError {
		t.Fatal("deferred: allowed more than the maximum allocation supported")
	}
}

func TestReadByteBuf(t *testing.T) {
	type testCase struct {
		name       string
		reader     io.Reader
		shouldFail bool
	}
	testCases := []testCase{
		{name: "Reader that returns EOF and n bytes read", reader: &testReader1Byte{b: 0x01}, shouldFail: false},
		{name: "Exhausted reader", reader: &testReader1Byte{b: 0x01, emptied: true}, shouldFail: true},
		{name: "Byte buffer", reader: bytes.NewBuffer([]byte{0x01}), shouldFail: false},
		{name: "Empty Byte buffer", reader: bytes.NewBuffer([]byte{}), shouldFail: true},
		{name: "Byte Reader", reader: bytes.NewReader([]byte{0x01}), shouldFail: false},
		{name: "bufio Reader", reader: bufio.NewReader(bytes.NewReader([]byte{0x01})), shouldFail: false},
		{name: "bufio Reader with testReader", reader: bufio.NewReader(&testReader1Byte{b: 0x01}), shouldFail: false},
		{name: "bufio Reader with exhausted testReader", reader: bufio.NewReader(&testReader1Byte{b: 0x01, emptied: true}), shouldFail: true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			b, err := readByteBuf(tc.reader, []byte{0x00})
			if tc.shouldFail && err == nil {
				t.Fatalf("Expected error. Got nil")
			} else if !tc.shouldFail && err != nil {
				t.Fatalf("Expected no error. Got %v", err)
			}

			// readByteBuf should return a nil error with the byte read.
			if err == nil {
				if b != 0x01 {
					t.Fatalf("Expected byte 0x01. Got %x", b)
				}
			}
		})
	}
}

type testReader1Byte struct {
	emptied bool
	b       byte
}

func (tr *testReader1Byte) Read(p []byte) (n int, err error) {
	if tr.emptied {
		return 0, io.EOF
	}

	written, err := bytes.NewReader([]byte{tr.b}).Read(p)
	if written != 1 {
		panic("unreachable. testReader1Byte has a single byte" + err.Error())
	}
	tr.emptied = true
	return 1, io.EOF
}
