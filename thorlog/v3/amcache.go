package thorlog

import (
	"time"

	"github.com/NextronSystems/jsonlog"
)

type AmcacheEntry struct {
	jsonlog.ObjectHeader

	File     *File  `json:"file" textlog:"file,expand"`
	SHA1     string `json:"sha1" textlog:"sha1"`
	Size     int64  `json:"size" textlog:"size"`
	Desc     string `json:"desc" textlog:"desc"`
	FirstRun Time   `json:"first_run" textlog:"first_run"`
	Created  Time   `json:"created" textlog:"created"`
	Product  string `json:"product" textlog:"product"`
	Company  string `json:"company" textlog:"company"`
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

func (AmcacheEntry) reportable() {}
