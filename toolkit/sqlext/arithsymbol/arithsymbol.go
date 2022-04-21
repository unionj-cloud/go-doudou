package arithsymbol

// ArithSymbol is alias type of string for arithmetic operators
type ArithSymbol string

const (
	// Eq equal
	Eq ArithSymbol = "="
	// Ne not equal
	Ne ArithSymbol = "!="
	// Gt greater than
	Gt ArithSymbol = ">"
	// Lt less than
	Lt ArithSymbol = "<"
	// Gte greater than and equal
	Gte ArithSymbol = ">="
	// Lte less than and equal
	Lte ArithSymbol = "<="
	// Is e.g. is null
	Is ArithSymbol = "is"
	// Not e.g. is not null
	Not ArithSymbol = "is not"
	// In contained by a slice
	In    ArithSymbol = "in"
	NotIn ArithSymbol = "not in"
	Like  ArithSymbol = "like"
)
