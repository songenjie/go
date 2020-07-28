package log

import (
	"fmt"
	"github.com/songenjie/go/common"
	"gopkg.in/alecthomas/kingpin.v2"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"sync"
	"time"
)

var (
	App              = kingpin.New(filepath.Base(os.Args[0]), "log")
	logDir           = App.Flag("log_dir", "storage log dir").Default(common.GetWordDir() + "/LOGS/LOG").String()
	logFlag          = App.Flag("log_flag", "log flag  log_stream/log_rotateTime/log_rotateLevel").Default("log_rotateLevel").String()
	logLevel         = App.Flag("log_level", "log level min Trace/Debug/Info/Warn/Error/Fatal").Default("Info").String()
	logRotateCount   = App.Flag("log_rotate_count", "log rotate count ").Default("10").Int()
	timeLevel        = App.Flag("time_level", "time level Second/Minute/Hour/Day/Week/Mounth/Year").Default("Second").String()
	timeInterval     = App.Flag("time_interval", "time interval").Default("1").Int64()
	levelRotateBytes = App.Flag("log_level_ttl_bytes", "log Level Rotate bytes").Default("1024").Int()
	logFlags         = [3]string{"log_stream", "log_rotateTime", "log_rotateLevel"}
	levelNames       = [6]string{"Trace", "Debug", "Info", "Warn", "Error", "Fatal"}
	timeLevels       = [7]string{"Second", "Minute", "Hour", "Day", "week", "Month", "Year"}
	timeSuffixs      = [7]string{"2006-01-02_15-04-05", "2006-01-02_15-04", "2006-01-02_15", "2006-01-02", "2006-01-02", "2006-01", "2006"}
	timeIntervals    = [7]int64{1, 60, 3600, 3600 * 24, 3600 * 24 * 7, 3600 * 24 * 30, 3600 * 24 * 30 * 365}
)

const maxBufPoolSize = 16
const timeFormat = "2006-01-02 15:04:05"

//struct
const (
	logStream = iota
	logRotateTime
	logRotateLevel
)

var jasonLog log

type log struct {
	level
	bClose   bool
	logFlag  int
	logLevel int
	hMutex   sync.Mutex
	bufMutex sync.Mutex
	bufs     [][]byte
}

type level interface {
	setFlag(int)
	getFlag() (int)
	init() error
	write([]byte, int) (int, error)
	close()
}

func init() {
	App.HelpFlag.Short('h')
	kingpin.MustParse(App.Parse(os.Args[1:]))

	jasonLog.initLogFlag()

	switch jasonLog.logFlag {
	case logStream:
		jasonLog.level = new(log_stream)
	case logRotateTime:
		jasonLog.level = new(log_rotateTime)
	case logRotateLevel:
		jasonLog.level = new(log_rotateLevel)
	default:
		return
	}
	common.IfErrorExit(jasonLog.level.init())

	jasonLog.level.setFlag(jasonLog.logFlag)

	jasonLog.initLogLevel()
}

func (this *log) getLogLevel() int {
	return this.logLevel
}

func (this *log) initLogLevel() {
	for index, levelname := range levelNames {
		if levelname == *logLevel {
			jasonLog.logLevel = index
			return
		}
	}
	os.Exit(0)
}

func (this *log) initLogFlag() {
	for index, levelflag := range logFlags {
		if levelflag == *logFlag {
			this.logFlag = index
			return
		}
	}
	os.Exit(0)
}

func (l *log) output(callDepth int, level int, format string, v ...interface{}) {
	if l.bclose() {
		return
	}

	if l.getLogLevel() > level {
		// higher level can be logged
		return
	}

	var s string
	if format == "" {
		s = fmt.Sprint(v...)
	} else {
		s = fmt.Sprintf(format, v...)
	}

	//time
	buf := l.popBuf()
	now := time.Now().Format(timeFormat)
	buf = append(buf, '[')
	buf = append(buf, now...)
	buf = append(buf, "] "...)

	if l.logLevel > 0 {
		_, file, line, ok := runtime.Caller(callDepth)
		if !ok {
			file = "???"
			line = 0
		} else {
			for i := len(file) - 1; i > 0; i-- {
				if file[i] == '/' {
					file = file[i+1:]
					break
				}
			}
		}

		buf = append(buf, file...)
		buf = append(buf, ':')

		buf = strconv.AppendInt(buf, int64(line), 10)
		buf = append(buf, ' ')
	}

	if l.logLevel > 0 {
		buf = append(buf, '[')
		buf = append(buf, levelNames[level]...)
		buf = append(buf, "] "...)
	}

	buf = append(buf, s...)

	if s[len(s)-1] != '\n' {
		buf = append(buf, '\n')
	}

	l.hMutex.Lock()
	common.IfErrorExitv(l.level.write(buf, level))
	l.hMutex.Unlock()
	l.putBuf(buf)
}

//log with Trace level
func Trace(v ...interface{}) {
	jasonLog.output(2, levelTrace, "", v...)
}

//log with Debug level
func Debug(v ...interface{}) {
	jasonLog.output(2, levelDebug, "", v...)
}

//log with info level
func Info(v ...interface{}) {
	jasonLog.output(2, levelInfo, "", v...)
}

//log with warn level
func Warn(v ...interface{}) {
	jasonLog.output(2, levelWarn, "", v...)
}

//log with error level
func Error(v ...interface{}) {
	jasonLog.output(2, levelError, "", v...)
}

//log with fatal level
func Fatal(v ...interface{}) {
	jasonLog.output(2, levelFatal, "", v...)
}

//log with Trace level
func Tracef(format string, v ...interface{}) {
	jasonLog.output(2, levelTrace, format, v...)
}

//log with Debug level
func Debugf(format string, v ...interface{}) {
	jasonLog.output(2, levelDebug, format, v...)
}

//log with info level
func Infof(format string, v ...interface{}) {
	jasonLog.output(2, levelInfo, format, v...)
}

//log with warn level
func Warnf(format string, v ...interface{}) {
	jasonLog.output(2, levelWarn, format, v...)
}

//log with error level
func Errorf(format string, v ...interface{}) {
	jasonLog.output(2, levelError, format, v...)
}

//log with fatal level
func Fatalf(format string, v ...interface{}) {
	jasonLog.output(2, levelFatal, format, v...)
}

func (l *log) popBuf() []byte {
	l.bufMutex.Lock()
	var buf []byte
	if len(l.bufs) == 0 {
		buf = make([]byte, 0, 1024)
	} else {
		buf = l.bufs[len(l.bufs)-1]
		l.bufs = l.bufs[0 : len(l.bufs)-1]
	}
	l.bufMutex.Unlock()

	return buf
}

func (l *log) putBuf(buf []byte) {
	l.bufMutex.Lock()
	if len(l.bufs) < maxBufPoolSize {
		buf = buf[0:0]
		l.bufs = append(l.bufs, buf)
	}
	l.bufMutex.Unlock()
}

func (this *log) bclose() bool {
	return jasonLog.bClose
}
