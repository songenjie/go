package main

import (
	"github.com/songenjie/go/log"
	"strconv"
	"time"
)

func main() {
	for i := 0; i < 10; i++ {
		log.Warn("songenjiehello" + strconv.Itoa(i))
		time.Sleep(time.Duration(2) * time.Second)
	}
}
