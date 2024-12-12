package thorlog

import (
	"strings"

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

	ActiveModules  []string `json:"active_modules"`
	ActiveFeatures []string `json:"active_features"`

	License LicenseInfo `json:"license" textlog:"license,expand"`

	FpFilters []string `json:"fp_filters"`
}

const typeScanInfo = "THOR invocation information"

func init() { AddLogObjectType(typeScanInfo, &ScanInfo{}) }

func NewScanInfo() *ScanInfo {
	return &ScanInfo{
		ObjectHeader: jsonlog.ObjectHeader{
			Type:    typeScanInfo,
			Summary: "Information about THOR invocation",
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
