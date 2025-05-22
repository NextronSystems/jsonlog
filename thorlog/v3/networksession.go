package thorlog

import (
	"time"
)

type NetworkSession struct {
	LogObjectHeader
	Client     string        `json:"client" textlog:"client"`
	Username   string        `json:"user_name" textlog:"username"`
	ClientType string        `json:"client_type" textlog:"client_type"`
	Active     time.Duration `json:"active" textlog:"active"`
	Idle       time.Duration `json:"idle" textlog:"idle"`
	NumOpens   int           `json:"num_opens" textlog:"num_opens"`
}

func (NetworkSession) reportable() {}

const typeNetworkSession = "network session"

func init() { AddLogObjectType(typeNetworkSession, &NetworkSession{}) }

func NewNetworkSession() *NetworkSession {
	return &NetworkSession{
		LogObjectHeader: LogObjectHeader{
			Type: typeNetworkSession,
		},
	}
}
