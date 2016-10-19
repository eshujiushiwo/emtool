package mlog

import (
	"io"
	"log"
	"os"
)

var logger *log.Logger

type Toollogger struct {
	logfiles io.Writer
}

func (t1 *Toollogger) Openfile() io.Writer {
	var multi_logfile []io.Writer

	logpath := "./log.log"

	logfile, err1 := os.OpenFile(logpath, os.O_RDWR|os.O_CREATE, 0666)
	defer logfile.Close()

	if err1 != nil {

		os.Exit(-1)
	}

	multi_logfile = []io.Writer{
		logfile,
		os.Stdout,
	}
	t1.logfiles = io.MultiWriter(multi_logfile...)
	return t1.logfiles

}

func (t1 *Toollogger) Logger() *log.Logger {

	logger = log.New(t1.logfiles, "\r\n", log.Ldate|log.Ltime|log.Lshortfile)

	return logger

}

var globalToolLogger *Toollogger

func init() {
	if globalToolLogger == nil {
		globalToolLogger = &Toollogger{}
	}
}

func Openfile() io.Writer {
	return globalToolLogger.Openfile()
}
func Logger() *log.Logger {
	return globalToolLogger.Logger()
}
