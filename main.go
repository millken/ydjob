package main

import (
	"fmt"
	"time"
	"Config"
	"Job"
)

func ready1() {
	for {
	time.Sleep(time.Second * Config.GetLoopTime().Mail)
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
	go (&Job.Jobmail{}).Run()
	go ready2()
    fmt.Printf("hello, world %s\n", Config.GetBeanstalk().MailQueue)
    fmt.Printf("%d", Config.GetLoopTime().Mail)
	time.Sleep(2e8 * time.Second)
}
