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

			// Pre-fill with garbage. We want optional fields to be reset to their
			// defaults.
			out := TupleWithOptionalFields{
				Int1: 0xf1,
				Int2: 0xf2,
				Int3: 0xf3,
				Int4: 0xf4,
			}
			err := out.UnmarshalCBOR(&buf)
			switch count {
			case 4:
				if out.Int4 != ints[3] {
					t.Errorf("field 4 should be %d, was %d", ints[3], out.Int4)
				}
				fallthrough
			case 3:
				if out.Int3 != ints[2] {
					t.Errorf("field 3 should be %d, was %d", ints[2], out.Int3)
				}
				fallthrough
			case 2:
				if out.Int2 != ints[1] {
					t.Errorf("field 2 should be %d, was %d", ints[1], out.Int2)
				}
				if out.Int1 != ints[0] {
					t.Errorf("field 1 should be %d, was %d", ints[0], out.Int1)
				}
				if (count == 4 || count == 3) && out.Int4 != 0 {
					t.Errorf("expected field 4 to be zero")
				} else if count == 3 && out.Int3 != 0 {
					t.Errorf("expected field 3 to be zero")
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
