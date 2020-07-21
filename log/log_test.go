package log

import (
	"fmt"
	"log"
	"os"
	"testing"
)

// func set log level
// init logger
// log.Info

func TestingLevelFileLog(t *testing.T) {
	h, err := NewStreamHandler(os.Stdout)
	if err != nil {

	}
	s := NewDefault(h)
	s.Info("hello world")
	s.Error("mmmm!")

	defer s.Close()

	s.Info("can not log")
}

func TestRotatingFileLog(t *testing.T) {
	path := "./test_log"
	os.RemoveAll(path)

	os.Mkdir(path, 0777)
	fileName := path + "/test"

	h, err := NewRotatingFileHandler(fileName, 10, 10)
	if err != nil {
		log.Println(err)
	}

	buf := make([]byte, 10)
	buf = []byte("akldjsdj")
	log.Println(len(buf))

	h.Write(buf)

	buf = []byte("heelo")
	log.Println(len(buf))

	h.Write(buf)

	log.Println(fileName)

	if _, err := os.Stat(fileName + ".1"); err != nil {
		log.Println(err)
	}
	log.Println(fileName)

	if _, err := os.Stat(fileName + ".2"); err == nil {
		log.Println(err)
	}
	log.Println(fileName)

	h.Write(buf)
	if _, err := os.Stat(fileName + ".2"); err != nil {
		log.Println(err)
	}

	defer h.Close()

	//os.RemoveAll(path)
}

func TestLevelRotate(t *testing.T) {
	path := "./LevelRotateLog"
	fileName := path + "/test"
	os.RemoveAll(path)

	os.Mkdir(path, 0777)

	h, err := NewRotatingFileHandler(fileName, 10, 10)
	if err != nil {
		fmt.Print(err)
	}
	LOG := New(h, Llevel)

	LOG.Info("hello")
	LOG.Error("mmmmads!")
	LOG.Error("abcedfadfkl!")

	LOG.Close()

	LOG.Info("can not log")
}
