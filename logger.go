package logx

import (
	"fmt"
	"os"
)

type (
	Logger struct {
		prefix string
		level  Level
		writer Writer
	}
)

func New(writer Writer) *Logger {
	if nil == writer {
		panic("log.New: writer cannot be empty")
	}
	return &Logger{
		writer: writer,
	}
}

func (logger *Logger) SetLevel(level Level) {
	logger.level = level
}

func (logger *Logger) GetLevel() Level {
	return logger.level
}

func (logger *Logger) SetPrefix(prefix string) {
	logger.prefix = prefix
}

func (logger *Logger) GetPrefix() string {
	return logger.prefix
}

func (logger *Logger) Debug(args ...interface{}) {
	logger.print(DebugLevel, args...)
}

func (logger *Logger) Info(args ...interface{}) {
	logger.print(InfoLevel, args...)
}

func (logger *Logger) Error(args ...interface{}) {
	logger.print(ErrorLevel, args...)
}

func (logger *Logger) Fatal(args ...interface{}) {
	logger.print(FatalLevel, args...)
}

func (logger *Logger) Panic(args ...interface{}) {
	logger.print(PanicLevel, args...)
}

func (logger *Logger) Debugf(format string, args ...interface{}) {
	logger.printf(DebugLevel, format, args...)
}

func (logger *Logger) Infof(format string, args ...interface{}) {
	logger.printf(InfoLevel, format, args...)
}

func (logger *Logger) Errorf(format string, args ...interface{}) {
	logger.printf(ErrorLevel, format, args...)
}

func (logger *Logger) Fatalf(format string, args ...interface{}) {
	logger.printf(FatalLevel, format, args...)
}

func (logger *Logger) Panicf(format string, args ...interface{}) {
	logger.printf(PanicLevel, format, args...)
}

func (logger *Logger) printf(level Level, format string, args ...interface{}) {
	if level >= logger.level && level < Disabled {
		logger.output(level, fmt.Sprintf(format, args...))
	}
}

func (logger *Logger) print(level Level, args ...interface{}) {
	if level >= logger.level && level < Disabled {
		logger.output(level, fmt.Sprint(args...))
	}
}

func (logger *Logger) output(level Level, message string) {
	r := newRecord(level, message, logger.prefix)
	if level >= ErrorLevel {
		r.pcs = getPcs(5)
	}
	if err := logger.writer.Write(r); err != nil {
		fmt.Fprintf(os.Stderr, "logger: write log record error, %s", err)
	}
	if FatalLevel == level {
		os.Exit(1)
	} else if PanicLevel == level {
		panic(message)
	}
}
