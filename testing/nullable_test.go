package testing

import (
	"bytes"
	"testing"
)

func TestNotNullMarshal(t *testing.T) {
	var not NotNull

	buf := new(bytes.Buffer)
	if err := not.MarshalCBOR(buf); err == nil {
		t.Fatal("should not have serialized null fields")
	}
	buf.Reset()

	not.Always = new(SimpleTypeOne)
	if err := not.MarshalCBOR(buf); err == nil {
		t.Fatal("should not have serialized null fields")
	}
	buf.Reset()

	not.Value = new(uint64)
	if err := not.MarshalCBOR(buf); err != nil {
		t.Fatal(err)
	}
	buf.Reset()

	not.Always = nil
	if err := not.MarshalCBOR(buf); err == nil {
		t.Fatal("should not have serialized null fields")
	}
	buf.Reset()

	not.Map = make(map[string]*NotNull)
	not.Map["foo"] = nil
	if err := not.MarshalCBOR(buf); err == nil {
		t.Fatal("should not have serialized null fields")
	}
	buf.Reset()

	not.Map["foo"] = &NotNull{Always: new(SimpleTypeOne), Value: new(uint64)}
	if err := not.MarshalCBOR(buf); err == nil {
		t.Fatal(err)
	}
	buf.Reset()

	not.Slice = append(not.Slice, nil)
	if err := not.MarshalCBOR(buf); err == nil {
		t.Fatal("should not have serialized null fields")
	}
	buf.Reset()

	not.Slice[0] = &NotNull{Always: new(SimpleTypeOne), Value: new(uint64)}
	if err := not.MarshalCBOR(buf); err == nil {
		t.Fatal(err)
	}
	buf.Reset()
}

func TestNotNullUnmarshal(t *testing.T) {
	var (
		yes YesNull
		not NotNull
	)

	buf := new(bytes.Buffer)
	buf.Reset()
	if err := yes.MarshalCBOR(buf); err != nil {
		t.Fatal(err)
	}

	if err := not.UnmarshalCBOR(buf); err == nil {
		t.Fatal("shouldn't decode null fields")
	}

	yes.Value = new(uint64)
	buf.Reset()
	if err := yes.MarshalCBOR(buf); err != nil {
		t.Fatal(err)
	}

	if err := not.UnmarshalCBOR(buf); err == nil {
		t.Fatal("shouldn't decode null fields")
	}

	yes.Always = new(SimpleTypeOne)
	buf.Reset()
	if err := yes.MarshalCBOR(buf); err != nil {
		t.Fatal(err)
	}
	if err := not.UnmarshalCBOR(buf); err != nil {
		t.Fatal(err)
	}

	yes.Map = make(map[string]*YesNull)
	yes.Map["foo"] = nil
	buf.Reset()
	if err := yes.MarshalCBOR(buf); err != nil {
		t.Fatal(err)
	}
	if err := not.UnmarshalCBOR(buf); err == nil {
		t.Fatal("has null value")
	}

	yes.Map["foo"] = &YesNull{Always: new(SimpleTypeOne), Value: new(uint64)}
	buf.Reset()
	if err := yes.MarshalCBOR(buf); err != nil {
		t.Fatal(err)
	}
	if err := not.UnmarshalCBOR(buf); err != nil {
		t.Fatal(err)
	}

	yes.Slice = append(yes.Slice, nil)
	buf.Reset()
	if err := yes.MarshalCBOR(buf); err != nil {
		t.Fatal(err)
	}
	if err := not.UnmarshalCBOR(buf); err == nil {
		t.Fatal("has null value")
	}

	yes.Slice[0] = &YesNull{Always: new(SimpleTypeOne), Value: new(uint64)}
	buf.Reset()
	if err := yes.MarshalCBOR(buf); err != nil {
		t.Fatal(err)
	}
	if err := not.UnmarshalCBOR(buf); err != nil {
		t.Fatal(err)
	}
}
