package thorlog

import (
	"time"

	"github.com/NextronSystems/jsonlog"
)

type LogLine struct {
	jsonlog.ObjectHeader

	LineIndex uint64 `json:"line_index" textlog:"-"`
	Line      string `json:"line" textlog:"line"`
	// Time contained in the log line, if the line has a timestamp that could be parsed.
	Time time.Time `json:"time,omitzero" textlog:"time,omitempty"`
}

func (LogLine) observed() {}

const TypeLogLine = "log line"

func init() { AddLogObjectType(TypeLogLine, &LogLine{}) }

func NewLogLine() *LogLine {
	return &LogLine{
		ObjectHeader: jsonlog.ObjectHeader{
			Type: TypeLogLine,
		},
	}
}
