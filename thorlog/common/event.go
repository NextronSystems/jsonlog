package common

import (
	"time"
)

type LogLevel string

const (
	Error   LogLevel = "Error"
	Alert   LogLevel = "Alert"
	Warning LogLevel = "Warning"
	Notice  LogLevel = "Notice"
	Info    LogLevel = "Info"
	Debug   LogLevel = "Debug"
)

type LogEventMetadata struct {
	Time   time.Time `json:"time" textlog:"-"`
	Lvl    LogLevel  `json:"level" textlog:"-"`
	Mod    string    `json:"module" textlog:"-"`
	ScanID string    `json:"scan_id" textlog:"scanid,omitempty"`
	GenID  string    `json:"event_id" textlog:"uid,omitempty"`
	Source string    `json:"hostname" textlog:"-"`
}

// Event describes the basic information of a THOR event that is available in all versions.
// The actual type can differ depending on the version.
type Event interface {
	// Metadata returns the metadata of the log event. It is never nil and changes to the metadata are reflected in the event.
	Metadata() *LogEventMetadata
	Message() string
	Version() Version
}
