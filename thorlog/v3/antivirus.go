package thorlog

import (
	"github.com/NextronSystems/jsonlog"
)

type AntiVirusProduct struct {
	LogObjectHeader

	Name            string `json:"name" textlog:"name"`
	Status          string `json:"status" textlog:"status"`
	SignatureStatus string `json:"signature_status" textlog:"signature_status"`
	Path            string `json:"path" textlog:"path"`
}

func (AntiVirusProduct) reportable() {}

const typeAntiVirusProduct = "antivirus product"

func init() { AddLogObjectType(typeAntiVirusProduct, &AntiVirusProduct{}) }

func NewAntiVirusProduct(name string) *AntiVirusProduct {
	return &AntiVirusProduct{
		LogObjectHeader: jsonlog.ObjectHeader{
			Type: typeAntiVirusProduct,
		},
		Name: name,
	}
}

type AntiVirusExclude struct {
	LogObjectHeader

	Type      string `json:"exclusion_type" textlog:"type"`
	Exclusion string `json:"exclusion" textlog:"exclusion"`
}

func (AntiVirusExclude) reportable() {}

const typeAntiVirusExclude = "antivirus exclusion"

func init() { AddLogObjectType(typeAntiVirusExclude, &AntiVirusExclude{}) }

func NewAntiVirusExclude(exclusionType string, exclusion string) *AntiVirusExclude {
	return &AntiVirusExclude{
		LogObjectHeader: jsonlog.ObjectHeader{
			Type: typeAntiVirusExclude,
		},
		Type:      exclusionType,
		Exclusion: exclusion,
	}
}
