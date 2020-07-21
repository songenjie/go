package test

import (
	"fmt"
	"os"
)

var (
	logLevelRotateBytes = App.Flag("log_level_ttl_bytes", "log Level Rotate bytes").Default("1024").Int()
)

//log level, from low to high, more high means more serious
const (
	LevelTrace = iota
	LevelDebug
	LevelInfo
	LevelWarn
	LevelError
	LevelFatal
)

var LevelName = [6]string{"TRACE", "DEBUG", "INFO", "WARN", "ERROR", "FATAL"}

func (this log) Info() {
	if !this.bclose {

	}
}

type LevelRotatingFileHandlers struct {
	logger [len(LevelName)]*LevelRotatingFileHandler
	//minLevel int
	//maxLevel int

	//maxBytes    int
	//backupCount int
}

type LevelRotatingFileHandler struct {
	fd       *os.File
	headname string
	curBytes int
}

func (this LevelRotatingFileHandlers) InitHeadName() {
	for i, name := range LevelName {
		this.logger[i].headname = name
	}
}

func (h LevelRotatingFileHandler) Rotate() {
	if h.curBytes < *logLevelRotateBytes {
		return
	}
	f, err := h.fd.Stat()
	if err != nil {
		return
	}

	if *logLevelRotateBytes <= 0 {
		return
	} else if f.Size() < int64(*logLevelRotateBytes) {
		h.curBytes = int(f.Size())
		return
	}

	if *logRotateCount > 0 {
		h.fd.Close()

		for i := *logRotateCount - 1; i > 0; i-- {
			sfn := fmt.Sprintf("%s.%d", h.headname, i)
			dfn := fmt.Sprintf("%s.%d", h.headname, i+1)

			os.Rename(sfn, dfn)
		}

		dfn := fmt.Sprintf("%s.1", h.headname)
		os.Rename(h.headname, dfn)

		h.fd, _ = os.OpenFile(h.headname, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		h.curBytes = 0
		f, err := h.fd.Stat()
		if err != nil {
			return
		}
		h.curBytes = int(f.Size())
	}
}
