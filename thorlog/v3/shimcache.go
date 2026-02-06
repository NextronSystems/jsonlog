package thorlog

import (
	"time"

	"github.com/NextronSystems/jsonlog"
)

type ShimCacheEntry struct {
	jsonlog.ObjectHeader

	Timestamp time.Time `json:"timestamp" textlog:"timestamp"`
	ExecFlag  *bool     `json:"exec_flag" textlog:"exec_flag,omitempty"`
	Path      string    `json:"path" textlog:"path"`
}

func (ShimCacheEntry) observed() {}

const typeShimCacheEntry = "shim cache entry"

func init() { AddLogObjectType(typeShimCacheEntry, &ShimCacheEntry{}) }

func NewShimCacheEntry() *ShimCacheEntry {
	return &ShimCacheEntry{
		ObjectHeader: jsonlog.ObjectHeader{
			Type: typeShimCacheEntry,
		},
	}
}

type ShimCache struct {
	jsonlog.ObjectHeader

	Entries          int `json:"entries" textlog:"entries"`
	LastKnownEntries int `json:"last_known_entries" textlog:"previous_entries,omitempty"`
}

func (ShimCache) observed() {}

const typeShimCache = "shim cache"

func init() { AddLogObjectType(typeShimCache, &ShimCache{}) }

func NewShimCache() *ShimCache {
	return &ShimCache{
		ObjectHeader: jsonlog.ObjectHeader{
			Type: typeShimCache,
		},
	}
}
