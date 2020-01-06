package testing

type SimpleTypeOne struct {
	Foo    string
	Value  uint64
	Binary []byte
}

type SimpleTypeTwo struct {
	Stuff  *SimpleTypeTwo
	Others []uint64
	Test   [][]byte
	Dog    string
}

type SimpleTypeTree struct {
	Stuff  *SimpleTypeTree
	Stufff *SimpleTypeTwo
	Others []uint64
	Test   [][]byte
	Dog    string
}
