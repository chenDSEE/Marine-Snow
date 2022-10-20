package log

type Level int

const (
	DebugLevel Level = iota
	InfoLevel
	WarnLevel
	ErrorLevel
	FatalLevel
	PanicLevel
)

func (level Level) String() string {
	return string(ParseLevel(level))
}

func ParseLevel(level Level) []byte {
	switch level {
	case DebugLevel:
		return []byte("debug")
	case InfoLevel:
		return []byte("info")
	case WarnLevel:
		return []byte("warning")
	case ErrorLevel:
		return []byte("error")
	case FatalLevel:
		return []byte("fatal")
	case PanicLevel:
		return []byte("panic")
	default:
		return []byte("unknown")
	}
}

// TODO: MarineSnow as a framework, should export a method for user to register self-define logger via Logger interface,
// for some app may want web framework can out put just same as app self-define log
// adapt for Go std log module
type Logger interface {
	GetLevel() Level
	SetLevel(Level)
	SetFormat(FormatType)
	WithFields(Fields) Logger // TODO: this should be another interface, like FieldLogger

	Log(Level, ...interface{})
	Logf(Level, string, ...interface{})
	Logln(Level, ...interface{})

	Print(v ...interface{})
	Printf(string, ...interface{})
	Println(v ...interface{})

	// TODO:
	//Debug(args ...interface{})
	Debugf(string, ...interface{})
	//Debugln(args ...interface{})
	//
	//Info(args ...interface{})
	Infof(string, ...interface{})
	//Infoln(args ...interface{})
	//
	//Warn(args ...interface{})
	//Warnf(format string, args ...interface{})
	//Warnln(args ...interface{})
	//
	//Error(args ...interface{})
	//Errorf(format string, args ...interface{})
	//Errorln(args ...interface{})
	//
	//Fatal(v ...interface{})
	//Fatalf(format string, v ...interface{})
	//Fatalln(v ...interface{})
	//
	//Panic(v ...interface{})
	//Panicf(format string, v ...interface{})
	//Panicln(v ...interface{})

	// TODO: like Logger.trace() should a method for logger struct, but not a method in this Logger interface
	// a Logger interface for a framework is for user to replace framework logger with self-define logger
	// you should not use this interface to limit and make sure logger method
}
