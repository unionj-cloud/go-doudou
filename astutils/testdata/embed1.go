package main

type TestBase1 struct {
	Index string
	Type  string
}

type TestEmbed1 struct {
	TestBase1 `json:"test_base_1,omitempty"`
	Fields    []Field `json:"fields"`
}
