package common

import (
	"fmt"
	"os"
)

func IfErrorExit(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
}

func IfErrorExitv(v interface{}, err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
}
