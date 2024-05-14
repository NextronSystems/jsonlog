package thorlog

import (
	"github.com/NextronSystems/jsonlog"
)

type FirewallRule struct {
	jsonlog.ObjectHeader

	Path            string     `json:"path" textlog:"path"`
	LocalPorts      StringList `json:"lport" textlog:"lport"`
	RemotePorts     StringList `json:"rport" textlog:"rport"`
	LocalAddresses  StringList `json:"lip" textlog:"lip"`
	RemoteAddresses StringList `json:"rip" textlog:"rip"`
	Name            string     `json:"name" textlog:"name"`
	Allow           bool       `json:"allow" textlog:"allow"`
	Enabled         bool       `json:"enabled" textlog:"enabled"`
	Inbound         bool       `json:"inbound" textlog:"inbound"`
	Protocol        string     `json:"protocol" textlog:"protocol"`
}

const typeFirewallRule = "firewall rule"

func init() { AddLogObjectType(typeFirewallRule, &FirewallRule{}) }

func NewFirewallRule() *FirewallRule {
	return &FirewallRule{
		ObjectHeader: jsonlog.ObjectHeader{
			Type: typeFirewallRule,
		},
	}
}

type RawFirewallRule struct {
	jsonlog.ObjectHeader

	Rule string `json:"rule" textlog:"rule"`
}

const typeRawFirewallRule = "raw firewall rule"

func init() { AddLogObjectType(typeRawFirewallRule, &RawFirewallRule{}) }

func NewRawFirewallRule(rule string) *RawFirewallRule {
	return &RawFirewallRule{
		ObjectHeader: jsonlog.ObjectHeader{
			Type:    typeRawFirewallRule,
			Summary: rule,
		},
		Rule: rule,
	}
}
