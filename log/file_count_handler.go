package log

import (
	"fmt"
	"os"
	"path"
)

type Jasonlogger struct {
	logger      map[int]*RotatingFileHandler
	headname    string
	maxBytes    int
	minLevel    int
	maxLevel    int
	backupCount int
}

func NewJasonLogger(headname string, minLevel int, maxLevel int, maxBytes int, backupCount int) Jasonlogger {
	jason := Jasonlogger{
		logger:      nil,
		headname:    headname,
		maxBytes:    maxBytes,
		minLevel:    minLevel,
		maxLevel:    maxLevel,
		backupCount: backupCount,
	}
	for i := minLevel; i < maxLevel; i++ {
		jason.logger[i] = &RotatingFileHandler{}
	}
	return jason
}

type LogRotateFile struct {
	filename string
	fd       *os.File
}

func NewLogRotateFile (filename string) LogRotateFile {
	return *new(LogRotateFile)
}



//RotatingFileHandler writes log a file, if file size exceeds maxBytes,
//it will backup current file and open a new one.
//
//max backup file number is set by backupCount, it will delete oldest if backups too many.
type RotatingFileHandler struct {
	fd *os.File

	fileName    string
	maxBytes    int
	curBytes    int
	backupCount int
}

func NewRotatingFileHandler(fileName string, maxBytes int, backupCount int) (*RotatingFileHandler, error) {
	dir := path.Dir(fileName)
	os.MkdirAll(dir, 0777)

	h := new(RotatingFileHandler)

	if maxBytes <= 0 {
		return nil, fmt.Errorf("invalid max bytes")
	}

	h.fileName = fileName
	h.maxBytes = maxBytes
	h.backupCount = backupCount

	var err error
	h.fd, err = os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	f, err := h.fd.Stat()
	if err != nil {
		return nil, err
	}
	h.curBytes = int(f.Size())

	return h, nil
}

func (h *RotatingFileHandler) Write(p []byte) (n int, err error) {
	h.doRollover()
	n, err = h.fd.Write(p)
	h.curBytes += n
	return
}

func (h *RotatingFileHandler) Close() error {
	if h.fd != nil {
		return h.fd.Close()
	}
	return nil
}

func (h *RotatingFileHandler) doRollover() {

	if h.curBytes < h.maxBytes {
		return
	}

	f, err := h.fd.Stat()
	if err != nil {
		return
	}

	if h.maxBytes <= 0 {
		return
	} else if f.Size() < int64(h.maxBytes) {
		h.curBytes = int(f.Size())
		return
	}

	if h.backupCount > 0 {
		h.fd.Close()

		for i := h.backupCount - 1; i > 0; i-- {
			sfn := fmt.Sprintf("%s.%d", h.fileName, i)
			dfn := fmt.Sprintf("%s.%d", h.fileName, i+1)

			os.Rename(sfn, dfn)
		}

		dfn := fmt.Sprintf("%s.1", h.fileName)
		os.Rename(h.fileName, dfn)

		h.fd, _ = os.OpenFile(h.fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		h.curBytes = 0
		f, err := h.fd.Stat()
		if err != nil {
			return
		}
		h.curBytes = int(f.Size())
	}
}
