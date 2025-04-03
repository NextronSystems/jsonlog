package thorlog

import (
	"time"

	"github.com/NextronSystems/jsonlog"
)

type HotfixSummary struct {
	jsonlog.ObjectHeader
	LastHotfix Time `json:"last_hotfix" textlog:"last_hotfix"`
}

func (HotfixSummary) reportable() {}

const typeHotfixSummary = "hotfix summary"

func init() { AddLogObjectType(typeHotfixSummary, &HotfixSummary{}) }

func NewHotfixSummary(lastHotfix time.Time) *HotfixSummary {
	return &HotfixSummary{
		ObjectHeader: LogObjectHeader{
			Type:    typeHotfixSummary,
			Summary: "last hotfix installed " + lastHotfix.Format("2006-01-02"),
		},
		LastHotfix: lastHotfix,
	}
}

type EndOfLifeReport struct {
	jsonlog.ObjectHeader

	Version   string `json:"version" textlog:"version"`
	EndOfLife Time   `json:"end_of_life" textlog:"end_time"`
}

func (EndOfLifeReport) reportable() {}

const typeEndOfLifeReport = "end of life report"

func init() { AddLogObjectType(typeEndOfLifeReport, &EndOfLifeReport{}) }

func NewEndOfLifeReport(version string, endOfLife time.Time) *EndOfLifeReport {
	return &EndOfLifeReport{
		ObjectHeader: LogObjectHeader{
			Type:    typeEndOfLifeReport,
			Summary: "end of life for " + version + " was " + endOfLife.Format("2006-01-02"),
		},
		Version:   version,
		EndOfLife: endOfLife,
	}
}
