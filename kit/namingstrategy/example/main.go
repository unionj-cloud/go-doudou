package main

import (
	"cloud/unionj/papilio/kit/namingstrategy/example/vo"
	"encoding/json"
	"fmt"
	"log"
)

func main() {
	stu := vo.Student{
		Name:      "wubin",
		Age:       30,
		TestScore: 99,
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
}
