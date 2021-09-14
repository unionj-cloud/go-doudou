package extraenum

// Extra part of defining column
type Extra string

const (
	Update Extra = "on update CURRENT_TIMESTAMP"
)
