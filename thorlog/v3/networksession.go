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
	// Established is the absolute time when the session was opened.
	Established time.Time `json:"established,omitzero" textlog:"established,omitempty"`
	NumOpens    int       `json:"num_opens" textlog:"num_opens"`
}

func (NetworkSession) observed() {}

const typeNetworkSession = "network session"

func init() { AddLogObjectType(typeNetworkSession, &NetworkSession{}) }

func NewNetworkSession() *NetworkSession {
	return &NetworkSession{
		LogObjectHeader: LogObjectHeader{
			Type: typeNetworkSession,
		},
	}
}
