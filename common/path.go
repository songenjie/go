package common

import (
	"errors"
	"os"
	"path/filepath"
)

func GetWordDir() string {
	path, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		os.Exit(0)
	}
	return path
}

func MkdirIfPathNotExist(path string) error {
	if IfExists(path) {
		if IsFile(path) {
			return errors.New(path + "is a file path !")
		}
		return nil
	}
	err := os.Mkdir(path, os.ModePerm)
	if err != nil {
		return errors.New("mkdir dir failed of :" + path + " !")
	} else {
		return nil
	}
}

func IfExists(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

func IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

func IsFile(path string) bool {
	return !IsDir(path)
}
