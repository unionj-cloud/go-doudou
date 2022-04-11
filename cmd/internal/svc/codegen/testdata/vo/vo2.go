package vo

type TestExprStringP struct {
	Age     age
	Hobbies [3]string
	Data    map[string]string
	School  []struct {
		Name string
		Addr struct {
			Zip   string
			Block string
			Full  string
		}
	}
}
