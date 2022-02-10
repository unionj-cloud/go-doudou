package extraenum

// Extra part of defining column
type Extra string

const (
	// Update used for update_at column
	Update Extra = "on update CURRENT_TIMESTAMP"
)
