package test

import "os"

const TimeFormat = "2006/01/02 15:04:05"

var (
	logTimeRotateTTl = App.Flag("log_time_ttl", "log Time TTl by hour/day/month/year").Default("day").String()
)

type TimeRotatingFileHandler struct {
	fd *os.File
	baseName   string
	interval   int64
	suffix     string
	rolloverAt int64
}

const (
	//WhenSecond = iota
	WhenMinute = iota
	WhenHour
	WhenDay
	WhenMonth
	WhenYear
)
