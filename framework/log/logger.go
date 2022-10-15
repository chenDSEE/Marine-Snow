package log

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

type logger struct {
	format formatter // not support hot update

	mu    sync.Mutex
	level Level
	out   io.Writer
}

var _ Logger = &logger{}

func (l *logger) GetLevel() Level {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.level
}

func (l *logger) SetLevel(level Level) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.level = level
}

func (l *logger) SetFormat(ft FormatType) {
	formatter := newFormatter(ft)
	if formatter == nil {
		return
	}

	l.mu.Lock()
	defer l.mu.Unlock()
	l.format = formatter
}

func (l *logger) WithFields(fs Fields) Logger {
	return &fieldLogger{
		l:      l,
		fields: addFields(nil, fs),
	}
}

func (l *logger) Print(v ...interface{}) {
	Output(l, l.GetLevel(), fmt.Sprint(v...), nil)
}

func (l *logger) Printf(format string, v ...interface{}) {
	Output(l, l.GetLevel(), fmt.Sprintf(format, v...), nil)
}

func (l *logger) Println(v ...interface{}) {
	Output(l, l.GetLevel(), fmt.Sprintln(v...), nil)
}

func (l *logger) Log(level Level, v ...interface{}) {
	Output(l, level, fmt.Sprint(v...), nil)
}

func (l *logger) Logf(level Level, format string, v ...interface{}) {
	Output(l, level, fmt.Sprintf(format, v...), nil)
}

func (l *logger) Logln(level Level, v ...interface{}) {
	Output(l, level, fmt.Sprintln(v...), nil)
}

func (l *logger) Debugf(format string, v ...interface{}) {
	if l.IsEnable(DebugLevel) {
		Output(l, DebugLevel, fmt.Sprintf(format, v...), nil)
	}
}

func (l *logger) Infof(format string, v ...interface{}) {
	if l.IsEnable(InfoLevel) {
		Output(l, InfoLevel, fmt.Sprintf(format, v...), nil)
	}
}

func Output(l *logger, level Level, msg string, fs Fields) {
	now := time.Now()

	l.mu.Lock()
	formatter := l.format
	l.mu.Unlock()

	data, err := formatter.Format(level, now, msg, fs)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Failed to format log data, %v\n", err)
		return
	}

	l.mu.Lock()
	defer l.mu.Unlock()
	if _, err := l.out.Write(data); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Failed to wirte log formatted data, %v\n", err)
	}
}

func (l *logger) IsEnable(level Level) bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	return level >= l.level
}
