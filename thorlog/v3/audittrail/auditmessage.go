package audittrail

import (
	"time"

	"github.com/NextronSystems/jsonlog"
	"github.com/NextronSystems/jsonlog/thorlog/common"
	"github.com/NextronSystems/jsonlog/thorlog/v3"
)

// AuditMessage is a less verbose variant of THOR Message that is written to the audit trail.
// It is derived from a THOR Message, with its metadata flattened and reduced to the time,
// level and module.
type AuditMessage struct {
	jsonlog.ObjectHeader

	// Time is the time at which the message was logged.
	Time time.Time `json:"time"`
	// Lvl is the level at which the message was logged.
	Lvl thorlog.LogLevel `json:"level"`
	// Mod is the THOR module that logged the message.
	Mod string `json:"module"`
	// Text is the message that was logged.
	Text string `json:"message"`
	// Fields contains additional structured fields that were logged. These
	// contain details about the Text displayed.
	Fields     thorlog.MessageFields `json:"fields" jsonschema:"nullable"`
	LogVersion common.Version        `json:"log_version"`
}

const TypeAuditMessage = "THOR audit message"

func init() { thorlog.AddLogObjectType(TypeAuditMessage, &AuditMessage{}) }

func NewAuditMessage(message *thorlog.Message) *AuditMessage {
	msg := &AuditMessage{
		ObjectHeader: thorlog.LogObjectHeader{
			Type: TypeAuditMessage,
		},
		Time:       message.Meta.Time,
		Lvl:        message.Meta.Lvl,
		Mod:        message.Meta.Mod,
		Text:       message.Text,
		Fields:     message.Fields,
		LogVersion: message.LogVersion,
	}
	return msg
}

func (a *AuditMessage) Timestamps() map[string]time.Time {
	return map[string]time.Time{
		"PRINTED": a.Time,
	}
}

func (a *AuditMessage) Version() common.Version {
	return a.LogVersion
}
