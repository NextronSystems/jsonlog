package thorlog

import (
	"github.com/NextronSystems/jsonlog"
)

type LogLine struct {
	jsonlog.ObjectHeader

	LineIndex uint64 `json:"line_index" textlog:"-"`
	Line      string `json:"line" textlog:"line"`
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
