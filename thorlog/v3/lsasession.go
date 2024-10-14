package thorlog

import (
	"time"
)

type LsaSession struct {
	LogObjectHeader

	LsaSession  string    `json:"lsa_session" textlog:"lsa_session"`
	User        string    `json:"user" textlog:"user"`
	AuthPackage string    `json:"auth_package" textlog:"auth_package"`
	Type        string    `json:"session_type" textlog:"type"`
	LogonTime   time.Time `json:"logon_time" textlog:"logon_time"`
	Domain      string    `json:"domain" textlog:"domain"`
	Server      string    `json:"server" textlog:"server"`
}

func (LsaSession) reportable() {}

const typeLsaSession = "lsa session"

func init() { AddLogObjectType(typeLsaSession, &LsaSession{}) }

func NewLsaSession() *LsaSession {
	return &LsaSession{
		LogObjectHeader: LogObjectHeader{
			Type: typeLsaSession,
		},
	}
}
