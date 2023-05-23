package formatter

import (
	"bytes"
	"fmt"
	"runtime"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
)

type LineFormatter struct{}

func (f *LineFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	data := ""
	// timestamp
	data += entry.Time.Format("2006-01-02 15:04:05.000")
	// level
	data += "|" + strings.ToUpper(entry.Level.String())
	// caller
	data += "|" + findCaller()
	// gid
	data += "|" + fmt.Sprintf("gid:%s", strconv.FormatUint(getGid(), 10))
	// message
	data += "|" + entry.Message
	return append([]byte(data), '\n'), nil
}

func findCaller() string {
	pc := make([]uintptr, 20)
	runtime.Callers(2, pc)
	frames := runtime.CallersFrames(pc)
	for {
		frame, more := frames.Next()
		if strings.Contains(frame.Function, "blotlog.") {
			frame, _ = frames.Next()
			fnName := ""
			if fnName == "0" {
				fnName = "init"
			}
			if frame.Func.Name() != "" {
				parts := strings.Split(frame.Func.Name(), ".")
				fnName = parts[len(parts)-1]
			}
			tmp := strings.Split(frame.File, "/")
			return fmt.Sprintf("[%s:%d:%s]", tmp[len(tmp)-1], frame.Line, fnName)
		}

		if !more {
			break
		}
	}

	return ""
}

func getGid() uint64 {
	b := make([]byte, 64)
	b = b[:runtime.Stack(b, false)]
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	b = b[:bytes.IndexByte(b, ' ')]
	n, err := strconv.ParseUint(string(b), 10, 64)
	if err != nil {
		return 0
	}
	return n
}
