package v2

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/NextronSystems/jsonlog/thorlog/common"
	"golang.org/x/mod/semver"
)

type Event struct {
	LogEventMetadata Metadata `textlog:",expand"`

	Data Fields `textlog:",expand"`

	EventVersion common.Version `json:"log_version"`
}

func (e *Event) Metadata() *common.LogEventMetadata {
	return (*common.LogEventMetadata)(&e.LogEventMetadata)
}

func (e *Event) Message() string {
	for _, kvPair := range e.Data {
		if kvPair.Key == "message" {
			return fmt.Sprint(kvPair.Value)
		}
	}
	return ""
}

func (e *Event) Version() common.Version {
	return common.Version(semver.Canonical(common.JsonV2))
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
		if kvPair.Key == "log_version" {
			var version common.Version
			switch v := kvPair.Value.(type) {
			case string:
				version = common.Version(v)
			case int:
				version = common.Version("v" + strconv.Itoa(v))
			case float64:
				version = common.Version("v" + strconv.Itoa(int(v)))
			default:
				return errors.New("invalid version type")
			}
			e.EventVersion = version
			continue
		} else if kvPair.Key == "time" || kvPair.Key == "level" || kvPair.Key == "module" || kvPair.Key == "scanid" || kvPair.Key == "uid" || kvPair.Key == "hostname" {
			continue
		}
		filteredData = append(filteredData, kvPair)
	}
	if e.EventVersion.Major() != common.JsonV2 {
		return fmt.Errorf("unsupported version %s", e.EventVersion)
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
	combinedData = append(combinedData, data2[1:len(data2)-1]...)
	combinedData = append(combinedData, fmt.Sprintf(`,"version":%q`, e.EventVersion)...)
	combinedData = append(combinedData, '}')
	return combinedData, nil
}

type Metadata struct {
	Time   time.Time       `json:"time" textlog:"-"`
	Lvl    common.LogLevel `json:"level" textlog:"-"`
	Mod    string          `json:"module" textlog:"module"`
	ScanID string          `json:"scanid" textlog:"scanid,omitempty"`
	GenID  string          `json:"uid" textlog:"uid,omitempty"`
	Source string          `json:"hostname" textlog:"-"`
}
