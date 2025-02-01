package typegen

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

// TypeExpr is the root of the AST. It represents a Go type expression that can be:
//   - A named type (with an optional package alias/path).
//   - A pointer to another TypeExpr.
//   - A slice or array of another TypeExpr.
//   - A map of (KeyExpr -> ValueExpr).
//   - A generic instantiation: BaseTypeExpr[T1, T2, ...]  (the base can itself be a named type).
type TypeExpr struct {
	// If non-empty, we have a pointer depth.
	// e.g. **[]Foo => PtrDepth=2, then an underlying that is "[]Foo"
	PtrDepth int

	// If ArrayLen >= 0, that means it's an array (fixed length). -1 means not an array.
	ArrayLen int

	// If IsSlice is true, we have a slice. (We do not store the length for slices.)
	IsSlice bool

	// If IsMap is true, then MapKey/MapValue are valid.
	IsMap    bool
	MapKey   *TypeExpr
	MapValue *TypeExpr

	// Base package path (or alias) and name, if this is a named type.
	// E.g. PkgPath="github.com/foo/bar", Name="MyType".
	// If the underlying is something else (map, slice, etc.), these may be "".
	PkgPath string
	Name    string // base type name, e.g. "CborLink"

	// If we have generic parameters, they are stored here. For example:
	//   "CborLink[*SimpleTypeTwo, SomeOther]" => Generics has 2 entries.
	// If we do not have generics, Generics is nil or empty.
	Generics []*TypeExpr

	// If we have a parent in the parse, it can be stored for convenience, but not required.
}

// parseTypeName parses a type expression (like reflect.Type.String()) into a *TypeExpr.
func parseTypeName(full string) (*TypeExpr, error) {
	// We’ll produce tokens, then parse them.
	tokens, err := tokenizeTypeString(full)
	if err != nil {
		return nil, fmt.Errorf("tokenize failed: %w", err)
	}
	p := &typeParser{tokens: tokens}
	expr, err := p.parseTypeExpr()
	if err != nil {
		return nil, fmt.Errorf("parseTypeExpr failed near token %d (\"%s\"): %w",
			p.pos, p.currToken(), err)
	}
	// We expect no trailing tokens.
	if p.pos < len(tokens) {
		return nil, fmt.Errorf("extra tokens after parse: %v", tokens[p.pos:])
	}
	return expr, nil
}

// tokenizeTypeString splits things like "map[string]*github.com/foo/bar.Baz" into
// tokens: "map", "[", "string", "]", "*", "github.com/foo/bar.Baz" etc.
// We keep bracket depth logic minimal. We'll handle identifiers, punctuation, brackets, etc.
func tokenizeTypeString(s string) ([]string, error) {
	var tokens []string
	i := 0
	for i < len(s) {
		c := s[i]
		switch c {
		case '[', ']', ',', '*', '(', ')':
			tokens = append(tokens, string(c))
			i++
		case ' ':
			// skip
			i++
		case '{', '}':
			// If reflect ever gives curly braces (uncommon), we’d parse them.
			// For now, let’s just treat them as punctuation tokens if needed.
			tokens = append(tokens, string(c))
			i++
		default:
			if c == ':' || c == '/' || c == '.' {
				// these are part of a "word" in import paths or type names
				// but let's parse them in the same token as the preceding alphanumerics
				// so we consume them in one token like "github.com/foo/bar.Baz"
				// We'll parse until next punctuation or bracket
				start := i
				for i < len(s) && !isPunctuation(s[i]) {
					i++
				}
				tokens = append(tokens, s[start:i])
			} else if isLetterDigit(c) {
				// parse an identifier (plus any embedded '.' or version segments, etc.)
				start := i
				for i < len(s) && !isPunctuation(s[i]) && s[i] != ' ' {
					i++
				}
				tokens = append(tokens, s[start:i])
			} else {
				return nil, fmt.Errorf("unexpected character '%c' at %d in %q", c, i, s)
			}
		}
	}
	return tokens, nil
}

func isPunctuation(c byte) bool {
	switch c {
	case '[', ']', ',', '*', '(', ')', '{', '}', ' ':
		return true
	}
	return false
}
func isLetterDigit(c byte) bool {
	return c == '_' || c == '/' || c == '.' || unicode.IsLetter(rune(c)) || unicode.IsDigit(rune(c)) || c == ':'
}

// typeParser is a simple index into the token list. We do a single pass parse.
type typeParser struct {
	tokens []string
	pos    int
}

func (p *typeParser) done() bool {
	return p.pos >= len(p.tokens)
}

func (p *typeParser) currToken() string {
	if p.done() {
		return ""
	}
	return p.tokens[p.pos]
}

func (p *typeParser) consume(tok string) bool {
	if p.currToken() == tok {
		p.pos++
		return true
	}
	return false
}

func (p *typeParser) expect(tok string) error {
	if p.consume(tok) {
		return nil
	}
	return fmt.Errorf("expected %q, got %q", tok, p.currToken())
}

