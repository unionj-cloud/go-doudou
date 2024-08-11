package main

import (
	"fmt"

	"github.com/unionj-cloud/go-doudou/v2/toolkit/form"
	"log"
	"net/url"
)

// <form method="POST">
//   <input type="text" name="Name" value="joeybloggs"/>
//   <input type="text" name="Age" value="3"/>
//   <input type="text" name="Gender" value="Male"/>
//   <input type="text" name="Address[0].Name" value="26 Here Blvd."/>
//   <input type="text" name="Address[0].Phone" value="9(999)999-9999"/>
//   <input type="text" name="Address[1].Name" value="26 There Blvd."/>
//   <input type="text" name="Address[1].Phone" value="1(111)111-1111"/>
//   <input type="text" name="active" value="true"/>
//   <input type="text" name="MapExample[key]" value="value"/>
//   <input type="text" name="NestedMap[key][key]" value="value"/>
//   <input type="text" name="NestedArray[0][0]" value="value"/>
//   <input type="submit"/>
// </form>

// Address contains address information
type Address struct {
	Name  string
	Phone string
}

// User contains user information
type User struct {
	Name        string
	Age         uint8
	Gender      string
	Address     []Address
	Active      bool `form:"active"`
	MapExample  map[string]string
	NestedMap   map[string]map[string]string
	NestedArray [][]string
	Extra       map[string]interface{} `form:"+"`
	Others      struct {
		Other1 int
		Other2 int
		Others map[string]interface{} `form:"+"`
	}
	Leader *User
}

// use a single instance of Decoder, it caches struct info
var decoder *form.Decoder

func main() {
	decoder = form.NewDecoder()
	decoder.SetNamespacePrefix("[")
	decoder.SetNamespaceSuffix("]")

	// this simulates the results of http.Request's ParseForm() function
	values := parseForm3()

	for i := 0; i < 20; i++ {
		var user ParameterWrapper

		// must pass a pointer
		err := decoder.Decode(&user, values)
		if err != nil {
			log.Panic(err)
		}

		fmt.Printf("%#v\n", user)
	}

	//userMap := make(map[string]interface{})
	//// must pass a pointer
	//err := decoder.Decode(&userMap, values)
	//if err != nil {
	//	log.Panic(err)
	//}
	//
	//fmt.Printf("%#v\n", userMap)

}

// this simulates the results of http.Request's ParseForm() function
func parseForm() url.Values {
	return url.Values{
		"Name":                 []string{"joeybloggs"},
		"Age":                  []string{"3"},
		"Gender":               []string{"Male"},
		"Address[0].Name":      []string{"26 Here Blvd."},
		"Address[0].Phone":     []string{"9(999)999-9999"},
		"Address[1].Name":      []string{"26 There Blvd."},
		"Address[1].Phone":     []string{"1(111)111-1111"},
		"active":               []string{"true"},
		"MapExample[key]":      []string{"value"},
		"MapExample[key1]":     []string{"value1"},
		"NestedMap[key][key]":  []string{"value"},
		"NestedMap[key][key1]": []string{"value1"},
		"NestedArray[0][0]":    []string{"value"},
		"other1":               []string{"1"},
		"other2":               []string{"2"},
		"Leader.Name":          []string{"jack"},
		"Others.Other1":        []string{"1"},
		"Others.Other2":        []string{"2"},
		"Others[Other3]":       []string{"3"},
		"Others[Other4]":       []string{"4"},
		"a":                    []string{"5"},
		"b":                    []string{"b"},
		"c[0]":                 []string{"c"},
		"c[1]":                 []string{"d"},
		"c[2]":                 []string{"e"},
		"d[key]":               []string{"value"},
		"d[key1]":              []string{"value1"},
		"e[key][key]":          []string{"value"},
		"e[key][key1]":         []string{"value1"},
		"f[0][0]":              []string{"value"},
		"f[0][1]":              []string{"value1"},
		"f[1][0]":              []string{"value10"},
	}
}

func parseForm2() url.Values {
	return url.Values{
		"Name":                []string{"joeybloggs"},
		"Age":                 []string{"3"},
		"Gender":              []string{"Male"},
		"Address[0][Name]":    []string{"26 Here Blvd."},
		"Address[0][Phone]":   []string{"9(999)999-9999"},
		"Address[1][Name]":    []string{"26 There Blvd."},
		"Address[1][Phone]":   []string{"1(111)111-1111"},
		"active":              []string{"true"},
		"MapExample[key]":     []string{"value"},
		"NestedMap[key][key]": []string{"value"},
		"NestedArray[0][0]":   []string{"value"},
		"other1":              []string{"1"},
		"other2":              []string{"2"},
		"Leader.Name":         []string{"jack"},
		"Leader[Name]":        []string{"jack"},
		"Others[Other1]":      []string{"1"},
		"Others[Other2]":      []string{"2"},
		"Others[Other3]":      []string{"3"},
		"Others[Other4]":      []string{"4"},
		"a":                   []string{"5"},
		"b":                   []string{"b"},
		"c[0]":                []string{"c"},
		"c[1]":                []string{"d"},
		"c[2]":                []string{"e"},
		"d[key]":              []string{"value"},
		"d[key1]":             []string{"value1"},
		"e[key][key]":         []string{"value"},
		"e[key][key1]":        []string{"value1"},
		"f[0][0]":             []string{"value"},
		"f[0][1]":             []string{"value1"},
		"f[1][0]":             []string{"value10"},
	}
}

type A struct {
	B string `json:"b" form:"b"`
}

type ParameterWrapper struct {
	Parameter A `json:"parameter" form:"parameter"`
}

func parseForm3() url.Values {
	// parameter[product_code]
	return url.Values{
		"parameter[b]": []string{"b"},
	}
}
