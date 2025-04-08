package thorlog

import (
	"strconv"
	"time"

	"github.com/NextronSystems/jsonlog"
)

type WindowsEventlogEntry struct {
	jsonlog.ObjectHeader

	File string `json:"file,omitempty" textlog:"file,omitempty"`

	EventId       uint16    `json:"-" textlog:"event_id"`
	EventLevel    int       `json:"-" textlog:"event_level"`
	EventTime     time.Time `json:"-" textlog:"event_time"`
	EventChannel  string    `json:"-" textlog:"event_channel,omitempty"`
	EventComputer string    `json:"-" textlog:"event_computer,omitempty"`

	Entry KeyValueList `json:"entry" textlog:"entry"`
}

func (WindowsEventlogEntry) reportable() {}

const TypeEventlogEntry = "eventlog entry"

func init() { AddLogObjectType(TypeEventlogEntry, &WindowsEventlogEntry{}) }

func NewEventlogEntry() *WindowsEventlogEntry {
	return &WindowsEventlogEntry{
		ObjectHeader: jsonlog.ObjectHeader{
			Type: TypeEventlogEntry,
		},
	}
}

type EventlogProcessStart struct {
	jsonlog.ObjectHeader
	Process    string      `json:"process" textlog:"process"`
	StartTimes []time.Time `json:"start_times" textlog:"-"`
	Count      int         `json:"-" textlog:"count"`
}

func (EventlogProcessStart) reportable() {}

const TypeProcessStart = "process start"

func init() { AddLogObjectType(TypeProcessStart, &EventlogProcessStart{}) }

func NewEventlogProcessStart(process string, startTimes []time.Time) *EventlogProcessStart {
	return &EventlogProcessStart{
		ObjectHeader: jsonlog.ObjectHeader{
			Type:    TypeProcessStart,
			Summary: process + " started " + strconv.Itoa(len(startTimes)) + " times",
		},
		Process:    process,
		StartTimes: startTimes,
		Count:      len(startTimes),
	}
}
