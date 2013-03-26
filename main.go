package main

import (
	"fmt"
	"time"
)

func ready1() {
	for {
	time.Sleep(time.Minute * 1)
	fmt.Printf("unix time:")
	}
}
func ready2() {
	for {
	time.Sleep(time.Second)
	fmt.Printf("%d\n", time.Now().Unix())
	}
}
func main() {
	go ready1()
	go ready2()
    fmt.Printf("hello, world\n")
	time.Sleep(2e8 * time.Second)
}
