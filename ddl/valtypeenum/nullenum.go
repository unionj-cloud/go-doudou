package valtypeenum

// ValType represents column value
type ValType int

const (
	// Func represents mysql built-in function
	Func ValType = iota
	// Null represents null value
	Null
	// Literal represents literal value
	Literal
)
