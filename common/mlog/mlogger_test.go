package mlog

import (
	"testing"
)

func Test_logger(t *testing.T) {
	//var t1 *Toollogger
	t1 := &Toollogger{}
	t1.Openfile()
	logger := t1.Logger()
	logger.Println("AAAAAA")
}
