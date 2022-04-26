package test

import (
	"time"

	v3 "github.com/unionj-cloud/go-doudou/toolkit/openapi/v3"
)

type Address struct {
	City *string `json:"city,omitempty" url:"city"`

	State *string `json:"state,omitempty" url:"state"`

	Street *string `json:"street,omitempty" url:"street"`

	Zip *string `json:"zip,omitempty" url:"zip"`
}

type ApiResponse struct {
	Code *int `json:"code,omitempty" url:"code"`

	Message *string `json:"message,omitempty" url:"message"`

	Type *string `json:"type,omitempty" url:"type"`
}

type Category struct {
	Id *int64 `json:"id,omitempty" url:"id"`

	Name *string `json:"name,omitempty" url:"name"`
}

type Customer struct {
	Address *[]Address `json:"address,omitempty" url:"address"`

	Id *int64 `json:"id,omitempty" url:"id"`

	Username *string `json:"username,omitempty" url:"username"`
}

type Order struct {
	Complete *bool `json:"complete,omitempty" url:"complete"`
	// 客户信息结构体
	// 用于描述客户相关的信息
	Customer *struct {
		// 用户名
		Username *string `json:"username,omitempty" url:"username"`
		// 用户地址
		// 例如：北京海淀区xxx街道
		// 某某小区
		Address *[]Address `json:"address,omitempty" url:"address"`
		// 用户ID
		Id *int64 `json:"id,omitempty" url:"id"`
	} `json:"customer,omitempty" url:"customer"`

	Id *int64 `json:"id,omitempty" url:"id"`

	// required
	PetId int64 `json:"petId,omitempty" url:"petId"`

	Quantity *int `json:"quantity,omitempty" url:"quantity"`

	// required
	ShipDate time.Time `json:"shipDate,omitempty" url:"shipDate"`
	// Order Status
	Status *string `json:"status,omitempty" url:"status"`
}

type Pet struct {
	Category *Category `json:"category,omitempty" url:"category"`

	Id *int64 `json:"id,omitempty" url:"id"`

	// required
	Name string `json:"name,omitempty" url:"name"`

	// required
	PhotoUrls []string `json:"photoUrls,omitempty" url:"photoUrls"`
	// pet status in the store
	// this is another line for test use
	Status *string `json:"status,omitempty" url:"status"`

	Tags *[]Tag `json:"tags,omitempty" url:"tags"`
}

type Tag struct {
	Id *int64 `json:"id,omitempty" url:"id"`

	Name *string `json:"name,omitempty" url:"name"`
}

type User struct {
	Additional1 *struct {
	} `json:"additional1,omitempty" url:"additional1"`

	Additional2 *struct {
	} `json:"additional2,omitempty" url:"additional2"`

	Avatar *v3.FileModel `json:"avatar,omitempty" url:"avatar"`

	Email *string `json:"email,omitempty" url:"email"`

	FirstName *string `json:"firstName,omitempty" url:"firstName"`

	Id *int64 `json:"id,omitempty" url:"id"`

	LastName *string `json:"lastName,omitempty" url:"lastName"`

	Password *string `json:"password,omitempty" url:"password"`

	Phone *string `json:"phone,omitempty" url:"phone"`
	// User Status
	UserStatus *int `json:"userStatus,omitempty" url:"userStatus"`

	Username *string `json:"username,omitempty" url:"username"`
}
