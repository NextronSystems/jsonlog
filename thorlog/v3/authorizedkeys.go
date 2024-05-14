package thorlog

import (
	"github.com/NextronSystems/jsonlog"
)

type AuthorizedKeysEntry struct {
	jsonlog.ObjectHeader

	Path    string `json:"path" textlog:"path"`
	Type    string `json:"type" textlog:"type"`
	Key     string `json:"key" textlog:"key"`
	Comment string `json:"comment" textlog:"comment"`
	Line    string `json:"line" textlog:"line"`
}

const typeAuthorizedKeysEntry = "authorized_keys entry"

func init() { AddLogObjectType(typeAuthorizedKeysEntry, &AuthorizedKeysEntry{}) }

func NewAuthorizedKeysEntry() *AuthorizedKeysEntry {
	return &AuthorizedKeysEntry{
		ObjectHeader: jsonlog.ObjectHeader{
			Type: typeAuthorizedKeysEntry,
		},
	}
}
