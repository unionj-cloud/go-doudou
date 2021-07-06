package client

import (
	"os"
	"time"
)

//go:generate go-doudou name --file $GOFILE

type Address struct {
	City string

	State string

	Street string

	Zip string
}

type ApiResponse struct {
	Code int

	Message string

	Type string
}

type Category struct {
	Id int64

	Name string
}

type Customer struct {
	Address []Address

	Id int64

	Username string
}

type Order struct {
	Complete bool
	// 客户信息结构体
	// 用于描述客户相关的信息
	Customer struct {
		// 用户ID
		Id int64
		// 用户名
		Username string
		// 用户地址
		// 例如：北京海淀区xxx街道
		// 某某小区
		Address []Address
	}

	Id int64

	PetId int64

	Quantity int

	ShipDate *time.Time
	// Order Status
	Status string
}

type Pet struct {
	Category Category

	Id int64

	Name string

	PhotoUrls []string
	// pet status in the store
	// this is another line for test use
	Status string

	Tags []Tag
}

type Tag struct {
	Id int64

	Name string
}

type User struct {
	Additional1 map[string]string

	Additional2 map[string]Tag

	Avatar *os.File

	Email string

	FirstName string

	Id int64

	LastName string

	Password string

	Phone string
	// User Status
	UserStatus int

	Username string
}
