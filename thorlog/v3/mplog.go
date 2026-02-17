package thorlog

import (
	"time"

	"github.com/NextronSystems/jsonlog"
)

// DetectionAddEntry represents a detection event in the Microsoft Protection Log.
type DetectionAddEntry struct {
	jsonlog.ObjectHeader

	Time       time.Time    `json:"time" textlog:"time"`
	ThreatName string       `json:"threat_name" textlog:"threat_name"`
	Detected   KeyValueList `json:"detected" textlog:",expand"`
}

func (DetectionAddEntry) observed() {}

const typeDetectionAdd = "DetectionAdd MPLog entry"

func init() { AddLogObjectType(typeDetectionAdd, &DetectionAddEntry{}) }

func NewDetectionAddEntry(t time.Time, threat string, detected KeyValueList) *DetectionAddEntry {
	return &DetectionAddEntry{
		ObjectHeader: jsonlog.ObjectHeader{
			Type: "DETECTION_ADD MPLog entry",
		},
		Time:       t,
		ThreatName: threat,
		Detected:   detected,
	}
}

// EstimatedImpactEntry represents an event in the Microsoft Protection Log that lists the impact of a specific file on the monitoring of a process.
type EstimatedImpactEntry struct {
	jsonlog.ObjectHeader

	Time             time.Time `json:"time" textlog:"time"`
	ProcessImageName string    `json:"image" textlog:"image"`
	Pid              int       `json:"pid" textlog:"pid"`
	AccessedFile     string    `json:"file" textlog:"file"`
}

func (EstimatedImpactEntry) observed() {}

const typeEstimatedImpact = "EstimatedImpact MPLog entry"

func init() { AddLogObjectType(typeEstimatedImpact, &EstimatedImpactEntry{}) }

func NewEstimatedImpactEntry(t time.Time, image string, pid int, file string) *EstimatedImpactEntry {
	return &EstimatedImpactEntry{
		ObjectHeader: jsonlog.ObjectHeader{
			Type: typeEstimatedImpact,
		},
		Time:             t,
		ProcessImageName: image,
		Pid:              pid,
		AccessedFile:     file,
	}
}

// SdnQueryEntry represents an event in the Microsoft Protection Log that lists a query to the Smart Data Network.
type SdnQueryEntry struct {
	jsonlog.ObjectHeader

	Time     time.Time `json:"time" textlog:"time"`
	Filepath string    `json:"file" textlog:"file"`
	Sha1     string    `json:"sha1" textlog:"sha1"`
	Sha256   string    `json:"sha256" textlog:"sha256"`
}

func (SdnQueryEntry) observed() {}

const typeSdnQuery = "SDN query MPLog entry"

func init() { AddLogObjectType(typeSdnQuery, &SdnQueryEntry{}) }

func NewSdnQueryEntry(t time.Time, file string, sha1 string, sha256 string) *SdnQueryEntry {
	return &SdnQueryEntry{
		ObjectHeader: jsonlog.ObjectHeader{
			Type: typeSdnQuery,
		},
		Time:     t,
		Filepath: file,
		Sha1:     sha1,
		Sha256:   sha256,
	}
}

// EmsDetectionEntry represents an event in the Microsoft Protection Log that lists a detection on process behaviour.
type EmsDetectionEntry struct {
	jsonlog.ObjectHeader

	Time       time.Time `json:"time" textlog:"time"`
	ThreatName string    `json:"threat_name" textlog:"threat"`
	Pid        int       `json:"pid" textlog:"pid"`
}

func (EmsDetectionEntry) observed() {}

const typeEmsDetection = "EMS detection MPLog entry"

func init() { AddLogObjectType(typeEmsDetection, &EmsDetectionEntry{}) }

func NewEmsDetection(timestamp time.Time, threatName string, pid int) *EmsDetectionEntry {
	return &EmsDetectionEntry{
		ObjectHeader: jsonlog.ObjectHeader{
			Type: "EMS detection MPLog entry",
		},
		Time:       timestamp,
		ThreatName: threatName,
		Pid:        pid,
	}
}
