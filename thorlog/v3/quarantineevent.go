package thorlog

import (
	"time"

	"github.com/NextronSystems/jsonlog"
)

type QuarantineEvent struct {
	jsonlog.ObjectHeader

	Id        string    `json:"id" textlog:"id"`
	Timestamp time.Time `json:"timestamp" textlog:"timestamp"`
	Name      string    `json:"name" textlog:"name"`
	Type      string    `json:"event_type" textlog:"type"`
	Url       string    `json:"url" textlog:"url,omitempty"`
}

func (QuarantineEvent) reportable() {}

const typeQuarantineEvent = "quarantine event"

func init() { AddLogObjectType(typeQuarantineEvent, &QuarantineEvent{}) }

func NewQuarantineEvent() *QuarantineEvent {
	return &QuarantineEvent{
		ObjectHeader: jsonlog.ObjectHeader{
			Type: typeQuarantineEvent,
		},
	}
}
