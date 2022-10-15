package log

import (
	"os"
)

var std = &logger{
	format: newFormatter(TextFormat),
	level:  DebugLevel,
	out:    os.Stdout,
}

func GetLevel() Level {
	std.mu.Lock()
	defer std.mu.Unlock()
	return std.level
}

func SetLevel(level Level) {
	std.mu.Lock()
	defer std.mu.Unlock()
	std.level = level
}

func WithFields(fs Fields) Logger {
	return std.WithFields(fs)
}

func Print(v ...interface{}) {
	std.Print(v...)
}

func Printf(format string, v ...interface{}) {
	std.Printf(format, v...)
}

func Println(v ...interface{}) {
	std.Println(v...)
}

func Log(level Level, v ...interface{}) {
	std.Log(level, v...)
}

func Logf(level Level, format string, v ...interface{}) {
	std.Logf(level, format, v...)
}

func Logln(level Level, v ...interface{}) {
	std.Logln(level, v...)
}

func Debugf(format string, v ...interface{}) {
	std.Debugf(format, v...)
}

func Infof(format string, v ...interface{}) {
	std.Infof(format, v...)
}
