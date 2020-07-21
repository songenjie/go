package test

import (
	"errors"
	"fmt"
	"github.com/songenjie/go/common"
	"gopkg.in/alecthomas/kingpin.v2"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

var LOG log

type log struct {
	dir     string
	flag    int
	level   int
	mutex   sync.Mutex
	bclose  bool
	handler Handler
}

var (
	App            = kingpin.New(filepath.Base(os.Args[0]), "log")
	logDir         = App.Flag("log_dir", "storage log dir").Default(common.GetWordDir() + "/log").String()
	logFlag        = App.Flag("log_flog", "log flag stream/level/time").Default("stream").String()
	logLevel       = App.Flag("log_level", "log level min trace/info/warn/err/panic").Default("info").String()
	logRotateCount = App.Flag("log rotate", "log rotate count ").Default("10").Int()
)

func init() {
	App.HelpFlag.Short('h')
	kingpin.MustParse(App.Parse(os.Args[1:]))
}

func (this log) Init() {

	// 2 flag
	switch strings.ToUpper(*logFlag) {
	case "STREAM":
		this.flag = Fstream
	case "LEVEL":
		this.flag = Flevel
	case "TIME":
		this.flag = Ftime
	default:
		fmt.Println(*logFlag + " not supported !")
		os.Exit(0)
	}

	//3  small level
	switch strings.ToUpper(*logLevel) {
	case LevelName[0]:
		this.level = LevelTrace
	case LevelName[1]:
		this.level = LevelDebug
	case LevelName[2]:
		this.level = LevelInfo
	case LevelName[3]:
		this.level = LevelWarn
	case LevelName[4]:
		this.level = LevelError
	case LevelName[5]:
		this.level = LevelFatal
	default:
		fmt.Println(*logLevel + " not supported !")
		os.Exit(0)
	}

	//1 logdir
	if this.flag != 1 {
		common.IfErrorExit(common.MkdirIfPathNotExist(*logDir))
		this.dir = *logDir
	}

	h, err := this.NewHandler()
	if err != nil {
		fmt.Println("New Handler failed! ")
		os.Exit(0)
	}
	this.handler = h
}

func (this log) NewHandler() (Handler, error) {
	switch this.flag {
	case Fstream:
		return NewStreamHandler(os.Stdout)
	case Ftime:
		return nil, nil
	case Flevel:
		return nil, nil
	}
	return nil, errors.New(" not support !")
}

/*
func hello() {
	LOG.Init()
}*/
