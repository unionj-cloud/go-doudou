package customtypes

import (
	"fmt"
	"testing"
	"time"

	"github.com/unionj-cloud/go-doudou/v2/toolkit/copier"
)

type User struct {
	CreatedAt Time `json:"created_at"`
}

type Author struct {
	CreatedAt Time `json:"created_at"`
}

func TestTime_UnmarshalJSON(t *testing.T) {
	u := User{}
	fmt.Println(u.CreatedAt)
	a := Author{}
	copier.DeepCopy(&u, &a)
	fmt.Println(a.CreatedAt)
}

func TestTime_UnmarshalJSON1(t *testing.T) {
	u := User{
		CreatedAt: Time(time.Now()),
	}
	fmt.Println(u.CreatedAt)
	a := Author{}
	copier.DeepCopy(&u, &a)
	fmt.Println(a.CreatedAt)
}
