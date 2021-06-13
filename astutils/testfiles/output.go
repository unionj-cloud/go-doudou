package main

import (
	"encoding/json"
	"fmt"
)

type UserVo struct {
	Id    int
	Name  string
	Phone string
	Dept  string
}

type Page struct {
	PageNo int
	Size   int
	Items  []UserVo
}

func main() {
	page := Page{
		PageNo: 10,
		Size:   30,
	}
	b, _ := json.Marshal(page)
	fmt.Println(string(b))
}
