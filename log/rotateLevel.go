package log

import "github.com/songenjie/go/common"

//for
const (
	levelTrace = iota
	levelDebug
	levelInfo
	levelWarn
	levelError
	levelFatal
)

type log_rotateLevel struct {
	log_rotateFiles map[int]*log_rotateFile
	jtype           int
	Level           int
}

func (this *log_rotateLevel) write(buf []byte, level int) (n int, err error) {
	for i := jasonLog.logLevel; i <= level; i++ {
		common.IfErrorExitv(this.log_rotateFiles[i].write(buf))
	}
	return 0, nil
}

func (this *log_rotateLevel) rotate() {
	for _, filerotate := range this.log_rotateFiles {
		filerotate.rotate()
	}
}

func (this *log_rotateLevel) init() error {
	this.log_rotateFiles = make(map[int]*log_rotateFile)
	for i := jasonLog.logLevel; i < len(levelNames); i++ {
		this.log_rotateFiles[i] = new(log_rotateFile)
		common.IfErrorExit(this.log_rotateFiles[i].init(levelNames[i]))
	}
	return nil
}

func (this *log_rotateLevel) close() {
	for _, fs := range this.log_rotateFiles {
		common.IfErrorExit(fs.fd.Close())
	}
	jasonLog.bClose = true
}

func (this *log_rotateLevel) getFlag() int {
	return this.jtype
}

func (this *log_rotateLevel) setFlag(jtype int) {
	this.jtype = jtype
}
