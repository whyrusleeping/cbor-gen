package typegen

import (
	"encoding/binary"
	"fmt"
	"io"
	"reflect"
	"strings"
	"text/template"

	cid "github.com/ipfs/go-cid"
)

const (
	MajUnsignedInt = 0
	MajNegativeInt = 1
	MajByteString  = 2
	MajTextString  = 3
	MajArray       = 4
	MajMap         = 5
	MajTag         = 6
	MajFloat       = 7
)

type ByteReader interface {
	io.ByteReader
	io.Reader
}

type Deferred struct {
	Raw []byte
}

func (d *Deferred) MarshalCBOR(w io.Writer) error {
	_, err := w.Write(d.Raw)
	return err
}

func (d *Deferred) UnmarshalCBOR(br ByteReader) error {
	// TODO: theres a more efficient way to implement this method, but for now
	// this is fine
	maj, extra, err := CborReadHeader(br)
	if err != nil {
		return err
	}
	header := CborEncodeMajorType(maj, extra)

	switch maj {
	case MajUnsignedInt:
		d.Raw = header
		return nil
	case MajByteString, MajTextString:
		buf := make([]byte, int(extra)+len(header))
		copy(buf, header)
		if _, err := io.ReadFull(br, buf[len(header):]); err != nil {
			return err
		}

		return nil
	case MajTag:
		sub := new(Deferred)
		if err := sub.UnmarshalCBOR(br); err != nil {
			return err
		}

		d.Raw = append(header, sub.Raw...)
		return nil
	case MajArray:
		d.Raw = header
		for i := 0; i < int(extra); i++ {
			sub := new(Deferred)
			if err := sub.UnmarshalCBOR(br); err != nil {
				return err
			}

			d.Raw = append(d.Raw, sub.Raw...)
		}
		return nil
	default:
		return fmt.Errorf("unhandled deferred cbor type: %d", maj)
	}
}

func CborReadHeader(br ByteReader) (byte, uint64, error) {
	first, err := br.ReadByte()
	if err != nil {
		return 0, 0, err
	}

	maj := (first & 0xe0) >> 5
	low := first & 0x1f

	switch {
	case low < 24:
		return maj, uint64(low), nil
	case low == 24:
		next, err := br.ReadByte()
		if err != nil {
			return 0, 0, err
		}
		return maj, uint64(next), nil
	case low == 25:
		buf := make([]byte, 2)
		if _, err := io.ReadFull(br, buf); err != nil {
			return 0, 0, err
		}
		return maj, uint64(binary.BigEndian.Uint16(buf)), nil
	case low == 26:
		buf := make([]byte, 4)
		if _, err := io.ReadFull(br, buf); err != nil {
			return 0, 0, err
		}
		return maj, uint64(binary.BigEndian.Uint32(buf)), nil
	case low == 27:
		buf := make([]byte, 8)
		if _, err := io.ReadFull(br, buf); err != nil {
			return 0, 0, err
		}
		return maj, binary.BigEndian.Uint64(buf), nil
	default:
		return 0, 0, fmt.Errorf("invalid header: (%x)", first)
	}
}

func CborWriteHeader(w io.Writer, t byte, val uint64) error {
	header := CborEncodeMajorType(t, val)
	if _, err := w.Write(header); err != nil {
		return err
	}

	return nil
}

func CborEncodeMajorType(t byte, l uint64) []byte {
	var b [9]byte
	switch {
	case l < 24:
		b[0] = (t << 5) | byte(l)
		return b[:1]
	case l < (1 << 8):
		b[0] = (t << 5) | 24
		b[1] = byte(l)
		return b[:2]
	case l < (1 << 16):
		b[0] = (t << 5) | 25
		binary.BigEndian.PutUint16(b[1:3], uint16(l))
		return b[:3]
	case l < (1 << 32):
		b[0] = (t << 5) | 26
		binary.BigEndian.PutUint32(b[1:5], uint32(l))
		return b[:5]
	default:
		b[0] = (t << 5) | 27
		binary.BigEndian.PutUint64(b[1:], uint64(l))
		return b[:]
	}
}

func ReadTaggedByteArray(br ByteReader, exptag uint64, maxlen uint64) ([]byte, error) {
	maj, extra, err := CborReadHeader(br)
	if err != nil {
		return nil, err
	}

	if maj != MajTag {
		return nil, fmt.Errorf("expected cbor type 'tag' in input")
	}

	if extra != exptag {
		return nil, fmt.Errorf("expected tag %d", exptag)
	}

	maj, extra, err = CborReadHeader(br)
	if err != nil {
		return nil, err
	}

	if maj != MajByteString {
		return nil, fmt.Errorf("expected cbor type 'byte string' in input")
	}

	if extra > 256*1024 {
		return nil, fmt.Errorf("string in cbor input too long")
	}

	buf := make([]byte, extra)
	if _, err := io.ReadFull(br, buf); err != nil {
		return nil, err
	}

	return buf, nil

}

