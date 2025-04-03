package thorlog

import (
	"time"

	"github.com/NextronSystems/jsonlog"
)

type ShimCacheEntry struct {
	jsonlog.ObjectHeader

	Timestamp Time   `json:"timestamp" textlog:"timestamp"`
	ExecFlag  *bool  `json:"exec_flag" textlog:"exec_flag,omitempty"`
	Path      string `json:"path" textlog:"path"`
	Hive      string `json:"hive" textlog:"hive"`
}

func (ShimCacheEntry) reportable() {}

const typeShimCacheEntry = "SHIM cache entry"

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

	Path             string `json:"path" textlog:"path"`
	Hive             string `json:"hive" textlog:"hive"`
	Entries          int    `json:"entries" textlog:"entries"`
	LastKnownEntries int    `json:"last_known_entries" textlog:"previous_entries,omitempty"`
}

func (ShimCache) reportable() {}

const typeShimCache = "SHIM cache"

func init() { AddLogObjectType(typeShimCache, &ShimCache{}) }

func NewShimCache() *ShimCache {
	return &ShimCache{
		ObjectHeader: jsonlog.ObjectHeader{
			Type: typeShimCache,
		},
	}
}
