package main

import (
	"log"
	"time"
)

//var notFoundErr = errors.New("not found！")
//
//const NO_DATA = "no，data！"

func main() {
	//fmt.Println("Hello，World！")
	logStr := "Current time：" + time.Now().String()
	log.Print(logStr)
	//fmt.Println(NO_DATA)
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
