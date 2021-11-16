package test

import (
	"os"
	"time"
)

type Address struct {
	City string `json:"city" url:"city"`

	State string `json:"state" url:"state"`

	Street string `json:"street" url:"street"`

	Zip string `json:"zip" url:"zip"`
}

type ApiResponse struct {
	Code int `json:"code" url:"code"`

	Message string `json:"message" url:"message"`

	Type string `json:"type" url:"type"`
}

type Category struct {
	Id int64 `json:"id" url:"id"`

	Name string `json:"name" url:"name"`
}

type Customer struct {
	Address []Address `json:"address" url:"address"`

	Id int64 `json:"id" url:"id"`

	Username string `json:"username" url:"username"`
}

type Order struct {
	Complete bool `json:"complete" url:"complete"`
	// 客户信息结构体
	// 用于描述客户相关的信息
	Customer struct {
		// 用户ID
		Id int64 `json:"id" url:"id"`
		// 用户名
		Username string `json:"username" url:"username"`
		// 用户地址
		// 例如：北京海淀区xxx街道
		// 某某小区
		Address []Address `json:"address" url:"address"`
	} `json:"customer" url:"customer"`

	Id int64 `json:"id" url:"id"`

	PetId int64 `json:"petId" url:"petId"`

	Quantity int `json:"quantity" url:"quantity"`

	ShipDate *time.Time `json:"shipDate" url:"shipDate"`
	// Order Status
	Status string `json:"status" url:"status"`
}

type Pet struct {
	Category Category `json:"category" url:"category"`

	Id int64 `json:"id" url:"id"`

	// required
	Name string `json:"name" url:"name"`

	// required
	PhotoUrls []string `json:"photoUrls" url:"photoUrls"`
	// pet status in the store
	// this is another line for test use
	Status string `json:"status" url:"status"`

	Tags []Tag `json:"tags" url:"tags"`
}

type Tag struct {
	Id int64 `json:"id" url:"id"`

	Name string `json:"name" url:"name"`
}

type User struct {
	Additional1 struct {
	} `json:"additional1" url:"additional1"`

	Additional2 struct {
	} `json:"additional2" url:"additional2"`

	Avatar *os.File `json:"avatar" url:"avatar"`

	Email string `json:"email" url:"email"`

	FirstName string `json:"firstName" url:"firstName"`

	Id int64 `json:"id" url:"id"`

	LastName string `json:"lastName" url:"lastName"`

	Password string `json:"password" url:"password"`

	Phone string `json:"phone" url:"phone"`
	// User Status
	UserStatus int `json:"userStatus" url:"userStatus"`

	Username string `json:"username" url:"username"`
}
