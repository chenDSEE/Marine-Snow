package log

import (
	"bytes"
	"sync"
	"time"
)

type textFormatter struct {
	separator []byte
	bufPool   sync.Pool
}

var _ formatter = &textFormatter{}

func newTextFormatter(separator string) formatter {
	return &textFormatter{
		separator: []byte(separator),
	}
}

func (fm *textFormatter) Format(level Level, now time.Time, msg string, fields Fields) ([]byte, error) {
	buf := bytes.NewBuffer([]byte{})

	buf.WriteString(now.Format(time.RFC3339))
	buf.Write(fm.separator)

	buf.WriteByte('[')
	buf.Write(ParseLevel(level))
	buf.WriteByte(']')
	buf.Write(fm.separator)

	buf.WriteString(msg)
	buf.Write(fm.separator)

	for k, v := range fields {
		buf.WriteByte('[')
		buf.WriteString(k)
		buf.WriteByte(':')
		buf.WriteString(v)
		buf.WriteByte(']')
		buf.Write(fm.separator)
	}

	buf.WriteByte('\n')
	return buf.Bytes(), nil
}
