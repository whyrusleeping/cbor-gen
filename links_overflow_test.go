package typegen

import (
	"bytes"
	"strings"
	"testing"
)

func TestScanForLinksContainerLengthOverflow(t *testing.T) {
	tests := []struct {
		name string
		data []byte
	}{
		{
			name: "map entry count multiplication",
			data: []byte{0xbb, 0x80, 0, 0, 0, 0, 0, 0, 0}, // 2^63 map pairs
		},
		{
			name: "array entry count addition",
			data: []byte{0x9b, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
		},
		{
			name: "nested array entry count addition",
			data: []byte{0x9b, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xfe, 0x82},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := ScanForLinks(bytes.NewReader(tc.data), nil)
			if err == nil || !strings.Contains(err.Error(), "overflows remaining item count") {
				t.Fatalf("expected container length overflow, got %v", err)
			}
		})
	}
}

func TestScanForLinksNestedContainersWithoutOverflow(t *testing.T) {
	// [ {"x": 1}, [] ]
	data := []byte{0x82, 0xa1, 0x61, 'x', 0x01, 0x80}
	if err := ScanForLinks(bytes.NewReader(data), nil); err != nil {
		t.Fatalf("unexpected scan error: %v", err)
	}
}
