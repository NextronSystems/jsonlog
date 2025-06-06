package thorlog

import (
	"time"

	"github.com/NextronSystems/jsonlog"
	"github.com/google/uuid"
)

type JumplistEntry struct {
	jsonlog.ObjectHeader

	Path          string    `json:"path" textlog:"path"`
	Pinned        bool      `json:"pinned" textlog:"pinned"`
	LastAccess    time.Time `json:"last_access" textlog:"last_access"`
	AccessCount   int       `json:"access_count" textlog:"access_count"`
	NetbiosName   string    `json:"netbios_name" textlog:"netbios_name"`
	ObjectID      uuid.UUID `json:"object_id" textlog:"object_id"`
	VolumeID      uuid.UUID `json:"volume_id" textlog:"volume_id"`
	BirthVolumeID uuid.UUID `json:"birth_volume_id" textlog:"birth_volume_id"`
	EntryID       uint64    `json:"entry_id" textlog:"entry_id"`
	Checksum      uint64    `json:"checksum" textlog:"checksum"`
}

func (JumplistEntry) reportable() {}

const typeJumplistEntry = "jumplist entry"

func init() { AddLogObjectType(typeJumplistEntry, &JumplistEntry{}) }

func NewJumplistEntry(path string) *JumplistEntry {
	return &JumplistEntry{
		ObjectHeader: jsonlog.ObjectHeader{
			Type: typeJumplistEntry,
		},
		Path: path,
	}
}