func ReadCid(br ByteReader) (cid.Cid, error) {
	buf, err := ReadTaggedByteArray(br, 42, 512)
	if err != nil {
		return cid.Undef, err
	}

	if len(buf) == 0 {
		return cid.Undef, nil
	}

	if len(buf) < 2 {
		return cid.Undef, fmt.Errorf("cbor serialized CIDs must have at least two bytes")
	}

	if buf[0] != 0 {
		return cid.Undef, fmt.Errorf("cbor serialized CIDs must have binary multibase")
	}

	return cid.Cast(buf[1:])
}

func WriteCid(w io.Writer, c cid.Cid) error {
	if err := CborWriteHeader(w, MajTag, 42); err != nil {
		return err
	}
	if c == cid.Undef {
		return CborWriteHeader(w, MajByteString, 0)

	}

	if err := CborWriteHeader(w, MajByteString, uint64(len(c.Bytes())+1)); err != nil {
		return err
	}

	// that binary multibase prefix...
	if _, err := w.Write([]byte{0}); err != nil {
		return err
	}

	if _, err := w.Write(c.Bytes()); err != nil {
		return err
	}

	return nil
}

func doTemplate(w io.Writer, info interface{}, templ string) error {
	t := template.Must(template.New("").
		Funcs(template.FuncMap{}).Parse(templ))

	return t.Execute(w, info)
}

func PrintHeaderAndUtilityMethods(w io.Writer, pkg string) error {
	data := struct {
		Package string
	}{pkg}
	return doTemplate(w, data, `package {{ .Package }}

import (
	"fmt"
	"io"
	cbg "github.com/whyrusleeping/cbor-gen"
)

/* This file was generated by github.com/whyrusleeping/cbor-gen */

`)
}

type Field struct {
	Name    string
	Pointer bool
	Type    reflect.Type
}

type GenTypeInfo struct {
	Name   string
	Fields []Field
}

func nameIsExported(name string) bool {
	return strings.ToUpper(name[0:1]) == name[0:1]
}

func ParseTypeInfo(i interface{}) (*GenTypeInfo, error) {
	t := reflect.TypeOf(i)

	out := GenTypeInfo{
		Name: t.Name(),
	}

	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		if !nameIsExported(f.Name) {
			continue
		}

		ft := f.Type
		var pointer bool
		if ft.Kind() == reflect.Ptr {
			ft = ft.Elem()
			pointer = true
		}

		out.Fields = append(out.Fields, Field{
			Name:    f.Name,
			Pointer: pointer,
			Type:    ft,
		})
	}

	return &out, nil
}

func (gti GenTypeInfo) Header() []byte {
	return CborEncodeMajorType(MajArray, uint64(len(gti.Fields)))
}

func (gti GenTypeInfo) HeaderAsByteString() string {
	h := gti.Header()
	s := "[]byte{"
	for _, b := range h {
		s += fmt.Sprintf("%d,", b)
	}
	s += "}"
	return s
}

func emitCborMarshalStringField(w io.Writer, f Field) error {
	if f.Pointer {
		return fmt.Errorf("pointers to strings not supported")
	}

	return doTemplate(w, f, `
	if _, err := w.Write(cbg.CborEncodeMajorType(cbg.MajTextString, uint64(len(t.{{ .Name }})))); err != nil {
		return err
	}
	if _, err := w.Write([]byte(t.{{ .Name }})); err != nil {
		return err
	}
`)
}
func emitCborMarshalStructField(w io.Writer, f Field) error {
	fname := f.Type.PkgPath() + "." + f.Type.Name()
	switch fname {
	case "math/big.Int":
		return doTemplate(w, f, `
	{
		if err := cbg.CborWriteHeader(w, cbg.MajTag, 2); err != nil {
			return err
		}
		var b []byte
		if t.{{ .Name }} != nil {
			b = t.{{ .Name }}.Bytes()
		}

		if err := cbg.CborWriteHeader(w, cbg.MajByteString, uint64(len(b))); err != nil {
			return err
		}
		if _, err := w.Write(b); err != nil {
			return err
		}
	}
`)

	case "github.com/ipfs/go-cid.Cid":
		return doTemplate(w, f, `
	if err := cbg.WriteCid(w, t.{{ .Name }}); err != nil {
		return err
	}
`)
	default:
		return doTemplate(w, f, `
{{ if .Pointer }}
	t.{{ .Name }} = new({{ .Type.Name }})
{{ end }}
	if err := t.{{ .Name }}.MarshalCBOR(w); err != nil {
		return err
	}
`)
	}

}

