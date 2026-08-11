package thorlog

import (
	"time"

	"github.com/NextronSystems/jsonlog"
)

type AntiVirusProduct struct {
	LogObjectHeader

	Name            string `json:"name" textlog:"name"`
	Status          string `json:"status" textlog:"status"`
	SignatureStatus string `json:"signature_status" textlog:"signature_status"`
	Path            string `json:"path" textlog:"path"`
	// SignatureUpdated is the time when the product's signatures were last updated.
	// It is only available for products that report this information.
	SignatureUpdated time.Time `json:"signature_updated,omitzero" textlog:"signature_updated,omitempty"`
}

func (AntiVirusProduct) observed() {}

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
	// Modified is the last write time of the registry key that holds this exclusion.
	// Since all exclusions of the same type share that key, this is the time when any
	// exclusion of this type was last added or removed, not necessarily this one.
	Modified time.Time `json:"modified,omitzero" textlog:"modified,omitempty"`
}

func (AntiVirusExclude) observed() {}

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
