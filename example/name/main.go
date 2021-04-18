package main

import (
	"encoding/json"
	"example/name/vo"
	"fmt"
	"log"
	"strings"
	"time"
)

func main() {
	sch := vo.School{
		Name:     "Beijing University",
		Address:  "Beijing",
		CreateAt: time.Now(),
	}
	company := vo.Company{
		Name:     "Stq",
		CreateAt: time.Now().AddDate(0, -1, 0),
	}
	stu := vo.Student{
		Name:      "wubin",
		Age:       30,
		TestScore: 99,
		School:    sch,
		Company:   company,
	}

	data, err := json.Marshal(stu)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(data))

	jsonStr := `{"age":30,"name":"wubin","testScore":99}`
	var stu1 vo.Student
	err = json.Unmarshal([]byte(jsonStr), &stu1)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(stu1)

	s := "vo.go"
	fmt.Println(strings.TrimSuffix(s, ".go"))
}
