package thorlog

import (
	"time"

	"github.com/NextronSystems/jsonlog"
)

type AmcacheEntry struct {
	jsonlog.ObjectHeader

	Path     string    `json:"path" textlog:"path"`
	SHA1     string    `json:"sha1" textlog:"sha1"`
	Size     int64     `json:"size" textlog:"size"`
	Desc     string    `json:"desc" textlog:"desc"`
	FirstRun time.Time `json:"first_run" textlog:"first_run"`
	Created  time.Time `json:"created" textlog:"created"`
	Product  string    `json:"product" textlog:"product"`
	Company  string    `json:"company" textlog:"company"`
}

const typeAmcacheEntry = "Amcache Entry"

func init() { AddLogObjectType(typeAmcacheEntry, &AmcacheEntry{}) }

func NewAmcacheEntry() *AmcacheEntry {
	return &AmcacheEntry{
		ObjectHeader: jsonlog.ObjectHeader{
			Type: typeAmcacheEntry,
		},
	}
}
