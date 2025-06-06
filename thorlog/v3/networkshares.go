package thorlog

import (
	"github.com/NextronSystems/jsonlog"
)

type NetworkShare struct {
	jsonlog.ObjectHeader
	Name        string     `json:"name" textlog:"share_name"`
	Path        string     `json:"path" textlog:"path"`
	Permissions AclEntries `json:"permissions" textlog:"share_perms,omitempty"`
}

func (NetworkShare) reportable() {}

const typeNetworkShare = "network share"

func init() { AddLogObjectType(typeNetworkShare, &NetworkShare{}) }

func NewNetworkShare(name, path string) *NetworkShare {
	return &NetworkShare{
		ObjectHeader: jsonlog.ObjectHeader{
			Type: typeNetworkShare,
		},
		Name: name,
		Path: path,
	}
}
