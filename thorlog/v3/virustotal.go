package thorlog

import (
	"time"

	"github.com/NextronSystems/jsonlog"
)

type VirusTotalResult struct {
	jsonlog.ObjectHeader
	LookupResult     string             `json:"result" textlog:"virustotal_result"`
	PositiveVerdicts int64              `json:"positive_verdicts" textlog:"virustotal_verdicts"`
	TotalVerdicts    int64              `json:"total_verdicts"`
	History          *VirusTotalHistory `json:"history,omitempty" textlog:",omitempty,expand"`
}

const typeVirusTotalResult = "VirusTotal information"

func init() { AddLogObjectType(typeVirusTotalResult, &VirusTotalResult{}) }

func NewVirusTotalResult(result string) *VirusTotalResult {
	return &VirusTotalResult{
		ObjectHeader: jsonlog.ObjectHeader{
			Summary: "VirusTotal result: " + result,
			Type:    typeVirusTotalResult,
		},
		LookupResult: result,
	}
}

type VirusTotalHistory struct {
	Names           StringList `json:"names,omitempty" textlog:"virustotal_names"`
	Tags            StringList `json:"tags,omitempty" textlog:"virustotal_tags"`
	Submissions     int64      `json:"submissions,omitempty"  textlog:"virustotal_submissions"`
	FirstSubmission *time.Time `json:"first_submission,omitempty" textlog:"virustotal_first_submission,omitempty"`
	LastSubmission  *time.Time `json:"last_submission,omitempty" textlog:"virustotal_last_submission,omitempty"`
}
