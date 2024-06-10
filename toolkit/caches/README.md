# Gorm Caches

Gorm Caches plugin using database request reductions (easer), and response caching mechanism provide you an easy way to optimize database performance.

## Features

- Database request reduction. If three identical requests are running at the same time, only the first one is going to be executed, and its response will be returned for all.
- Database response caching. By implementing the Cacher interface, you can easily setup a caching mechanism for your database queries.
- Supports all databases that are supported by gorm itself.

## Install

```bash
go get -u github.com/go-gorm/caches/v2
```

## Usage

Configure the `easer`, and the `cacher`, and then load the plugin to gorm.

```go
package main

import (
	"fmt"
	"sync"

	"github.com/go-gorm/caches"
	"github.com/wubin1989/mysql"
	"github.com/wubin1989/gorm"
)

func main() {
	db, _ := gorm.Open(
		mysql.Open("DATABASE_DSN"),
		&gorm.Config{},
	)
	cachesPlugin := &caches.Caches{Conf: &caches.Config{
		Easer: true,
		Cacher: &yourCacherImplementation{},
	}}
	_ = db.Use(cachesPlugin)
}
```

## Easer Example

```go
package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/go-gorm/caches"
	"github.com/wubin1989/mysql"
	"github.com/wubin1989/gorm"
)

type UserRoleModel struct {
	gorm.Model
	Name string `gorm:"unique"`
}

type UserModel struct {
	gorm.Model
	Name   string
	RoleId uint
	Role   *UserRoleModel `gorm:"foreignKey:role_id;references:id"`
}

func main() {
	db, _ := gorm.Open(
		mysql.Open("DATABASE_DSN"),
		&gorm.Config{},
	)

	cachesPlugin := &caches.Caches{Conf: &caches.Config{
		Easer: true,
	}}

	_ = db.Use(cachesPlugin)

	_ = db.AutoMigrate(&UserRoleModel{})

	_ = db.AutoMigrate(&UserModel{})

	adminRole := &UserRoleModel{
		Name: "Admin",
	}
	db.FirstOrCreate(adminRole, "Name = ?", "Admin")

	guestRole := &UserRoleModel{
		Name: "Guest",
	}
	db.FirstOrCreate(guestRole, "Name = ?", "Guest")

	db.Save(&UserModel{
		Name: "ktsivkov",
		Role: adminRole,
	})
	db.Save(&UserModel{
		Name: "anonymous",
		Role: guestRole,
	})

	var (
		q1Users []UserModel
		q2Users []UserModel
	)
	wg := &sync.WaitGroup{}
	wg.Add(2)
	go func() {
		db.Model(&UserModel{}).Joins("Role").Find(&q1Users, "Role.Name = ? AND Sleep(1) = false", "Admin")
		wg.Done()
	}()
	go func() {
		time.Sleep(500 * time.Millisecond)
		db.Model(&UserModel{}).Joins("Role").Find(&q2Users, "Role.Name = ? AND Sleep(1) = false", "Admin")
		wg.Done()
	}()
	wg.Wait()

	fmt.Println(fmt.Sprintf("%+v", q1Users))
	fmt.Println(fmt.Sprintf("%+v", q2Users))
}
```

## Cacher Example

```go
package main

import (
	"fmt"
	"sync"

	"github.com/go-gorm/caches"
	"github.com/wubin1989/mysql"
	"github.com/wubin1989/gorm"
)

type UserRoleModel struct {
	gorm.Model
	Name string `gorm:"unique"`
}

type UserModel struct {
	gorm.Model
	Name   string
	RoleId uint
	Role   *UserRoleModel `gorm:"foreignKey:role_id;references:id"`
}

type dummyCacher struct {
	store *sync.Map
}

func (c *dummyCacher) init() {
	if c.store == nil {
		c.store = &sync.Map{}
	}
}

func (c *dummyCacher) Get(key string) *caches.Query {
	c.init()
	val, ok := c.store.Load(key)
	if !ok {
		return nil
	}

	return val.(*caches.Query)
}

func (c *dummyCacher) Store(key string, val *caches.Query) error {
	c.init()
	c.store.Store(key, val)
	return nil
}

func main() {
	db, _ := gorm.Open(
		mysql.Open("DATABASE_DSN"),
		&gorm.Config{},
	)

	cachesPlugin := &caches.Caches{Conf: &caches.Config{
		Cacher: &dummyCacher{},
	}}

	_ = db.Use(cachesPlugin)

	_ = db.AutoMigrate(&UserRoleModel{})

	_ = db.AutoMigrate(&UserModel{})

	adminRole := &UserRoleModel{
		Name: "Admin",
	}
	db.FirstOrCreate(adminRole, "Name = ?", "Admin")

	guestRole := &UserRoleModel{
		Name: "Guest",
	}
	db.FirstOrCreate(guestRole, "Name = ?", "Guest")

	db.Save(&UserModel{
		Name: "ktsivkov",
		Role: adminRole,
	})
	db.Save(&UserModel{
		Name: "anonymous",
		Role: guestRole,
	})

	var (
		q1Users []UserModel
		q2Users []UserModel
	)

	db.Model(&UserModel{}).Joins("Role").Find(&q1Users, "Role.Name = ? AND Sleep(1) = false", "Admin")
	fmt.Println(fmt.Sprintf("%+v", q1Users))

	db.Model(&UserModel{}).Joins("Role").Find(&q2Users, "Role.Name = ? AND Sleep(1) = false", "Admin")
	fmt.Println(fmt.Sprintf("%+v", q2Users))
}
```

## License

MIT license.

## Easer
The easer is an adjusted version of the [ServantGo](https://github.com/ktsivkov/servantgo) library to fit the needs of this plugin.
