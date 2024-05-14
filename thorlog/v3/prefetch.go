package thorlog

import (
	"github.com/NextronSystems/jsonlog"
)

type PrefetchElement struct {
	jsonlog.ObjectHeader

	Dir  string `json:"dir" textlog:"dir"`
	File *File  `json:"file" textlog:",expand"`
}

const typePrefetchElement = "prefetch element"

func init() { AddLogObjectType(typePrefetchElement, &PrefetchElement{}) }

func NewPrefetchElement() *PrefetchElement {
	return &PrefetchElement{
		ObjectHeader: jsonlog.ObjectHeader{
			Type: typePrefetchElement,
		},
	}
}

const typePrefetchFile = "prefetch file"

func init() { AddLogObjectType(typePrefetchFile, &PrefetchElement{}) }

func NewPrefetchFile() *PrefetchElement {
	return &PrefetchElement{
		ObjectHeader: jsonlog.ObjectHeader{
			Type: typePrefetchFile,
		},
	}
}
