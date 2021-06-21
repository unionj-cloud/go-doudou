package main

type base struct {
	Index string
	Type  string
}

type struct1 struct {
	base
	Name       string `json:"good"`
	StructType int    `json:"struct_type" dd:"awesomtag"`
	Format     string `dd:"anothertag"`
	Pos        int
}
