package valtypeenum

type ValType int

const (
	Func ValType = iota
	Null
	Literal
)
