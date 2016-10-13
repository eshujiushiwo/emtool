package mlogger

import (
	"github.com/uber-go/zap"
	"os"
)

func Logger(loglevel, loginfo string) {

	logpath := "./log.log"

	logfile, err1 := os.OpenFile(logpath, os.O_RDWR|os.O_CREATE, 0666)
	defer logfile.Close()

	if err1 != nil {

		os.Exit(-1)
	}

	logger := zap.New(
		zap.NewJSONEncoder(zap.RFC3339Formatter("@timestamp")),
		zap.Output(logfile),
		zap.Output(os.Stdout),
	)
	if loglevel == "i" {
		logger.Info(loginfo)
	} else if loglevel == "w" {
		logger.Warn(loginfo)
	} else if loglevel == "e" {
		logger.Error(loginfo)
	}

}
