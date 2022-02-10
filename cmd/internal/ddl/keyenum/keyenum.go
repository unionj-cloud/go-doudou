package keyenum

// Key define index type
type Key string

const (
	// Pri primary key
	Pri Key = "PRI"
	// Uni unique index
	Uni Key = "UNI"
	// Mul composite index
	Mul Key = "MUL"
	// Empty blank
	Empty Key = ""
)
