package log

import (
	"fmt"
	"github.com/songenjie/go/common"
	"os"
	"path"
)

// write implements Handler interface
func (h *log_rotateFile) write(p []byte) (n int, err error) {
	h.rotate()
	n, err = h.fd.Write(p)
	h.curBytes += n
	return
}


type log_rotateFile struct {
	fd          *os.File
	fileName    string
	maxBytes    int
	curBytes    int
	backupCount int
}

func (this *log_rotateFile) rotate() {
	if this.curBytes < this.maxBytes {
		return
	}

	f, err := this.fd.Stat()
	if err != nil {
		return
	}

	if this.maxBytes <= 0 {
		return
	} else if f.Size() < int64(this.maxBytes) {
		this.curBytes = int(f.Size())
		return
	}

	if this.backupCount > 0 {
		this.fd.Close()

		for i := this.backupCount - 1; i > 0; i-- {
			sfn := fmt.Sprintf("%s.%d", this.fileName, i)
			dfn := fmt.Sprintf("%s.%d", this.fileName, i+1)

			os.Rename(sfn, dfn)
		}

		dfn := fmt.Sprintf("%s.1", this.fileName)
		os.Rename(this.fileName, dfn)

		this.fd, _ = os.OpenFile(this.fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		this.curBytes = 0
		f, err := this.fd.Stat()
		if err != nil {
			return
		}
		this.curBytes = int(f.Size())
	}
}


func (this *log_rotateFile) init(filename string) error {
	this.fileName = *logDir + "." + filename
	common.IfErrorExit(common.MkdirIfPathNotExist(path.Dir(this.fileName)))

	if *levelRotateBytes <= 0 {
		return fmt.Errorf("invalid max bytes")
	}

	this.maxBytes = *levelRotateBytes
	this.backupCount = *logRotateCount

	var err error
	this.fd, err = os.OpenFile(this.fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return err
	}

	f, err := this.fd.Stat()
	if err != nil {
		return err
	}
	this.curBytes = int(f.Size())

	return nil
}
