package main

type TestBase3 struct {
	index string `json:"index"`
	Type  string `json:"type"`
}

type TestEmbed3 struct {
	TestBase3
	Fields []Field `json:"fields"`
}
