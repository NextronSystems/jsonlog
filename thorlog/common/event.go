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

// LogEventMetadata contains the metadata of a log event.
// It is used to store common fields that are available in all log events.
//
// In a textlog formatted event, some of these fields are part of the header and do
// not occur "normally" as a KEY: VALUE pair in the event body.
// These fields are marked with the `textlog:"-"` tag to prevent them from being
// included in the event body.
type LogEventMetadata struct {
	Time   time.Time `json:"time" textlog:"-"`
	Lvl    LogLevel  `json:"level" textlog:"-"`
	Mod    string    `json:"module" textlog:"module"`
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
