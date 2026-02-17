package thorlog

import (
	"github.com/NextronSystems/jsonlog"
)

type WindowsPipe struct {
	LogObjectHeader

	Pipe string `json:"pipe" textlog:"pipe"`
}

func (WindowsPipe) observed() {}

const typeWindowsPipe = "named pipe"

func init() { AddLogObjectType(typeWindowsPipe, &WindowsPipe{}) }

func NewWindowsPipe(pipe string) *WindowsPipe {
	return &WindowsPipe{
		LogObjectHeader: LogObjectHeader{
			Type: typeWindowsPipe,
		},
		Pipe: pipe,
	}
}

type WindowsPipeList struct {
	jsonlog.ObjectHeader
	Pipes StringList `json:"pipes" textlog:"pipes"`
}

func (WindowsPipeList) observed() {}

const typeWindowsPipeList = "pipe list"

func init() { AddLogObjectType(typeWindowsPipeList, &WindowsPipeList{}) }

func NewWindowsPipeList() *WindowsPipeList {
	return &WindowsPipeList{
		ObjectHeader: jsonlog.ObjectHeader{
			Type: typeWindowsPipeList,
		},
	}
}
