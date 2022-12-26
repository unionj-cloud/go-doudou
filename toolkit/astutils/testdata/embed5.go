package main

type testBase5 struct {
	Index string `json:"index"`
	Type  string `json:"type"`
}

type TestEmbed5 struct {
	testBase5 `json:"testBase"`
	Fields    []Field `json:"fields"`
}
