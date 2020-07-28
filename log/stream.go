package log

import (
	"io"
	"os"
)

type log_stream struct {
	jtype int
	w     io.Writer
}

func (this *log_stream) write(buf []byte, level int) (n int, err error) {
	return this.w.Write(buf)
}
func (this *log_stream) rotate(level int) {
	return
}
func (this *log_stream) init() error {
	this.setFlag(logStream)
	this.w = os.Stdout
	return nil
}
func (this *log_stream) close() {
	jasonLog.bClose = true
}
func (this *log_stream) getFlag() int {
	return this.jtype
}
func (this *log_stream) setFlag(jtype int) {
	this.jtype = jtype
}
