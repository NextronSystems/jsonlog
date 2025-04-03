package thorlog

import (
	"time"
)

type VirusTotalInformation struct {
	LookupResult     string             `json:"result" textlog:"result"`
	PositiveVerdicts int64              `json:"positive_verdicts" textlog:"verdicts"`
	TotalVerdicts    int64              `json:"total_verdicts"`
	History          *VirusTotalHistory `json:"history,omitempty" textlog:",omitempty,expand"`
}

type VirusTotalHistory struct {
	Names           StringList `json:"names,omitempty" textlog:"names" jsonschema:"nullable"`
	Tags            StringList `json:"tags,omitempty" textlog:"tags" jsonschema:"nullable"`
	Submissions     int64      `json:"submissions,omitempty"  textlog:"submissions"`
	FirstSubmission *Time      `json:"first_submission,omitempty" textlog:"first_submission,omitempty"`
	LastSubmission  *Time      `json:"last_submission,omitempty" textlog:"last_submission,omitempty"`
}
