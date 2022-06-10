package main

import "time"

func testCall(a, b string) string {
	return a + "," + b
}
func testMain() {
	testCall("a", testCall("b", testCall(time.Now().String(), "！")))
	testCall("a！b"+testCall("x", "y"), "z")
	if "aa！" != "b" {
		testCall("a", func() string { return "bb！" }())
	}
}