func emitCborMarshalUint64Field(w io.Writer, f Field) error {
	if f.Pointer {
		return fmt.Errorf("pointers to integers not supported")
	}
	return doTemplate(w, f, `
	if _, err := w.Write(cbg.CborEncodeMajorType(cbg.MajUnsignedInt, t.{{ .Name }})); err != nil {
		return err
	}
`)
}

func emitCborMarshalSliceField(w io.Writer, f Field) error {
	if f.Pointer {
		return fmt.Errorf("pointers to slices not supported")
	}
	e := f.Type.Elem()

	if e.Kind() == reflect.Uint8 {
		return doTemplate(w, f, `
	if _, err := w.Write(cbg.CborEncodeMajorType(cbg.MajByteString, uint64(len(t.{{ .Name }})))); err != nil {
		return err
	}
	if _, err := w.Write(t.{{ .Name}}); err != nil {
		return err
	}
`)
	}

	if e.Kind() == reflect.Ptr {
		e = e.Elem()
	}

	switch e.Kind() {
	default:
		return fmt.Errorf("do not yet support slices of non-structs: %s %s", f.Type.Elem(), e.Kind())
	case reflect.Struct:
		// ok
	}

	return doTemplate(w, f, `
	if _, err := w.Write(cbg.CborEncodeMajorType(cbg.MajArray, uint64(len(t.{{ .Name }})))); err != nil {
		return err
	}
	for _, v := range t.{{ .Name }} {
		if err := v.MarshalCBOR(w); err != nil {
			return err
		}
	}
`)
}

func emitCborMarshalStructTuple(w io.Writer, gti *GenTypeInfo) error {
	err := doTemplate(w, gti, `func (t *{{ .Name }}) MarshalCBOR(w io.Writer) error {
	if _, err := w.Write({{ .HeaderAsByteString }}); err != nil {
		return err
	}
`)
	if err != nil {
		return err
	}

	for _, f := range gti.Fields {
		fmt.Fprintf(w, "\n\t// t.%s (%s)", f.Name, f.Type)

		switch f.Type.Kind() {
		case reflect.String:
			if err := emitCborMarshalStringField(w, f); err != nil {
				return err
			}
		case reflect.Struct:
			if err := emitCborMarshalStructField(w, f); err != nil {
				return err
			}
		case reflect.Uint64:
			if err := emitCborMarshalUint64Field(w, f); err != nil {
				return err
			}
		case reflect.Slice:
			if err := emitCborMarshalSliceField(w, f); err != nil {
				return err
			}
		default:
			return fmt.Errorf("field %q of %q has unsupported kind %q", f.Name, gti.Name, f.Type.Kind())
		}
	}

	fmt.Fprintf(w, "\treturn nil\n}\n\n")
	return nil
}

func emitCborUnmarshalStringField(w io.Writer, f Field) error {
	if f.Pointer {
		return fmt.Errorf("pointers to strings not supported")
	}
	return doTemplate(w, f, `
	maj, extra, err = cbg.CborReadHeader(br)
	if err != nil {
		return err
	}

	if maj != cbg.MajTextString {
		return fmt.Errorf("expected cbor type 'text string' in input")
	}

	if extra > 256 * 1024 {
		return fmt.Errorf("string in cbor input too long")
	}

	{
		buf := make([]byte, extra)
		if _, err := io.ReadFull(br, buf); err != nil {
			return err
		}

		t.{{ .Name }} = string(buf)
	}
`)
}

func emitCborUnmarshalStructField(w io.Writer, f Field) error {
	fname := f.Type.PkgPath() + "." + f.Type.Name()
	switch fname {
	case "math/big.Int":
		return doTemplate(w, f, `
	maj, extra, err = cbg.CborReadHeader(br)
	if err != nil {
		return err
	}

	if maj != cbg.MajTag || extra != 2 {
		return fmt.Errorf("big ints should be cbor bignums")
	}

	maj, extra, err = cbg.CborReadHeader(br)
	if err != nil {
		return err
	}

	if maj != cbg.MajByteString {
		return fmt.Errorf("big ints should be tagged cbor byte strings")
	}

	if extra > 256 {
		return fmt.Errorf("cbor bignum was too large")
	}

	if extra > 0 {
		buf := make([]byte, extra)
		if _, err := io.ReadFull(br, buf); err != nil {
			return err
		}
		t.{{ .Name }} = big.NewInt(0).SetBytes(buf)
	}
`)
	case "github.com/ipfs/go-cid.Cid":
		return doTemplate(w, f, `
	{
		c, err := cbg.ReadCid(br)
		if err != nil {
			return err
		}
		t.{{ .Name }} = c
	}
`)
	default:
		return doTemplate(w, f, `
{{ if f.Pointer }}
	t.{{ .Name }} = new({{ .Type.Name }})
{{ end }}
	if err := t.{{ .Name }}.UnmarshalCBOR(br); err != nil {
		return err
	}
`)
	}
}

