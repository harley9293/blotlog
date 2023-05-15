package log

import "time"

type Level string

const (
	DebugLevel Level = "DEBUG"
	InfoLevel  Level = "INFO"
	WarnLevel  Level = "WARN"
	ErrorLevel Level = "ERROR"
)

type RotateConf struct {
	Path  string
	Time  time.Duration
	Count uint
	Pass  bool
}

func (r *RotateConf) Fill() {
	if r.Path == "" {
		r.Path = "./log"
	}
	if r.Count == 0 {
		r.Count = 7
	}
	if r.Time == 0 {
		r.Time = time.Hour * 24
	}
}
