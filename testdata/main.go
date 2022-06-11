package main

import (
	"errors"
	"fmt"
	"log"
	"time"
)

var notFoundErr = errors.New("not found！")

const NO_DATA = "no，data！"

func main() {
	fmt.Println("Hello，World！")
	logStr := "Current time：" + time.Now().String()
	log.Print(logStr)
	fmt.Println(NO_DATA)
	if logStr == "한국어" {
		fmt.Println("にほんご")
	}
	log.Println(":) 😁😁😁")
}

//
//func testCase(url string) string {
//	test := `test`
//	if url == "test" {
//		return test
//	}
//	switch url {
//	case "moo":
//		return ""
//	}
//	return "foo"
//}
