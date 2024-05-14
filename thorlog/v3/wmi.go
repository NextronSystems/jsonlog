package thorlog

import (
	"github.com/NextronSystems/jsonlog"
)

type WmiElement struct {
	jsonlog.ObjectHeader

	Key        string `json:"key" textlog:"key"`
	Filtertype string `json:"filtertype" textlog:"filtertype"`

	Eventfiltername   string `json:"eventfiltername" textlog:"eventfiltername"`
	Eventconsumername string `json:"eventconsumername" textlog:"eventconsumername"`
	Eventfilter       string `json:"eventfilter" textlog:"eventfilter"`
	Eventconsumer     string `json:"eventconsumer" textlog:"eventconsumer"`
}

const typeWmiElement = "WMI element"

func init() { AddLogObjectType(typeWmiElement, &WmiElement{}) }

func NewWmiElement() *WmiElement {
	return &WmiElement{
		ObjectHeader: jsonlog.ObjectHeader{
			Type: typeWmiElement,
		},
	}
}

type WmiStartupCommand struct {
	jsonlog.ObjectHeader
	Location string `json:"location" textlog:"location"`
	Caption  string `json:"caption" textlog:"caption"`
	Command  string `json:"command" textlog:"command"`
}

const typeWmiStartupCommand = "WMI startup command"

func init() { AddLogObjectType(typeWmiStartupCommand, &WmiStartupCommand{}) }

func NewWmiStartupCommand() *WmiStartupCommand {
	return &WmiStartupCommand{
		ObjectHeader: jsonlog.ObjectHeader{
			Type: typeWmiStartupCommand,
		},
	}
}