//
// The grammar we parse (simplified):
//
//  TypeExpr := (Pointer)* (ArrayOrSlice | MapType | NamedType) (GenericParams?) ?
//
//  (Pointer)* means zero or more "*" tokens
//
//  ArrayOrSlice := "[" (number?) "]" TypeExpr
//  MapType := "map" "[" TypeExpr "]" TypeExpr
//  NamedType := ident ( '.' ident )*   (like "github.com/foo.bar.Baz")
//
//  GenericParams := "[" TypeExpr ( "," TypeExpr )* "]"
//
// We parse top-down, with a little loop for leading pointer operators `*` and
// leading `[]`/`[N]`, or leading `map[...]`.
//

// parseTypeExpr is the top-level entry
func (p *typeParser) parseTypeExpr() (*TypeExpr, error) {
	expr := &TypeExpr{
		ArrayLen: -1, // -1 => not an array
	}

	// 1) gather leading pointers
	for p.consume("*") {
		expr.PtrDepth++
	}

	// 2) check if next is "map"
	if p.consume("map") {
		// parse: "[" KeyType "]" ValueType
		expr.IsMap = true
		if err := p.expect("["); err != nil {
			return nil, err
		}
		keyExpr, err := p.parseTypeExpr()
		if err != nil {
			return nil, err
		}
		expr.MapKey = keyExpr
		if err := p.expect("]"); err != nil {
			return nil, err
		}
		valExpr, err := p.parseTypeExpr()
		if err != nil {
			return nil, err
		}
		expr.MapValue = valExpr
		return expr, nil
	}

	// 3) check if next is "[" => slice or array
	if p.consume("[") {
		// Could be "[N]" or "[]"
		tk := p.currToken()
		arrLen, err := strconv.Atoi(tk)
		if err == nil {
			// we have an array
			expr.ArrayLen = arrLen
			p.pos++
		} else {
			// otherwise a slice
			expr.IsSlice = true
		}
		if err := p.expect("]"); err != nil {
			return nil, err
		}
		subExpr, err := p.parseTypeExpr()
		if err != nil {
			return nil, err
		}
		// subExpr is the element type
		// merge pointer depth, etc. Actually we just want to store subExpr as the child of expr
		// But we stored pointers, map, etc. in this same struct. We can nest the expr:
		// Instead, let's nest. We do: "expr" is a slice/array, so store subExpr inside it.
		// We'll do that by reusing the subExpr structure. But we have to combine them carefully.
		expr.mergeUnder(subExpr)
		return expr, nil
	}

	// 4) otherwise we must parse a named base. e.g. "github.com/whyrusleeping/cbor-gen/testing.CborLink"
	// We'll consume tokens until we see something that indicates the next stage.
	baseNameTokens := []string{}
	for !p.done() {
		tk := p.currToken()
		// stop if bracket or punctuation that indicates generics or something
		if tk == "[" || tk == "]" || tk == "," || tk == "*" || tk == "map" || tk == "[" {
			break
		}
		// else it's part of the base name
		baseNameTokens = append(baseNameTokens, tk)
		p.pos++
	}
	if len(baseNameTokens) == 0 {
		return expr, nil // or error? Means something is off
	}
	baseName := strings.Join(baseNameTokens, "")

	// Now store baseName in expr, splitting any last '.' from the type name if needed
	// or if it’s a “slash/dot” for pkg path. We don’t know exactly the boundary for pkg vs type,
	// but the standard reflect.Type.String() typically ends with ".Foo" or for local type "Foo".
	// A quick approach is to find the *last* '.' that is followed by a proper Go identifier.
	idx := strings.LastIndex(baseName, ".")
	if idx < 0 {
		expr.PkgPath = "" // local or built-in
		expr.Name = baseName
	} else {
		// pkg part is everything up to the last dot
		expr.PkgPath = baseName[:idx]
		expr.Name = baseName[idx+1:]
	}

	// 5) If next token is "[", then we have generics
	if p.consume("[") {
		// parse params until "]"
		// e.g. "T1, T2, ..."
		for {
			subT, err := p.parseTypeExpr()
			if err != nil {
				return nil, err
			}
			expr.Generics = append(expr.Generics, subT)
			if p.consume("]") {
				// done
				break
			}
			if !p.consume(",") {
				return nil, fmt.Errorf("expected ',' or ']' in generic params, got %q", p.currToken())
			}
		}
	}

	return expr, nil
}

