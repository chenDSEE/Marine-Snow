package log

import "time"

type formatter interface {
	Format(Level, time.Time, string, Fields) ([]byte, error)
}

type FormatType int

const (
	TextFormat FormatType = iota
	JSONFormat
)

func newFormatter(ft FormatType) formatter {
	switch ft {
	case TextFormat:
		return newTextFormatter(" ")
	case JSONFormat:
		return newJSONFormatter(" ")
	default:
		return nil
	}
}
