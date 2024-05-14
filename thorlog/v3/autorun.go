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
}

const typeAutorunEntry = "autorun entry"

func init() { AddLogObjectType(typeAutorunEntry, &AutorunEntry{}) }

func NewAutorunEntry() *AutorunEntry {
	return &AutorunEntry{
		ObjectHeader: jsonlog.ObjectHeader{
			Type: typeAutorunEntry,
		},
	}
}

type Md5Change struct {
	jsonlog.ObjectHeader
	OldMd5     string `json:"old_md5" textlog:"md5_before"`
	CurrentMd5 string `json:"current_md5" textlog:"-"`
}

const typeMd5Change = "MD5 hash change"

func init() { AddLogObjectType(typeMd5Change, &Md5Change{}) }

func NewMd5Change(oldMd5, newMd5 string) *Md5Change {
	return &Md5Change{
		ObjectHeader: jsonlog.ObjectHeader{
			Type:    typeAutorunEntry,
			Summary: "Hash changed from " + oldMd5 + " to " + newMd5,
		},
		OldMd5:     oldMd5,
		CurrentMd5: newMd5,
	}
}