// mergeUnder merges 'under' type expression into the current one, effectively
// meaning "expr" is an array/slice, pointer, or map that was partially discovered.
func (te *TypeExpr) mergeUnder(sub *TypeExpr) {
	// Our type is "already known" to be a slice/array/pointer. We want to set the
	// underlying type to `sub`. But we must preserve sub’s pointer/slice/etc. as well.
	// The easiest is just store sub as a “child” field, but our struct is “flattened.”
	// We can keep the top-level isSlice, arrayLen, ptrDepth, etc. and nest sub in a single child.
	// But for simpler code generation, let’s do it this way:
	// - The top’s .PtrDepth, .IsSlice, .ArrayLen are already set
	// - The sub might also have pointer depth, generics, etc.
	// So we consider the top’s type to be “prefix” on sub. Then we produce a new top that is sub.
	// But we want “ptrDepth + sub.PtrDepth,” etc.
	//
	// Actually simpler: set te’s “Named” fields from sub’s named fields, set te.Generics from sub, etc.
	// plus we combine pointer depth. In effect, the array (and pointer) is the outer type, the sub is the next layer.

	// We combine pointer depth, map, slice, array, generics, etc.
	// This is one way:
	te.PkgPath = sub.PkgPath
	te.Name = sub.Name
	te.Generics = sub.Generics
	te.IsMap = sub.IsMap
	te.MapKey = sub.MapKey
	te.MapValue = sub.MapValue

	// The sub may also have array, slice, pointer, so we add them.
	te.PtrDepth += sub.PtrDepth

	// If sub is an array, sub.ArrayLen>=0 => we have to reflect that the final type is “array of sub’s deeper type”.
	// But we can keep the top as array/slice, so te.ArrayLen is the “outer,” we want to nest sub’s array deeper.
	// This gets complicated if sub is also an array. That yields multi-dimensional arrays, e.g. `[3][4]T`.
	// So for truly robust: we keep going until we find the base.
	// For simplicity here, we assume the user rarely does `[N][M]T` in generics.
	// If you want to handle multi-dim arrays properly, you’d do a loop.

	// If sub is a slice:
	if sub.IsSlice {
		// we are “outer slice/array” of sub’s element => multi-dimensional slices.
		// For simplicity, let’s just do a chain approach or an error.
		// Real code should do a deeper nested structure.
		te.IsSlice = true
	} else if sub.ArrayLen >= 0 {
		te.ArrayLen = sub.ArrayLen
	}
}

// RenderString renders a TypeExpr back to valid Go syntax. The “shortener”
// that uses your `resolvePkgName(...)` can be a separate pass, or you can do it inline.
func (te *TypeExpr) RenderString() string {
	// pointer(s)
	ptr := strings.Repeat("*", te.PtrDepth)

	if te.IsMap {
		return fmt.Sprintf("%smap[%s]%s",
			ptr,
			te.MapKey.RenderString(),
			te.MapValue.RenderString(),
		)
	}
	if te.IsSlice {
		return fmt.Sprintf("%s[]%s",
			ptr,
			// Render sub-part. But we only stored it in the same struct.
			// Actually if we have multiple slices or pointer layers, we’re flattening them.
			// So do we just rely on the next fields? That might lose multi-slice dimension
			// unless we are storing it deeper. The above simplified approach merges them.
			// So effectively we do: `ptr + "[]" + (the rest)`.
			te.renderBaseAndGenerics(),
		)
	}
	if te.ArrayLen >= 0 {
		return fmt.Sprintf("%s[%d]%s",
			ptr,
			te.ArrayLen,
			te.renderBaseAndGenerics(),
		)
	}
	return ptr + te.renderBaseAndGenerics()
}

func (te *TypeExpr) renderBaseAndGenerics() string {
	// If no name => could be an error or built in. Usually we do e.g. "int" if te.Name="int", te.PkgPath=""
	base := te.Name
	if te.PkgPath != "" {
		base = te.PkgPath + "." + base
	}
	if len(te.Generics) == 0 {
		return base
	}
	// Render generics as e.g. "Base[T1, T2]"
	var sub []string
	for _, g := range te.Generics {
		sub = append(sub, g.RenderString())
	}
	return fmt.Sprintf("%s[%s]", base, strings.Join(sub, ", "))
}

// ShortenPackages traverses the AST and, for any named type that comes
// from a package that is not the current package, rewrites its PkgPath using resolver.
// If the PkgPath equals the current package, it is cleared so that no package prefix is rendered.
func (te *TypeExpr) ShortenPackages(currentPkg string, resolver func(pkgPath, typeName string) string) {
	if te.PkgPath != "" {
		// Extract the last segment from currentPkg.
		lastSlash := strings.LastIndex(currentPkg, "/")
		var currentLast string
		if lastSlash >= 0 {
			currentLast = currentPkg[lastSlash+1:]
		} else {
			currentLast = currentPkg
		}
		// If the package matches the current package or the last segment thereof,
		// clear the package so no prefix is rendered.
		if te.PkgPath == currentPkg || te.PkgPath == currentLast {
			te.PkgPath = ""
		} else if strings.Contains(te.PkgPath, "/") {
			alias := resolver(te.PkgPath, te.Name)
			te.PkgPath = alias
		}
	}

	// Recurse into generic parameters.
	for _, g := range te.Generics {
		g.ShortenPackages(currentPkg, resolver)
	}
	// Recurse into map key/value if applicable.
	if te.IsMap {
		te.MapKey.ShortenPackages(currentPkg, resolver)
		te.MapValue.ShortenPackages(currentPkg, resolver)
	}
}
