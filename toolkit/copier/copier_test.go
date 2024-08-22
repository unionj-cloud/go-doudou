package copier

import (
	"fmt"
	"github.com/samber/lo"
	"github.com/unionj-cloud/go-doudou/v2/toolkit/customtypes"
	"github.com/wubin1989/gorm"
	"testing"
)

type Family struct {
	Father string
	Mather string
	Pets   map[string]string
	Toys   []string
}

type TestStruct struct {
	Name   string
	Age    int
	Family Family
}

type FamilyShadow struct {
	Father string
	Mather string
	Pets   map[string]string
}

type TestStructShadow struct {
	Name   string
	Family FamilyShadow
}

func TestDeepCopy(t *testing.T) {
	pets := make(map[string]string)
	pets["a"] = "dog"
	pets["b"] = "cat"

	family := Family{
		Father: "Jack",
		Mather: "Lily",
		Pets:   pets,
		Toys: []string{
			"car",
			"lego",
		},
	}
	src := TestStruct{
		Name:   "Rose",
		Age:    18,
		Family: family,
	}

	var target TestStructShadow

	type args struct {
		src    interface{}
		target interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			args: args{
				src:    src,
				target: &target,
			},
			wantErr: false,
		},
		{
			args: args{
				src:    nil,
				target: nil,
			},
			wantErr: false,
		},
		{
			name: "",
			args: args{
				src:    make(chan string),
				target: &target,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := DeepCopy(tt.args.src, tt.args.target); (err != nil) != tt.wantErr {
				t.Errorf("DeepCopy() error = %v, wantErr %v", err, tt.wantErr)
			}
			fmt.Printf("%#v\n", tt.args.target)
		})
	}
}

func TestDeepCopy_ShouldHasError(t *testing.T) {
	pets := make(map[string]string)
	pets["a"] = "dog"
	pets["b"] = "cat"

	family := Family{
		Father: "Jack",
		Mather: "Lily",
		Pets:   pets,
		Toys: []string{
			"car",
			"lego",
		},
	}
	src := TestStruct{
		Name:   "Rose",
		Age:    18,
		Family: family,
	}

	var target TestStructShadow

	type args struct {
		src    interface{}
		target interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "TestDeepCopy",
			args: args{
				src:    src,
				target: target,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := DeepCopy(tt.args.src, tt.args.target); (err != nil) != tt.wantErr {
				t.Errorf("DeepCopy() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDeepCopy2(t *testing.T) {
	t1 := `{"name":"jack", "age": "18.0"}`
	type Person struct {
		Name string
		Age  float64 `json:"age,string"`
	}
	var p Person
	json.Unmarshal([]byte(t1), &p)

	type Student struct {
		Name string
		Age  int `json:"age,string"`
	}
	var s Student
	DeepCopy(p, &s)

	fmt.Println(p)
	fmt.Println(s)

	// Output:
	// {jack 18}
	//{jack 18}

}

func TestDeepCopy3(t *testing.T) {
	//t1 := `{"name":"jack", "age": 18.0, "school": "beijing"}`
	p := make(map[string]interface{})
	p["name"] = nil
	p["age"] = lo.ToPtr(18)
	//dec := decoder.NewDecoder(t1)
	//dec.UseInt64()
	//dec.Decode(&p)
	//ddd, _ := json.Marshal(p)
	//fmt.Println(string(ddd))
	type Student struct {
		Name *string `json:"name"`
		Age  *int    `json:"age,string"`
	}
	var s Student
	if err := DeepCopy(p, &s); err != nil {
		panic(err)
	}

	fmt.Println(p)
	fmt.Println(s)

	// Output:
	// {jack 18}
	//{jack 18}

}

func TestDeepCopy4(t *testing.T) {
	//t1 := `{"name":"jack", "age": 18.0, "school": "beijing"}`
	p := make(map[string]interface{})
	p["name"] = nil
	p["age"] = "18"
	type Student struct {
		Name *string `json:"name"`
		Age  *int64  `json:"age,string"`
	}
	var s Student
	if err := DeepCopy(p, &s); err != nil {
		panic(err)
	}

	fmt.Println(p)
	fmt.Println(s)

	// Output:
	// {jack 18}
	//{jack 18}

}

func TestDeepCopy5(t *testing.T) {
	t1 := `{"updated_at":"2024-07-27 13:34:27","deleted_at":"2006-01-02T15:04:05Z"}`
	p := make(map[string]interface{})
	json.Unmarshal([]byte(t1), &p)

	type Student struct {
		UpdatedAt *customtypes.Time `json:"updated_at"`
		DeletedAt gorm.DeletedAt    `json:"deleted_at"`
	}
	var s Student
	if err := DeepCopy(p, &s); err != nil {
		panic(err)
	}

	fmt.Println(p)
	fmt.Println(s)

	// Output:
	// {jack 18}
	//{jack 18}

}
