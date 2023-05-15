package log

import (
	"github.com/harley9293/blotlog/formatter"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"sync"
	"testing"
)

var wg sync.WaitGroup

func TestCoroutine(t *testing.T) {
	reInit()
	defer clear()
	AddRotateHook(&RotateConf{})
	for i := 1; i <= 5; i++ {
		wg.Add(1)
		go func(i int) {
			Debug("test print: index=%d", i)
			wg.Done()
		}(i)
	}
	wg.Wait()

	body, err := ioutil.ReadFile("./log/debug.log")
	if err != nil {
		t.Fatalf("log file not exist, err=%s", err)
	}

	s := string(body)
	lines := strings.Split(s, "\n")
	if len(lines) != 6 {
		t.Fatalf("log line error, line=%d", len(lines))
	}

	index, gid := 0, 0
	for i := 0; i < 5; i++ {
		partList := strings.Split(lines[i], "|")
		if len(partList) != 5 {
			t.Fatalf("log format error, partNum=%d", len(partList))
		}
		curIndex, _ := strconv.Atoi(string(lines[i][len(lines[i])-1]))
		gidFieldList := strings.Split(partList[3], "gid:")
		if len(gidFieldList) != 2 {
			t.Fatalf("gid format error, gid:%s", partList[3])
		}
		curGid, _ := strconv.Atoi(gidFieldList[1])
		if index == 0 {
			index = curIndex
			gid = curGid
		}

		if index-curIndex != gid-curGid {
			t.Fatalf("gid error, index1=%d, index2=%d, gid1=%d, gid2=%d", index, curIndex, gid, curGid)
		}
	}
}

func TestRotateHookNoPass(t *testing.T) {
	reInit()
	defer clear()
	AddRotateHook(&RotateConf{})
	Debug("test print1")
	Error("test print1")

	body, err := ioutil.ReadFile("./log/debug.log")
	if err != nil {
		t.Fatalf("log file not exist, err=%s", err)
	}
	debugStr := string(body)
	lines := strings.Split(debugStr, "\n")
	if len(lines) != 2 {
		t.Fatalf("Debug line error, line=%d", len(lines))
	}

	body, err = ioutil.ReadFile("./log/error.log")
	if err != nil {
		t.Fatalf("log file not exist, err=%s", err)
	}
	errorStr := string(body)
	lines = strings.Split(errorStr, "\n")
	if len(lines) != 2 {
		t.Fatalf("Error line error, line=%d", len(lines))
	}
}

func TestRotateHookWithPass(t *testing.T) {
	reInit()
	defer clear()
	AddRotateHook(&RotateConf{Pass: true})
	Debug("test print2")
	Error("test print2")

	body, err := ioutil.ReadFile("./log/debug.log")
	if err != nil {
		t.Fatalf("log file not exist, err=%s", err)
	}
	debugStr := string(body)
	lines := strings.Split(debugStr, "\n")
	if len(lines) != 3 {
		t.Fatalf("Debug line num error, line=%d", len(lines))
	}

	body, err = ioutil.ReadFile("./log/error.log")
	if err != nil {
		t.Fatalf("log file not exist, err=%s", err)
	}
	errorStr := string(body)
	lines = strings.Split(errorStr, "\n")
	if len(lines) != 2 {
		t.Fatalf("Error line num error, line=%d", len(lines))
	}
}

func TestChangeLevelAfterAddRotateHook(t *testing.T) {
	reInit()
	defer clear()
	AddRotateHook(&RotateConf{})
	Debug("test print3")
	Error("test print3")
	SetLevel(WarnLevel)
	Debug("test print3")
	Error("test print3")

	body, err := ioutil.ReadFile("./log/debug.log")
	if err != nil {
		t.Fatalf("log file not exist, err=%s", err)
	}
	debugStr := string(body)
	lines := strings.Split(debugStr, "\n")
	if len(lines) != 2 {
		t.Fatalf("Debug line num error, line=%d", len(lines))
	}

	body, err = ioutil.ReadFile("./log/error.log")
	if err != nil {
		t.Fatalf("log file not exist, err=%s", err)
	}
	errorStr := string(body)
	lines = strings.Split(errorStr, "\n")
	if len(lines) != 3 {
		t.Fatalf("Error line num error, line=%d", len(lines))
	}
}

// -------------------inner----------------------

func reInit() {
	logger = logrus.New()
	logger.SetFormatter(&formatter.LineFormatter{
		Skip: 10,
	})

	logger.SetLevel(logrus.DebugLevel)
	ConsoleOff()
}

func clear() {
	os.RemoveAll("./log")
}
