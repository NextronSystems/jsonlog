package thorlog

import (
	"fmt"
	"time"

	"github.com/NextronSystems/jsonlog"
	"github.com/NextronSystems/jsonlog/thorlog/truncate"
)

type RegistryValue struct {
	jsonlog.ObjectHeader

	File        string    `json:"file,omitempty" textlog:"file"`
	Key         string    `json:"key" textlog:"key"`
	Modified    time.Time `json:"modified" textlog:"modified"`
	ParsedValue string    `json:"value" textlog:"value"`
	Size        uint64    `json:"size" textlog:"size"`
}

func (RegistryValue) reportable() {}

const TypeRegistryValue = "registry value"

func init() { AddLogObjectType(TypeRegistryValue, &RegistryValue{}) }

func NewRegistryValue() *RegistryValue {
	return &RegistryValue{
		ObjectHeader: jsonlog.ObjectHeader{
			Type: TypeRegistryValue,
		},
	}
}

type RegistryKey struct {
	jsonlog.ObjectHeader
	File            string    `json:"path,omitempty" textlog:"path,omitempty"`
	Key             string    `json:"key" textlog:"key"`
	Modified        time.Time `json:"modified" textlog:"modified"`
	FormattedValues string    `json:"values" textlog:"values"`
}

func (RegistryKey) reportable() {}

func (s *RegistryKey) Truncate(matches []jsonlog.FieldMatch, truncateLimit int, stringContext int) jsonlog.Object {
	var ourMatches []truncate.Match
	for _, match := range matches {
		if match.FieldPointer == &s.FormattedValues {
			ourMatches = append(ourMatches, match.Match)
		}
	}
	var copiedKey = *s
	copiedKey.FormattedValues = truncate.TruncateWithNewlines(s.FormattedValues, ourMatches, truncateLimit, stringContext)
	return &copiedKey
}

func (s *RegistryKey) RawEvent() (string, *jsonlog.Reference) {
	return s.FormattedValues, jsonlog.NewReference(s, &s.FormattedValues)
}

const TypeRegistryKey = "registry key"

func init() { AddLogObjectType(TypeRegistryKey, &RegistryKey{}) }

func NewRegistryKey() *RegistryKey {
	return &RegistryKey{
		ObjectHeader: jsonlog.ObjectHeader{
			Type: TypeRegistryKey,
		},
	}
}

const TypeMsOfficeConnectionCacheEntry = "MS Office connection cache entry"

func init() {
	AddLogObjectType(TypeMsOfficeConnectionCacheEntry, &MsOfficeConnectionCacheEntry{})
}

type MsOfficeConnectionCacheEntry struct {
	jsonlog.ObjectHeader
	RegistryHive string    `json:"registry_hive" textlog:"path"`
	Entry        string    `json:"entry" textlog:"entry"`
	Modified     time.Time `json:"modified" textlog:"modified"`
	Key          string    `json:"key" textlog:"key"`
}

func (MsOfficeConnectionCacheEntry) reportable() {}

func NewMsOfficeConnectionCacheEntry() *MsOfficeConnectionCacheEntry {
	return &MsOfficeConnectionCacheEntry{
		ObjectHeader: jsonlog.ObjectHeader{
			Type: TypeMsOfficeConnectionCacheEntry,
		},
	}
}

type RegisteredDebugger struct {
	jsonlog.ObjectHeader
	Executable string `json:"executable" textlog:"file"`
	Debugger   string `json:"debugger" textlog:"element"`
}

func (RegisteredDebugger) reportable() {}

const TypeRegisteredDebugger = "registered debugger"

func init() { AddLogObjectType(TypeRegisteredDebugger, &RegisteredDebugger{}) }

func NewRegisteredDebugger(target string, debugger string) *RegisteredDebugger {
	return &RegisteredDebugger{
		ObjectHeader: jsonlog.ObjectHeader{
			Type:    TypeRegisteredDebugger,
			Summary: fmt.Sprintf("%q registered as debugger for %q", debugger, target),
		},
		Executable: target,
		Debugger:   debugger,
	}
}
