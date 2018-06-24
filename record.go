package logx

import (
	"time"
)

type (
	Record interface {
		Level() Level
		Line() int
		Time() time.Time
		Stack([]byte) int
		Prefix() string
		Message() string
		File() string
		Func() string
	}
)

type (
	record struct {
		t       time.Time
		line    uint32
		level   Level
		prefix  string
		fn      string
		file    string
		message string
		pcs     []uintptr
	}
)

func (r *record) Time() time.Time {
	return r.t
}

func (r *record) Message() string {
	return r.message
}

func (r *record) Prefix() string {
	return r.prefix
}

func (r *record) Line() int {
	return int(r.line)
}

func (r *record) File() string {
	return r.file
}

func (r *record) Func() string {
	return r.fn
}

func (r *record) Level() Level {
	return r.level
}

func (r *record) Stack(buffer []byte) int {
	if len(buffer) == 0 {
		return 0
	}
	if len(r.pcs) == 0 {
		return 0
	}
	return getStack(r.pcs, buffer)
}

func newRecord(level Level, line int, message, prefix, fn, file string) *record {
	return &record{
		t:       time.Now(),
		level:   level,
		line:    uint32(line),
		prefix:  prefix,
		fn:      fn,
		file:    file,
		message: message,
	}
}
