package testing

import (
	"bytes"
	"fmt"
	"testing"
)

func TestOptionalFields(t *testing.T) {
	ints := []int64{1, 2, 3, 4, 5}

	objects := make([][]byte, 5)
	for i := range objects {
		count := i + 1
		t.Run(fmt.Sprintf("length-%d", count), func(t *testing.T) {
			var buf bytes.Buffer
			obj := IntArray{Ints: ints[:count]}
			if err := obj.MarshalCBOR(&buf); err != nil {
				t.Fatal(err)
			}

			var out TupleWithOptionalFields
			err := out.UnmarshalCBOR(&buf)
			switch count {
			case 4:
				if out.Int4 != ints[3] {
					t.Errorf("field 4 should be %d, was %d", ints[3], out.Int4)
				}
				fallthrough
			case 3:
				if out.Int3 != ints[2] {
					t.Errorf("field 4 should be %d, was %d", ints[3], out.Int4)
				}
				fallthrough
			case 2:
				if out.Int2 != ints[1] {
					t.Errorf("field 4 should be %d, was %d", ints[1], out.Int2)
				}
				if out.Int1 != ints[0] {
					t.Errorf("field 4 should be %d, was %d", ints[0], out.Int1)
				}
				if err != nil {
					t.Errorf("expected no error when unmarshaling, got: %s", err)
				}
			case 1, 0:
				if err == nil {
					t.Errorf("expected an error when unmarshaling with too few fields")
				}
			default:
				if err == nil {
					t.Errorf("expected an error when unmarshaling with too many fields")
				}
			}
		})
	}

}
