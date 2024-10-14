package thorlog

import (
	"github.com/NextronSystems/jsonlog"
)

type PeSieveReport struct {
	jsonlog.ObjectHeader

	SuspiciousSections int `json:"suspicious_sections" textlog:"suspicious_sections"`
	Replaced           int `json:"replaced" textlog:"replaced"`
	HdrMod             int `json:"hdr_mod" textlog:"hdr_mod"`
	UnreachableFile    int `json:"unreachable_file" textlog:"unreachable_file"`
	Patched            int `json:"patched" textlog:"patched"`
	IatHooked          int `json:"iat_hooked" textlog:"iat_hooked"`
	Implanted          int `json:"implanted" textlog:"implanted"`
	Other              int `json:"other" textlog:"other"`
	Skipped            int `json:"skipped" textlog:"skipped"`
	Errors             int `json:"errors" textlog:"errors"`
}

func (PeSieveReport) reportable() {}

const typePeSieveReport = "PE-Sieve report"

func init() { AddLogObjectType(typePeSieveReport, &PeSieveReport{}) }

func NewPeSieveReport() *PeSieveReport {
	return &PeSieveReport{
		ObjectHeader: jsonlog.ObjectHeader{
			Type: typePeSieveReport,
		},
	}
}