func emitCborUnmarshalUint64Field(w io.Writer, f Field) error {
	return doTemplate(w, f, `
	maj, extra, err = cbg.CborReadHeader(br)
	if err != nil {
		return err
	}
	if maj != cbg.MajUnsignedInt {
		return fmt.Errorf("wrong type for uint64 field")
	}
	t.{{ .Name }} = extra
`)
}

func emitCborUnmarshalSliceField(w io.Writer, f Field) error {
	e := f.Type.Elem()
	var pointer bool
	if e.Kind() == reflect.Ptr {
		pointer = true
		e = e.Elem()
	}

	err := doTemplate(w, nil, `
	maj, extra, err = cbg.CborReadHeader(br)
	if err != nil {
		return err
	}
	if extra > 8192 {
		return fmt.Errorf("array too large")
	}
`)
	if err != nil {
		return err
	}

	if e.Kind() == reflect.Uint8 {
		return doTemplate(w, f, `
	if maj != cbg.MajByteString {
		return fmt.Errorf("expected byte array")
	}
	t.{{ .Name }} = make([]byte, extra)
	if _, err := io.ReadFull(br, t.{{ .Name }}); err != nil {
		return err
	}
`)
	}

	err = doTemplate(w, f, `
	if maj != cbg.MajArray {
		return fmt.Errorf("expected cbor array")
	}
	if extra > 0 {
		t.{{ .Name }} = make({{ f.Type }}, 0, extra)
	}
	for i := 0; i < int(extra); i++ {
`)
	if err != nil {
		return err
	}

	switch e.Kind() {
	case reflect.Struct:
		fmt.Fprintf(w, "\t\tvar v %s\n", e.Name())
		fmt.Fprintf(w, "\t\tif err := v.UnmarshalCBOR(br); err != nil {\n\t\t\treturn err\n\t\t}\n\n")

		var ptrfix string
		if pointer {
			ptrfix = "&"
		}
		fmt.Fprintf(w, "\t\tt.%s = append(t.%s, %sv)\n", f.Name, f.Name, ptrfix)
	default:
		return fmt.Errorf("do not yet support slices of non-structs: %s", e)
	}
	fmt.Fprintf(w, "\t}\n\n")

	return nil
}

// Generates 'tuple representation' cbor encoders for the given type
func GenTupleEncodersForType(i interface{}, w io.Writer) error {
	gti, err := ParseTypeInfo(i)
	if err != nil {
		return err
	}

	if err := emitCborMarshalStructTuple(w, gti); err != nil {
		return err
	}

	// Now for the unmarshal

	err = doTemplate(w, gti, `
func (t *{{ .Name}}) UnmarshalCBOR(br cbg.ByteReader) error {

	maj, extra, err := cbg.CborReadHeader(br)
	if err != nil {
		return err
	}
	if maj != cbg.MajArray {
		return fmt.Errorf("cbor input should be of type array")
	}

	if extra != {{ len .Fields }} {
		return fmt.Errorf("cbor input had wrong number of fields")
	}

`)
	if err != nil {
		return err
	}

	for _, f := range gti.Fields {
		fmt.Fprintf(w, "\t// t.%s (%s)\n", f.Name, f.Type)
		switch f.Type.Kind() {
		case reflect.String:
			if err := emitCborUnmarshalStringField(w, f); err != nil {
				return err
			}
		case reflect.Struct:

			if err := emitCborUnmarshalStructField(w, f); err != nil {
				return err
			}

		case reflect.Uint64:

			if err := emitCborUnmarshalUint64Field(w, f); err != nil {
				return err
			}
		case reflect.Slice:
			if err := emitCborUnmarshalSliceField(w, f); err != nil {
				return err
			}
		default:
			return fmt.Errorf("field %q of %q has unsupported kind %q", f.Name, gti.Name, f.Type.Kind())
		}
	}

	fmt.Fprintf(w, "\treturn nil\n}\n\n")

	return nil
}
