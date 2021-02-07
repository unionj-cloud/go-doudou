package vo

//go:generate ns -file $GOFILE
type Student struct {
	Name      string
	Age       int
	TestScore int
}