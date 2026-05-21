package typegen

import "testing"

// TestParseTypeInfoOpaqueWrapper covers the "opaque wrapper" support: a struct
// whose sole unexported field is tagged transparent, optionally with the
// bytes/nullable/parse options.
func TestParseTypeInfoOpaqueWrapper(t *testing.T) {
	type wrapper struct {
		str string `cborgen:"transparent,bytes,nullable,parse=Parse,maxlen=64"`
	}

	gti, err := ParseTypeInfo(wrapper{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !gti.Transparent {
		t.Fatal("expected type to be transparent")
	}
	if len(gti.Fields) != 1 {
		t.Fatalf("expected exactly one field, got %d", len(gti.Fields))
	}

	f := gti.Fields[0]
	if f.Name != "str" {
		t.Fatalf("expected unexported field 'str' to be included, got %q", f.Name)
	}
	if !f.Bytes || !f.Nullable {
		t.Fatalf("expected Bytes and Nullable to be set, got Bytes=%v Nullable=%v", f.Bytes, f.Nullable)
	}
	if f.ParseFunc != "Parse" {
		t.Fatalf("expected ParseFunc 'Parse', got %q", f.ParseFunc)
	}
	if f.MaxLen != 64 {
		t.Fatalf("expected MaxLen 64, got %d", f.MaxLen)
	}
}

// TestParseTypeInfoUntaggedUnexportedSkipped confirms an unexported field with
// no transparent tag is still ignored (existing behavior preserved).
func TestParseTypeInfoUntaggedUnexportedSkipped(t *testing.T) {
	type hidden struct {
		secret string
	}

	gti, err := ParseTypeInfo(hidden{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(gti.Fields) != 0 {
		t.Fatalf("expected unexported untagged field to be skipped, got %d fields", len(gti.Fields))
	}
}

// TestParseTypeInfoOpaqueOptionMisuse confirms the new options are rejected
// outside of their supported context.
func TestParseTypeInfoOpaqueOptionMisuse(t *testing.T) {
	t.Run("bytes without transparent", func(t *testing.T) {
		type bad struct {
			S string `cborgen:"bytes"`
		}
		if _, err := ParseTypeInfo(bad{}); err == nil {
			t.Fatal("expected error for bytes option without transparent")
		}
	})

	t.Run("parse on non-string", func(t *testing.T) {
		type bad struct {
			N uint64 `cborgen:"transparent,parse=Parse"`
		}
		if _, err := ParseTypeInfo(bad{}); err == nil {
			t.Fatal("expected error for parse option on non-string field")
		}
	})
}
