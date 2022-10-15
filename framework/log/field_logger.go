package log

import "fmt"

type Fields map[string]string

// all configuration copy from fieldLogger.
// Level, formatter should not change for fieldLogger
type fieldLogger struct {
	l      *logger
	fields Fields
}

func (fl *fieldLogger) GetLevel() Level {
	return fl.l.GetLevel()
}

// level can not change via fieldLogger
func (fl *fieldLogger) SetLevel(Level) {
	return
}

// format can not change via fieldLogger
func (fl *fieldLogger) SetFormat(FormatType) {
	return
}

// COW to avoid Fields be modified and race condition
func (fl *fieldLogger) WithFields(fs Fields) Logger {
	return &fieldLogger{
		l:      fl.l,
		fields: addFields(fl.fields, fs),
	}
}

func addFields(oldFields Fields, newFields Fields) Fields {
	fields := make(Fields, len(oldFields)+len(newFields)+3) // 3 for time, msg and level

	// deep copy to avoid Fields be modified and race condition(COW, copy-on-write)
	for k, v := range oldFields {
		fields[k] = v
	}

	for k, v := range newFields {
		fields[k] = v
	}

	return fields
}

func (fl *fieldLogger) Print(v ...interface{}) {
	Output(fl.l, fl.GetLevel(), fmt.Sprint(v...), fl.fields)
}

func (fl *fieldLogger) Printf(format string, v ...interface{}) {
	Output(fl.l, fl.GetLevel(), fmt.Sprintf(format, v...), fl.fields)
}

func (fl *fieldLogger) Println(v ...interface{}) {
	Output(fl.l, fl.GetLevel(), fmt.Sprintln(v...), fl.fields)
}

func (fl *fieldLogger) Log(level Level, v ...interface{}) {
	Output(fl.l, level, fmt.Sprint(v...), fl.fields)
}

func (fl *fieldLogger) Logf(level Level, format string, v ...interface{}) {
	Output(fl.l, level, fmt.Sprintf(format, v...), fl.fields)
}

func (fl *fieldLogger) Logln(level Level, v ...interface{}) {
	Output(fl.l, level, fmt.Sprintln(v...), fl.fields)
}

func (fl *fieldLogger) Debugf(format string, v ...interface{}) {
	if fl.l.IsEnable(DebugLevel) {
		Output(fl.l, DebugLevel, fmt.Sprintf(format, v...), fl.fields)
	}
}

func (fl *fieldLogger) Infof(format string, v ...interface{}) {
	if fl.l.IsEnable(InfoLevel) {
		Output(fl.l, InfoLevel, fmt.Sprintf(format, v...), fl.fields)
	}
}
