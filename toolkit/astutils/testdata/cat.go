package main

import (
	"context"
	"fmt"
)

type User struct {
	Name string `json:"name"`
}

type Cat struct {
	Hobbies map[string]interface{}
	Sleep   func() bool
	Run     chan string
}

// eat execute eat behavior for Cat
func (c *Cat) eat(food string) (
	// not hungry
	full bool,
	// how feel
	mood string) {
	fmt.Println("eat " + food)
	return true, "happy"
}

func (c *Cat) PostSelectVersionPage(ctx context.Context, body PageDTO[VersionDTO, Cat, User]) (code int, message string, data PageDTO[VersionDTO, Cat, User], err error) {
	return 0, "", PageDTO[VersionDTO, Cat, User]{}, err
}

type VersionDTO struct {
	VersionName string      `json:"versionName"`
	VersionId   interface{} `json:"versionId"`
}

type PageDTO[T any, R any, K any] struct {
	TotalRow   int    `json:"totalRow"`
	PageNumber int    `json:"pageNumber"`
	TotalPage  int    `json:"totalPage"`
	PageSize   int    `json:"pageSize"`
	ReturnMsg  string `json:"return_msg"`
	List       []T    `json:"list"`
	Item       []R    `json:"item"`
	Many       []K    `json:"many"`
	ReturnCode string `json:"return_code"`
}
