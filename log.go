package log

import (
	"io"

	"github.com/harley9293/blotlog/formatter"
	"github.com/harley9293/blotlog/hook"
	"github.com/sirupsen/logrus"
)

var logger *logrus.Logger = nil

func init() {
	logger = logrus.New()
	logger.SetFormatter(&formatter.LineFormatter{})

	logger.SetLevel(logrus.DebugLevel)
}

// -------------------config----------------------

func SetLevel(level Level) {
	lvl, _ := logrus.ParseLevel(string(level))
	logger.SetLevel(lvl)
}

func ConsoleOff() {
	logger.SetOutput(io.Discard)
}

func AddRotateHook(conf *RotateConf) {
	if conf == nil {
		conf = new(RotateConf)
	}
	conf.Fill()
	for level := logger.GetLevel(); level >= logrus.ErrorLevel; level-- {
		logger.AddHook(hook.NewLevelRotateHook(level, conf.Path, conf.Time, conf.Count, conf.Pass))
	}
}

// -------------------print----------------------

func Debug(str string, args ...interface{}) {
	logger.Debugf(str, args...)
}

func Info(str string, args ...interface{}) {
	logger.Infof(str, args...)
}

func Warn(str string, args ...interface{}) {
	logger.Warnf(str, args...)
}

func Error(str string, args ...interface{}) {
	logger.Errorf(str, args...)
}
