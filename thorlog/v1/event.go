package v1

import (
	"encoding/json"

	"github.com/NextronSystems/jsonlog/thorlog/common"
	"golang.org/x/mod/semver"
)

type Event struct {
	common.LogEventMetadata

	Data Fields `textlog:",expand"`
}

func (e *Event) Metadata() *common.LogEventMetadata {
	return &e.LogEventMetadata
}

func (e *Event) Message() string {
	for _, kvPair := range e.Data {
		if kvPair.Key == "message" {
			return kvPair.Value
		}
	}
	return ""
}

func (e *Event) Version() common.Version {
	return common.Version(semver.Canonical(common.JsonV1))
}

func (e *Event) UnmarshalJSON(data []byte) error {
	if err := json.Unmarshal(data, &e.LogEventMetadata); err != nil {
		return err
	}
	if err := json.Unmarshal(data, &e.Data); err != nil {
		return err
	}
	// Remove the fields that are already in the metadata
	var filteredData Fields
	for _, kvPair := range e.Data {
		if kvPair.Key != "time" && kvPair.Key != "level" && kvPair.Key != "module" && kvPair.Key != "scanid" && kvPair.Key != "uid" && kvPair.Key != "hostname" && kvPair.Key != "log_version" {
			filteredData = append(filteredData, kvPair)
		}
	}
	e.Data = filteredData
	return nil
}

func (e Event) MarshalJSON() ([]byte, error) {
	data, err := json.Marshal(e.LogEventMetadata)
	if err != nil {
		return nil, err
	}
	// Prevent duplicate fields (which could happen if the same field is in both e.LogEventMetadata and e.Data)
	var filteredData Fields
	for _, kvPair := range e.Data {
		if kvPair.Key != "time" && kvPair.Key != "level" && kvPair.Key != "module" && kvPair.Key != "scanid" && kvPair.Key != "uid" && kvPair.Key != "hostname" && kvPair.Key != "log_version" {
			filteredData = append(filteredData, kvPair)
		}
	}
	data2, err := json.Marshal(filteredData)
	if err != nil {
		return nil, err
	}
	if data[len(data)-1] != '}' {
		panic("metadata was not marshalled correctly")
	}
	if data2[0] != '{' {
		panic("data was not marshalled correctly")
	}
	var combinedData = data[:len(data)-1]
	combinedData = append(combinedData, ',')
	combinedData = append(combinedData, data2[1:]...)
	return combinedData, nil
}
