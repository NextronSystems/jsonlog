package thorlog

import (
	"github.com/NextronSystems/jsonlog"
)

type WmiElement struct {
	jsonlog.ObjectHeader

	Key        string `json:"key" textlog:"key"`
	Filtertype string `json:"filter_type" textlog:"filtertype"`

	Eventfiltername   string `json:"event_filter_name" textlog:"eventfiltername"`
	Eventconsumername string `json:"event_consumer_name" textlog:"eventconsumername"`
	Eventfilter       string `json:"event_filter" textlog:"eventfilter"`
	Eventconsumer     string `json:"event_consumer" textlog:"eventconsumer"`
}

func (WmiElement) observed() {}

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

func (WmiStartupCommand) observed() {}

const typeWmiStartupCommand = "WMI startup command"

func init() { AddLogObjectType(typeWmiStartupCommand, &WmiStartupCommand{}) }

func NewWmiStartupCommand() *WmiStartupCommand {
	return &WmiStartupCommand{
		ObjectHeader: jsonlog.ObjectHeader{
			Type: typeWmiStartupCommand,
		},
	}
}
