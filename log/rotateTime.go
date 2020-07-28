package log

import (
	"github.com/songenjie/go/common"
	"io/ioutil"
	"os"
	"path"
	"sort"
	"strings"
	"time"
)

type log_rotateTime struct {
	jtype      int
	level      int // suffix == level
	fd         *os.File
	baseName   string
	interval   int64
	suffix     string
	rolloverAt int64
}

func (this *log_rotateTime) initTimeLevel() {
	for index, timelevel := range timeLevels {
		if timelevel == *timeLevel {
			this.level = index
			this.suffix = timeSuffixs[index]
			this.interval = timeIntervals[index]
		}
	}
}

func (this *log_rotateTime) write(buf []byte, _ int) (n int, err error) {
	this.rotate()
	return this.fd.Write(buf)
	//code
}

func (this *log_rotateTime) rotate() {
	now := time.Now()

	if this.rolloverAt <= now.Unix() {
		this.baseName = *logDir + "T_" + now.Format(this.suffix)
		this.fd, _ = os.OpenFile(this.baseName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		this.rolloverAt = time.Now().Unix() + this.interval
		this.rrotate()
	}

}

func (this *log_rotateTime) rrotate() {
	dir, _ := ioutil.ReadDir(path.Dir(*logDir))
	var filenames = []string{}
	for _, fi := range dir {
		if !fi.IsDir() && strings.Contains(fi.Name(), "LOGT") {
			filenames = append(filenames, fi.Name())
		}
	}
	sort.Strings(filenames)
	for i := 0; i < len(filenames)-*logRotateCount; i++ {
		common.IfErrorExit(os.Remove(path.Dir(*logDir) + "/" + filenames[i]))
	}
}

func (this *log_rotateTime) init() error {
	now := time.Now()
	this.setFlag(logRotateTime)
	this.initTimeLevel()
	this.baseName = *logDir + "T_" + now.Format(this.suffix)

	common.MkdirIfPathNotExist(path.Dir(*logDir))

	this.interval = this.interval * (*timeInterval)

	var err error
	this.fd, err = os.OpenFile(this.baseName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return err
	}

	fInfo, _ := this.fd.Stat()
	this.rolloverAt = fInfo.ModTime().Unix() + this.interval
	this.rrotate()
	return nil
}

func (this *log_rotateTime) close() {
	common.IfErrorExit(this.fd.Close())
	jasonLog.bClose = true
}

func (this *log_rotateTime) getFlag() int {
	return this.jtype
}

func (this *log_rotateTime) setFlag(jtype int) {
	this.jtype = jtype
}
