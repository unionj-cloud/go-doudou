package vo

type age int

type Event struct {
	Name      string
	EventType int
}

type TestAlias struct {
	Age    age
	School []struct {
		Name string
		Addr struct {
			Zip   string
			Block string
			Full  string
		}
	}
}
