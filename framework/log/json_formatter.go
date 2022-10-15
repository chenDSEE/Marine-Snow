package log

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

type JSONFormatter struct {
	bufPool sync.Pool
}

var _ formatter = &JSONFormatter{}

func newJSONFormatter(separator string) formatter {
	return &JSONFormatter{}
}

func (fm *JSONFormatter) Format(level Level, now time.Time, msg string, fields Fields) ([]byte, error) {
	// add default field
	if fields == nil {
		fields = make(Fields, 3)
	}

	fields["message"] = msg
	fields["level"] = level.String()
	fields["time"] = now.Format(time.RFC3339)

	buf := bytes.NewBuffer([]byte{})
	encoder := json.NewEncoder(buf)
	if err := encoder.Encode(fields); err != nil {
		return nil, fmt.Errorf("failed to marshal fields to JSON, %w", err)
	}

	return buf.Bytes(), nil
}
