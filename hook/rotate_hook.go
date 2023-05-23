package hook

import (
	"time"

	"github.com/harley9293/blotlog/formatter"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
)

var logLevels = map[logrus.Level]string{
	logrus.DebugLevel: "debug",
	logrus.InfoLevel:  "info",
	logrus.WarnLevel:  "warn",
	logrus.ErrorLevel: "error",
}

func NewLevelRotateHook(l logrus.Level, path string, duration time.Duration, count uint, pass bool) logrus.Hook {
	name := path + "/" + logLevels[l]
	writer, err := rotatelogs.New(
		name+"_%Y%m%d.log",
		rotatelogs.WithLinkName(name+".log"),
		rotatelogs.WithRotationTime(duration),
		rotatelogs.WithRotationCount(count),
	)
	if err != nil {
		logrus.Errorf("config local file system for logger error: %v", err)
	}

	var writerMap = lfshook.WriterMap{}
	writerMap[l] = writer
	if pass {
		for level := l - 1; level >= logrus.ErrorLevel; level-- {
			writerMap[level] = writer
		}
	}

	rotateHook := lfshook.NewHook(writerMap, &formatter.LineFormatter{})

	return rotateHook
}
