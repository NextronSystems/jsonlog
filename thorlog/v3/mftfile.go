package thorlog

import (
	"time"

	"github.com/NextronSystems/jsonlog"
)

type MftFileEntry struct {
	jsonlog.ObjectHeader

	Path     string    `json:"path" textlog:"path"`
	Size     int64     `json:"size" textlog:"size"`
	Dir      bool      `json:"dir" textlog:"dir"`
	Modified time.Time `json:"modified" textlog:"modified"`
	Created  time.Time `json:"created" textlog:"created"`
	Accessed time.Time `json:"accessed" textlog:"accessed"`
	Changed  time.Time `json:"changed" textlog:"changed"`
	Filename string    `json:"filename" textlog:"filename"`
	Deleted  bool      `json:"deleted,omitempty" textlog:"deleted,omitempty"`
	Flags    *uint64   `json:"flags,omitempty" textlog:"flags,omitempty"`
}

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

	EventTime time.Time  `json:"event_time" textlog:"event_time"`
	Filename  string     `json:"filename" textlog:"filename"`
	Reasons   StringList `json:"reasons" textlog:"reason"`
}

const typeUsnEntry = "USN entry"

func init() { AddLogObjectType(typeUsnEntry, &UsnEntry{}) }

func NewUsnEntry() *UsnEntry {
	return &UsnEntry{
		LogObjectHeader: LogObjectHeader{
			Type: typeUsnEntry,
		},
	}
}
