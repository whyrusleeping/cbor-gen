package testing

type SignedArray struct {
	Signed []uint64
}

type SimpleTypeOne struct {
	Foo    string
	Value  uint64
	Binary []byte
	Signed int64
}

type SimpleTypeTwo struct {
	Stuff        *SimpleTypeTwo
	Others       []uint64
	SignedOthers []int64
	Test         [][]byte
	Dog          string
}

type SimpleTypeTree struct {
	Stuff  *SimpleTypeTree
	Stufff *SimpleTypeTwo
	Others []uint64
	Test   [][]byte
	Dog    string
}
