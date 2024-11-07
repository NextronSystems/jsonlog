package thorlog

import (
	"github.com/NextronSystems/jsonlog"
)

type AutorunEntry struct {
	jsonlog.ObjectHeader

	Type         string `json:"autorun_type" textlog:"autorun_type"`
	Location     string `json:"location" textlog:"location"`
	Image        *File  `json:"image" textlog:",expand"`
	Arguments    string `json:"arguments" textlog:"arguments"`
	Entry        string `json:"entry" textlog:"entry"`
	LaunchString string `json:"launch_string" textlog:"launch_string"`

	OldMd5 string `json:"old_md5,omitempty" textlog:"md5_before,omitempty"`
}

func (AutorunEntry) reportable() {}

const typeAutorunEntry = "autorun entry"

func init() { AddLogObjectType(typeAutorunEntry, &AutorunEntry{}) }

func NewAutorunEntry() *AutorunEntry {
	return &AutorunEntry{
		ObjectHeader: jsonlog.ObjectHeader{
			Type: typeAutorunEntry,
		},
	}
}
