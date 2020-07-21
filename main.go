package main

import (
	"fmt"
	"github.com/songenjie/go/log"
	"net/http"
	"os"
	"path/filepath"
)

var (
	engine_name     = "go"
	log_file_path   = "./logs"
	log_low_level   = log.LevelInfo
	log_high_level  = log.LevelError
	log_rotate_size = 10;
	log_rotate_num  = 10;
	LOG             = &log.Logger{}
)

func init() {

	fmt.Println("hello1")
	path := log_file_path
	engine_name = path + "/" + engine_name
	os.RemoveAll(path)
	os.Mkdir(path, 0777)

	h, err := log.NewRotatingFileHandler(engine_name, log_rotate_size, log_rotate_num)
	if err != nil {
		fmt.Print(err)
	}
	LOG = log.New(h, log.Llevel)

}

/*
syn.groups
*/



func main() {
	fmt.Println(os.Args[0])
	fmt.Println(filepath.Abs(filepath.Dir(os.Args[0])))
	//defer LOG.Close()
	fmt.Println("hello111")
	fmt.Println(os.Args[0])
	fmt.Println("hello111")
	LOG.Info("hello")
	go logs()
	http.ListenAndServe(":9090", nil)

}

func logs() {
	LOG.Error(log.LevelError)
	LOG.Error(log.Llevel)
}
