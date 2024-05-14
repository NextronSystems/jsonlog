package parser

import (
	"encoding/json"
	"errors"

	"github.com/NextronSystems/jsonlog/thorlog/common"
	thorlogv1 "github.com/NextronSystems/jsonlog/thorlog/v1"
	thorlogv2 "github.com/NextronSystems/jsonlog/thorlog/v2"
	"github.com/NextronSystems/jsonlog/thorlog/v3"
)

func ParseEvent(data []byte) (common.Event, error) {
	var versionedEvent struct {
		Version common.Version `json:"log_version"`
		Type    string         `json:"type"`
	}
	if err := json.Unmarshal(data, &versionedEvent); err != nil {
		return nil, err
	}
	switch versionedEvent.Version.Major() {
	case common.JsonV1, "":
		var event thorlogv1.Event
		if err := json.Unmarshal(data, &event); err != nil {
			return nil, err
		}
		return &event, nil
	case common.JsonV2:
		var event thorlogv2.Event
		if err := json.Unmarshal(data, &event); err != nil {
			return nil, err
		}
		return &event, nil
	case common.JsonV3:
		var logObject thorlog.EmbeddedObject
		if err := json.Unmarshal(data, &logObject); err != nil {
			return nil, err
		}
		event, isEvent := logObject.Object.(common.Event)
		if !isEvent {
			return nil, errors.New("json v3 log object is not an event")
		}
		return event, nil
	}
	return nil, errors.New("unknown JSON version")
}
