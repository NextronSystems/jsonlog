package thorlog

import "github.com/NextronSystems/jsonlog"

// Issue describes a problem that occurred during the analysis of a scan target like a file or process.
// Often this will be an issue with displaying the results,
// e.g. the results may be truncated due to size limitations.
type Issue struct {
	// Affected is the path to the substructure that is related to the issue.
	// If the issue can't be related to a specific substructure, this may be null.
	Affected *jsonlog.Reference `json:"affected" textlog:"affected"`
	// Category is a human-readable description of the issue category.
	Category string `json:"category" textlog:"category"`
	// Description is a human-readable description of the issue.
	Description string `json:"description" textlog:"description"`
}

const (
	// IssueCategoryTruncated indicates that a value was truncated due to its size.
	IssueCategoryTruncated = "truncated"
	// IssueCategoryOutOfRange indicates that a value can't be represented in the format that the log uses.
	IssueCategoryOutOfRange = "out_of_range"
)
