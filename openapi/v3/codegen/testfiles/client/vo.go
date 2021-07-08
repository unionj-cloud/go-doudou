package client

import (
	"os"
	"time"
)

type Address struct {
	City string `json:"city,omitempty"`

	State string `json:"state,omitempty"`

	Street string `json:"street,omitempty"`

	Zip string `json:"zip,omitempty"`
}

type ApiResponse struct {
	Code int `json:"code,omitempty"`

	Message string `json:"message,omitempty"`

	Type string `json:"type,omitempty"`
}

type Category struct {
	Id int64 `json:"id,omitempty"`

	Name string `json:"name,omitempty"`
}

type Customer struct {
	Address []Address `json:"address,omitempty"`

	Id int64 `json:"id,omitempty"`

	Username string `json:"username,omitempty"`
}

type Order struct {
	Complete bool `json:"complete,omitempty"`
	// 客户信息结构体
	// 用于描述客户相关的信息
	Customer struct {
		// 用户ID
		Id int64 `json:"id,omitempty"`
		// 用户名
		Username string `json:"username,omitempty"`
		// 用户地址
		// 例如：北京海淀区xxx街道
		// 某某小区
		Address []Address `json:"address,omitempty"`
	} `json:"customer,omitempty"`

	Id int64 `json:"id,omitempty"`

	PetId int64 `json:"petId,omitempty"`

	Quantity int `json:"quantity,omitempty"`

	ShipDate *time.Time `json:"shipDate,omitempty"`
	// Order Status
	Status string `json:"status,omitempty"`
}

type Pet struct {
	Category Category `json:"category,omitempty"`

	Id int64 `json:"id,omitempty"`

	// required
	Name string `json:"name,omitempty"`

	// required
	PhotoUrls []string `json:"photoUrls,omitempty"`
	// pet status in the store
	// this is another line for test use
	Status string `json:"status,omitempty"`

	Tags []Tag `json:"tags,omitempty"`
}

type Tag struct {
	Id int64 `json:"id,omitempty"`

	Name string `json:"name,omitempty"`
}

type User struct {
	Additional1 map[string]string `json:"additional1,omitempty"`

	Additional2 map[string]Tag `json:"additional2,omitempty"`

	Avatar *os.File `json:"avatar,omitempty"`

	Email string `json:"email,omitempty"`

	FirstName string `json:"firstName,omitempty"`

	Id int64 `json:"id,omitempty"`

	LastName string `json:"lastName,omitempty"`

	Password string `json:"password,omitempty"`

	Phone string `json:"phone,omitempty"`
	// User Status
	UserStatus int `json:"userStatus,omitempty"`

	Username string `json:"username,omitempty"`
}
