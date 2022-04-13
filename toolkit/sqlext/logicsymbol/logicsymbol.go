package logicsymbol

// LogicSymbol define logic operators
type LogicSymbol string

const (
	// And logic
	And LogicSymbol = "and"
	// Or logic
	Or LogicSymbol = "or"
	// Append logic
	Append LogicSymbol = " "
	// End is similar with Append except that End does not add parentheses
	End LogicSymbol = ""
)
