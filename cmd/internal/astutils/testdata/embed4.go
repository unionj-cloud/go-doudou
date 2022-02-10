package main

type TestBase4 struct {
	index string
	Type  string
}

type TestEmbed4 struct {
	TestBase1 `json:"test_base_1,omitempty"`
	Fields    []Field `json:"fields"`
}
