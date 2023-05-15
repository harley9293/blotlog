package formatter

import (
	"bytes"
	"fmt"
	"github.com/sirupsen/logrus"
	"runtime"
	"strconv"
	"strings"
)

type LineFormatter struct {
	Skip int
}

func (f *LineFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	data := ""
	// timestamp
	data += entry.Time.Format("2006-01-02 15:04:05.000")
	// level
	data += "|" + strings.ToUpper(entry.Level.String())
	// caller
	data += "|" + findCaller(f.Skip)
	// gid
	data += "|" + fmt.Sprintf("gid:%s", strconv.FormatUint(getGid(), 10))
	// message
	data += "|" + entry.Message
	return append([]byte(data), '\n'), nil
}

func findCaller(skip int) string {
	file, line, pc := getCaller(skip)
	fullFnName := runtime.FuncForPC(pc)

	fnName := ""
	if fullFnName != nil {
		fnNameStr := fullFnName.Name()
		parts := strings.Split(fnNameStr, ".")
		fnName = parts[len(parts)-1]
	}

	if fnName == "0" {
		fnName = "init"
	}

	tmp := strings.Split(file, "/")
	return fmt.Sprintf("[%s:%d:%s]", tmp[len(tmp)-1], line, fnName)
}

func getCaller(skip int) (string, int, uintptr) {
	pc, file, line, ok := runtime.Caller(skip)
	if !ok {
		return "", 0, pc
	}
	n := 0

	for i := len(file) - 1; i > 0; i-- {
		if file[i] == '/' {
			n++
			if n >= 2 {
				file = file[i+1:]
				break
			}
		}
	}
	return file, line, pc
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
