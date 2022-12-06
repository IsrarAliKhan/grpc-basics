package log

import (
	"io"
	"log"
	"os"
	"time"

	"gorm.io/gorm/logger"
)

var (
	Warn          *log.Logger
	Info          *log.Logger
	Debug         *log.Logger
	Error         *log.Logger
	defaultWriter io.Writer
)

func init() {
	defaultWriter = os.Stdout

	//log.SetOutput(defaultWriter)
	Info = log.New(defaultWriter, "[INFO]\t", log.LstdFlags|log.Lshortfile)
	Warn = log.New(defaultWriter, "[WARN]\t", log.LstdFlags|log.Lshortfile)
	Debug = log.New(defaultWriter, "[DEBUG]\t", log.LstdFlags|log.Lshortfile)
	Error = log.New(defaultWriter, "[ERROR]\t", log.LstdFlags|log.Lshortfile)
}

func GetQueryLogger() logger.Interface {
	logConfig := logger.Config{
		SlowThreshold: time.Second,
		LogLevel:      logger.Info,
	}

	queryLogger := log.New(defaultWriter, "[QUERY]\t", log.LstdFlags)

	return logger.New(queryLogger, logConfig)
}
