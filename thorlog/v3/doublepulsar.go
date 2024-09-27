package thorlog

import (
	"fmt"
)

type DoublePulsarHandshake struct {
	LogObjectHeader

	Type string    `json:"handshake_type" textlog:"type"` // SMB or RDP
	Key  HexNumber `json:"key,omitempty" textlog:"key,omitempty"`
}

const typeDoublePulsarHandshake = "DoublePulsar Handshake"

func init() { AddLogObjectType(typeDoublePulsarHandshake, &DoublePulsarHandshake{}) }

func NewDoublePulsarHandshake(handshakeType string, key uint64) *DoublePulsarHandshake {
	return &DoublePulsarHandshake{
		LogObjectHeader: LogObjectHeader{
			Type:    typeDoublePulsarHandshake,
			Summary: fmt.Sprintf("DoublePulsar Handshake via %s succeeded", handshakeType),
		},
		Key:  HexNumber(key),
		Type: handshakeType,
	}
}
