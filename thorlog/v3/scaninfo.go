package thorlog

import (
	"fmt"
	"strings"
	"time"

	"github.com/NextronSystems/jsonlog"
)

type SpaceSeparatedList []string

func (s SpaceSeparatedList) String() string {
	return strings.Join(s, " ")
}

type ScanInfo struct {
	jsonlog.ObjectHeader

	Versions  VersionInfo        `json:"versions" textlog:",expand"`
	Arguments SpaceSeparatedList `json:"arguments" textlog:"arguments"`
	ScanID    string             `json:"scan_id" textlog:"scan_id"`
	ThorDir   string             `json:"thor_dir" textlog:"thor_dir"`
	User      string             `json:"user" textlog:"user"`
	Elevated  bool               `json:"elevated" textlog:"elevated"`

	Outputs []ScannerOutput `json:"outputs"`

	ActiveModules  StringList `json:"active_modules" textlog:"active_modules"`
	ActiveFeatures StringList `json:"active_features" textlog:"active_features"`

	Threads int `json:"threads" textlog:"threads"`

	Timeout         time.Duration `json:"timeout" textlog:"timeout"`
	CPULimit        int           `json:"cpu_limit" textlog:"cpu_limit"`
	FreeMemoryLimit Memory        `json:"free_memory_limit" textlog:"free_memory_limit"`
	FileSizeLimit   Memory        `json:"file_size_limit" textlog:"file_size_limit"`

	License LicenseInfo `json:"license" textlog:"license,expand"`

	FpFilters []string `json:"fp_filters"`
}

const typeScanInfo = "THOR invocation information"

func init() { AddLogObjectType(typeScanInfo, &ScanInfo{}) }

func NewScanInfo() *ScanInfo {
	return &ScanInfo{
		ObjectHeader: jsonlog.ObjectHeader{
			Type: typeScanInfo,
		},
	}
}

type ScannerOutput struct {
	Kind   string `json:"kind"`
	Output string `json:"output"`
}

type LicenseInfo struct {
	Owner   string `json:"owner" textlog:"owner"`
	Type    string `json:"license_type" textlog:"type"`
	Starts  string `json:"starts" textlog:"starts"`
	Expires string `json:"expires" textlog:"expires"`
	Scanner string `json:"scanner" textlog:"scanner"`
	Hash    string `json:"hash" textlog:"hash"`
}

type VersionInfo struct {
	Thor       string `json:"thor" textlog:"version"`
	Build      string `json:"build" textlog:"build"`
	Signatures string `json:"signatures" textlog:"signature_version"`
	Sigma      string `json:"sigma_rules" textlog:"sigma_version"`
}

type Memory uint64

func (m Memory) String() string {
	var usedSuffix = "B"
	var divisor uint64 = 1
	bytes := uint64(m)
	for suffix, multiplier := range multipliers {
		if multiplier < bytes && multiplier > divisor {
			divisor = multiplier
			usedSuffix = suffix + "B"
		}
	}
	return fmt.Sprintf("%d%s", int64(float64(bytes)/float64(divisor)), usedSuffix)
}

const (
	kb = 1024
	mb = 1024 * kb
	gb = 1024 * mb
	tb = 1024 * gb
	pb = 1024 * tb
)

var multipliers = map[string]uint64{
	"K": kb,
	"M": mb,
	"G": gb,
	"T": tb,
	"P": pb,
}
