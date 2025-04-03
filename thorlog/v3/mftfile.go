package thorlog

import (
	"time"

	"github.com/NextronSystems/jsonlog"
)

type MftFileEntry struct {
	jsonlog.ObjectHeader

	Path     string  `json:"path" textlog:"path"`
	Size     int64   `json:"size" textlog:"size"`
	Dir      bool    `json:"dir" textlog:"dir"`
	Modified Time    `json:"modified" textlog:"modified"`
	Created  Time    `json:"created" textlog:"created"`
	Accessed Time    `json:"accessed" textlog:"accessed"`
	Changed  Time    `json:"changed" textlog:"changed"`
	Filename string  `json:"filename" textlog:"filename"`
	Deleted  bool    `json:"deleted,omitempty" textlog:"deleted,omitempty"`
	Flags    *uint64 `json:"flags,omitempty" textlog:"flags,omitempty"`
}

func (MftFileEntry) reportable() {}

const typeMftFileEntry = "MFT entry"

func init() { AddLogObjectType(typeMftFileEntry, &MftFileEntry{}) }

func NewMftFileEntry() *MftFileEntry {
	return &MftFileEntry{
		ObjectHeader: jsonlog.ObjectHeader{
			Type: typeMftFileEntry,
		},
	}
}

type UsnEntry struct {
	LogObjectHeader

	EventTime Time       `json:"event_time" textlog:"event_time"`
	Filename  string     `json:"filename" textlog:"filename"`
	Reasons   StringList `json:"reasons" textlog:"reason"`
}

func (UsnEntry) reportable() {}

const typeUsnEntry = "USN entry"

func init() { AddLogObjectType(typeUsnEntry, &UsnEntry{}) }

func NewUsnEntry() *UsnEntry {
	return &UsnEntry{
		LogObjectHeader: LogObjectHeader{
			Type: typeUsnEntry,
		},
	}
}
