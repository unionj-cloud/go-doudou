package main

type testBase struct {
	Index string `json:"index"`
	Type  string `json:"type"`
}

type TestEmbed2 struct {
	testBase
	Fields []Field `json:"fields"`
}
