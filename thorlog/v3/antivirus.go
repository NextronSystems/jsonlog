package thorlog

import (
	"github.com/NextronSystems/jsonlog"
)

type AntiVirusProduct struct {
	LogObjectHeader

	Name            string `json:"Name" textlog:"name"`
	Status          string `json:"Status" textlog:"status"`
	SignatureStatus string `json:"SignatureStatus" textlog:"signature_status"`
	Path            string `json:"Path" textlog:"path"`
}

const typeAntiVirusProduct = "Antivirus product"

func init() { AddLogObjectType(typeAntiVirusProduct, AntiVirusProduct{}) }

func NewAntiVirusProduct(name string) *AntiVirusProduct {
	return &AntiVirusProduct{
		LogObjectHeader: jsonlog.ObjectHeader{
			Type:    typeAntiVirusProduct,
			Summary: name,
		},
		Name: name,
	}
}

type AntiVirusExclude struct {
	LogObjectHeader

	Type      string `json:"type" textlog:"type"`
	Exclusion string `json:"exclusion" textlog:"exclusion"`
}

const typeAntiVirusExclude = "Antivirus exclusion"

func init() { AddLogObjectType(typeAntiVirusExclude, AntiVirusExclude{}) }

func NewAntiVirusExclude(exclusionType string, exclusion string) *AntiVirusExclude {
	return &AntiVirusExclude{
		LogObjectHeader: jsonlog.ObjectHeader{
			Type:    typeAntiVirusExclude,
			Summary: exclusionType + " " + exclusion,
		},
		Type:      exclusionType,
		Exclusion: exclusion,
	}
}
