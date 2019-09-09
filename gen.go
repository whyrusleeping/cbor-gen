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

const MaxLength = 8192

const (
	MajUnsignedInt = 0
	MajNegativeInt = 1
	MajByteString  = 2
	MajTextString  = 3
	MajArray       = 4
	MajMap         = 5
	MajTag         = 6
	MajOther       = 7
)

type CBORUnmarshaler interface {
	UnmarshalCBOR(io.Reader) error
}

type CBORMarshaler interface {
	MarshalCBOR(io.Writer) error
}

type Deferred struct {
	Raw []byte
}

func (d *Deferred) MarshalCBOR(w io.Writer) error {
	_, err := w.Write(d.Raw)
	return err
}

func (d *Deferred) UnmarshalCBOR(br io.Reader) error {
	// TODO: theres a more efficient way to implement this method, but for now
	// this is fine
	maj, extra, err := CborReadHeader(br)
	if err != nil {
		return err
	}
	header := CborEncodeMajorType(maj, extra)

	switch maj {
	case MajUnsignedInt, MajNegativeInt, MajOther:
		d.Raw = header
		return nil
	case MajByteString, MajTextString:
		buf := make([]byte, int(extra)+len(header))
		copy(buf, header)
		if _, err := io.ReadFull(br, buf[len(header):]); err != nil {
			return err
		}

		d.Raw = buf

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
	case MajMap:
		d.Raw = header
		sub := new(Deferred)
		for i := 0; i < int(extra*2); i++ {
			sub.Raw = sub.Raw[:0]
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

func readByte(r io.Reader) (byte, error) {
	if br, ok := r.(io.ByteReader); ok {
		return br.ReadByte()
	}
	var b [1]byte
	_, err := io.ReadFull(r, b[:])
	return b[0], err
}

func CborReadHeader(br io.Reader) (byte, uint64, error) {
	first, err := readByte(br)
	if err != nil {
		return 0, 0, err
	}

	maj := (first & 0xe0) >> 5
	low := first & 0x1f

	switch {
	case low < 24:
		return maj, uint64(low), nil
	case low == 24:
		next, err := readByte(br)
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

func ReadTaggedByteArray(br io.Reader, exptag uint64, maxlen uint64) ([]byte, error) {
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

func WriteBool(w io.Writer, b bool) error {
	if b {
		if _, err := w.Write([]byte{0xf5}); err != nil {
			return err
		}
	} else {
		if _, err := w.Write([]byte{0xf4}); err != nil {
			return err
		}
	}
	return nil
}

func ReadString(r io.Reader) (string, error) {
	maj, l, err := CborReadHeader(r)
	if err != nil {
		return "", err
	}

	if maj != MajTextString {
		return "", fmt.Errorf("got tag %d while reading string value", maj)
	}

	if l > MaxLength {
		return "", fmt.Errorf("string in input was too long")
	}

	buf := make([]byte, l)
	_, err = io.ReadFull(r, buf)
	if err != nil {
		return "", err
	}

	return string(buf), nil
}

func ReadCid(br io.Reader) (cid.Cid, error) {
	buf, err := ReadTaggedByteArray(br, 42, 512)
	if err != nil {
		return cid.Undef, err
	}

	if len(buf) == 0 {
		return cid.Undef, fmt.Errorf("undefined cid")
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
		return fmt.Errorf("undefined cid")
		//return CborWriteHeader(w, MajByteString, 0)
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
	xerrors "golang.org/x/xerrors"
)

/* This file was generated by github.com/whyrusleeping/cbor-gen */

var _ = xerrors.Errorf

`)
}

type Field struct {
	Name    string
	Pointer bool
	Type    reflect.Type

	IterLabel string
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
			Name:    "t." + f.Name,
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
	if _, err := w.Write(cbg.CborEncodeMajorType(cbg.MajTextString, uint64(len({{ .Name }})))); err != nil {
		return err
	}
	if _, err := w.Write([]byte({{ .Name }})); err != nil {
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
		if {{ .Name }} != nil {
			b = {{ .Name }}.Bytes()
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
	if err := cbg.WriteCid(w, {{ .Name }}); err != nil {
		return xerrors.Errorf("failed to write cid field {{ .Name }}: %w", err)
	}
`)
	default:
		return doTemplate(w, f, `
	if err := {{ .Name }}.MarshalCBOR(w); err != nil {
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
	if _, err := w.Write(cbg.CborEncodeMajorType(cbg.MajUnsignedInt, {{ .Name }})); err != nil {
		return err
	}
`)
}

func emitCborMarshalBoolField(w io.Writer, f Field) error {
	return doTemplate(w, f, `
	if err := cbg.WriteBool(w, {{ .Name }}); err != nil {
		return err
	}
`)
}

func emitCborMarshalMapField(w io.Writer, f Field) error {
	err := doTemplate(w, f, `
	if err := cbg.CborWriteHeader(w, cbg.MajMap, uint64(len({{ .Name }}))); err != nil {
		return err
	}

	for k, v := range {{ .Name }} {
`)
	if err != nil {
		return err
	}

	// Map key
	switch f.Type.Key().Kind() {
	case reflect.String:
		if err := emitCborMarshalStringField(w, Field{Name: "k"}); err != nil {
			return err
		}
	default:
		return fmt.Errorf("non-string map keys are not yet supported")
	}

	// Map value
	switch f.Type.Elem().Kind() {
	case reflect.Ptr:
		if f.Type.Elem().Elem().Kind() != reflect.Struct {
			return fmt.Errorf("unsupported map elem ptr type: %s", f.Type.Elem())
		}

		fallthrough
	case reflect.Struct:
		if err := emitCborMarshalStructField(w, Field{Name: "v", Type: f.Type.Elem()}); err != nil {
			return err
		}
	default:
		return fmt.Errorf("currently unsupported map elem type: %s", f.Type.Elem())
	}

	return doTemplate(w, f, `
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
	if _, err := w.Write(cbg.CborEncodeMajorType(cbg.MajByteString, uint64(len({{ .Name }})))); err != nil {
		return err
	}
	if _, err := w.Write({{ .Name }}); err != nil {
		return err
	}
`)
	}

	if e.Kind() == reflect.Ptr {
		e = e.Elem()
	}

	err := doTemplate(w, f, `
	if _, err := w.Write(cbg.CborEncodeMajorType(cbg.MajArray, uint64(len({{ .Name }})))); err != nil {
		return err
	}
	for _, v := range {{ .Name }} {`)
	if err != nil {
		return err
	}

	switch e.Kind() {
	default:
		return fmt.Errorf("do not yet support slices of %s yet", e.Kind())
	case reflect.Struct:
		fname := e.PkgPath() + "." + e.Name()
		switch fname {
		case "github.com/ipfs/go-cid.Cid":
			err := doTemplate(w, f, `
		if err := cbg.WriteCid(w, v); err != nil {
			return xerrors.Errorf("failed writing cid field {{ .Name }}: %w", err)
		}
`)
			if err != nil {
				return err
			}

		default:
			err := doTemplate(w, f, `
		if err := v.MarshalCBOR(w); err != nil {
			return err
		}
`)
			if err != nil {
				return err
			}
		}
	case reflect.Uint64:
		err := doTemplate(w, f, `
		if err := cbg.CborWriteHeader(w, cbg.MajUnsignedInt, v); err != nil {
			return err
		}
`)
		if err != nil {
			return err
		}
	case reflect.Slice:
		subf := Field{Name: "v", Type: e}
		if err := emitCborMarshalSliceField(w, subf); err != nil {
			return err
		}
	}

	// array end
	fmt.Fprintf(w, "\t}\n")
	return nil
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
		case reflect.Bool:
			if err := emitCborMarshalBoolField(w, f); err != nil {
				return err
			}
		case reflect.Map:
			if err := emitCborMarshalMapField(w, f); err != nil {
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
	{{ .Name }}, err = cbg.ReadString(br)
	if err != nil {
		return err
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
		{{ .Name }} = big.NewInt(0).SetBytes(buf)
	}
`)
	case "github.com/ipfs/go-cid.Cid":
		return doTemplate(w, f, `
	{
		c, err := cbg.ReadCid(br)
		if err != nil {
			return xerrors.Errorf("failed to read cid field {{ .Name }}: %w", err)
		}
		{{ .Name }} = c
	}
`)
	default:
		return doTemplate(w, f, `
{{ if .Pointer }}
	{{ .Name }} = new({{ .Type }})
{{ end }}
	if err := {{ .Name }}.UnmarshalCBOR(br); err != nil {
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
	{{ .Name }} = extra
`)
}

func emitCborUnmarshalBoolField(w io.Writer, f Field) error {
	return doTemplate(w, f, `
	maj, extra, err = cbg.CborReadHeader(br)
	if err != nil {
		return err
	}
	if maj != cbg.MajOther {
		return fmt.Errorf("booleans must be major type 7")
	}
	switch extra {
	case 20:
		{{ .Name }} = false
	case 21:
		{{ .Name }} = true
	default:
		return fmt.Errorf("booleans are either major type 7, value 20 or 21 (got %d)", extra)
	}
`)
}

func emitCborUnmarshalMapField(w io.Writer, f Field) error {
	err := doTemplate(w, f, `
	maj, extra, err = cbg.CborReadHeader(br)
	if err != nil {
		return err
	}
	if maj != cbg.MajMap {
		return fmt.Errorf("expected a map (major type 5)")
	}
	if extra > 4096 {
		return fmt.Errorf("map too large")
	}

	{{ .Name }} = make({{ .Type }}, extra)


	for i, l := 0, int(extra); i < l; i++ {
`)
	if err != nil {
		return err
	}

	switch f.Type.Key().Kind() {
	case reflect.String:
		if err := doTemplate(w, f, `
	var k string
`); err != nil {
			return err
		}
		if err := emitCborUnmarshalStringField(w, Field{Name: "k"}); err != nil {
			return err
		}
	default:
		return fmt.Errorf("maps with non-string keys are not yet supported")
	}

	var pointer bool
	t := f.Type.Elem()
	switch t.Kind() {
	case reflect.Ptr:
		if t.Elem().Kind() != reflect.Struct {
			return fmt.Errorf("unsupported map elem ptr type: %s", t)
		}

		t = t.Elem()
		pointer = true
		fallthrough
	case reflect.Struct:
		if err := doTemplate(w, f, `
	var v {{ .Type.Elem }}
`); err != nil {
			return err
		}
		if err := emitCborUnmarshalStructField(w, Field{Name: "v", Pointer: pointer, Type: t}); err != nil {
			return err
		}
		if err := doTemplate(w, f, `
	{{ .Name }}[k] = v
`); err != nil {
			return err
		}
	default:
		return fmt.Errorf("currently only support maps of structs")
	}

	return doTemplate(w, f, `
	}
`)
}

func emitCborUnmarshalSliceField(w io.Writer, f Field) error {
	if f.IterLabel == "" {
		f.IterLabel = "i"
	}

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
	{{ .Name }} = make([]byte, extra)
	if _, err := io.ReadFull(br, {{ .Name }}); err != nil {
		return err
	}
`)
	}

	err = doTemplate(w, f, `
	if maj != cbg.MajArray {
		return fmt.Errorf("expected cbor array")
	}
	if extra > 0 {
		{{ .Name }} = make({{ .Type }}, extra)
	}
	for {{ .IterLabel }} := 0; {{ .IterLabel }} < int(extra); {{ .IterLabel}}++ {
`)
	if err != nil {
		return err
	}

	switch e.Kind() {
	case reflect.Struct:
		fname := e.PkgPath() + "." + e.Name()
		switch fname {
		case "github.com/ipfs/go-cid.Cid":
			err := doTemplate(w, f, `
		c, err := cbg.ReadCid(br)
		if err != nil {
			return xerrors.Errorf("reading cid field {{ .Name }} failed: %w", err)
		}
		{{ .Name }}[{{ .IterLabel }}] = c
`)
			if err != nil {
				return err
			}
		default:
			fmt.Fprintf(w, "\t\tvar v %s\n", e.Name())
			fmt.Fprintf(w, "\t\tif err := v.UnmarshalCBOR(br); err != nil {\n\t\t\treturn err\n\t\t}\n\n")

			var ptrfix string
			if pointer {
				ptrfix = "&"
			}
			fmt.Fprintf(w, "\t\t%s[%s] = %sv\n", f.Name, f.IterLabel, ptrfix)
		}
	case reflect.Uint64:
		err := doTemplate(w, f, `
		maj, val, err := cbg.CborReadHeader(br)
		if err != nil {
			return xerrors.Errorf("failed to read uint64 for {{ .Name }} slice: %w", err)
		}

		if maj != cbg.MajUnsignedInt {
			return xerrors.Errorf("value read for array {{ .Name }} was not a uint, instead got %d", maj)
		}
		
		{{ .Name }}[{{ .IterLabel}}] = val
`)
		if err != nil {
			return err
		}
	case reflect.Slice:
		nextIter := string([]byte{f.IterLabel[0] + 1})
		subf := Field{
			Name:      fmt.Sprintf("%s[%s]", f.Name, f.IterLabel),
			Type:      e,
			IterLabel: nextIter,
		}
		fmt.Fprintf(w, "\t\t{\n\t\t\tvar maj byte\n\t\tvar extra uint64\n\t\tvar err error\n")
		if err := emitCborUnmarshalSliceField(w, subf); err != nil {
			return err
		}
		fmt.Fprintf(w, "\t\t}\n")

	default:
		return fmt.Errorf("do not yet support slices of %s yet", e.Elem())
	}
	fmt.Fprintf(w, "\t}\n\n")

	return nil
}

func emitCborUnmarshalStructTuple(w io.Writer, gti *GenTypeInfo) error {
	err := doTemplate(w, gti, `
func (t *{{ .Name}}) UnmarshalCBOR(br io.Reader) error {

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
		case reflect.Bool:
			if err := emitCborUnmarshalBoolField(w, f); err != nil {
				return err
			}
		case reflect.Map:
			if err := emitCborUnmarshalMapField(w, f); err != nil {
				return err
			}
		default:
			return fmt.Errorf("field %q of %q has unsupported kind %q", f.Name, gti.Name, f.Type.Kind())
		}
	}

	fmt.Fprintf(w, "\treturn nil\n}\n\n")

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

	if err := emitCborUnmarshalStructTuple(w, gti); err != nil {
		return err
	}

	return nil
}
