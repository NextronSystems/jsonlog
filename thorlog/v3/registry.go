package thorlog

import (
	"time"

	"github.com/NextronSystems/jsonlog"
)

type RegistryValue struct {
	jsonlog.ObjectHeader

	Key         string    `json:"key" textlog:"key"`
	Modified    time.Time `json:"modified" textlog:"modified"`
	ParsedValue string    `json:"value" textlog:"value,omitincontext"`
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
	Key             string    `json:"key" textlog:"key"`
	Modified        time.Time `json:"modified" textlog:"modified"`
	FormattedValues string    `json:"values" textlog:"values,omitincontext"`
}

func (RegistryKey) reportable() {}

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
	Entry    string    `json:"entry" textlog:"entry"`
	Modified time.Time `json:"modified" textlog:"modified"`
	Key      string    `json:"key" textlog:"key"`
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
			Type: TypeRegisteredDebugger,
		},
		Executable: target,
		Debugger:   debugger,
	}
}
