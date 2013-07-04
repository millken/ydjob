package main

import (
	"time"
	"Job"
)


func main() {
	go (&Job.Jobmail{}).Run()
	go (&Job.Jobsms{}).Run()
	go (&Job.Firewall{}).Run()
	
	time.Sleep(2e8 * time.Second)
}
