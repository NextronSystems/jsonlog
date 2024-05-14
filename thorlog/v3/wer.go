package thorlog

import (
	"time"
)

type WerCrashReport struct {
	Type        string    `json:"-" textlog:"-"`
	Exe         string    `json:"exe" textlog:"exe"`
	Date        time.Time `json:"date" textlog:"date"`
	AppPath     string    `json:"apppath" textlog:"apppath"`
	Error       string    `json:"error" textlog:"error"`
	FaultModule string    `json:"fault_in_module" textlog:"fault_in_module"`
}

type AnalysisResult struct {
	Exe         string    `json:"exe"`
	Date        time.Time `json:"date"`
	AppPath     string    `json:"apppath"`
	Error       string    `json:"error"`
	FaultModule string    `json:"fault_in_module"`
}
