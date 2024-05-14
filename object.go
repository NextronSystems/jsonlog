package jsonlog

// Object is the interface that all log objects must implement.
// Each log object has a `type` and `summary` field in its JSON representation.
// The type field is used to identify the object type, and the summary field is
// a human-readable summary of the object's contents.
type Object interface {
	// EmbeddedHeader returns the header of the log object.
	EmbeddedHeader() ObjectHeader
	// isObject is a marker method that ensures that all log objects must embed the ObjectHeader.
	isObject()
}

// ObjectHeader is the header of a log object. It must be embedded in all log objects.
type ObjectHeader struct {
	// Summary is a human-readable summary of the object's contents.
	Summary string `json:"summary"`
	// Type is the type of the object. It should be unique across all log objects
	// and can be used to identify the object type that has embedded this header.
	Type string `json:"type"`
}

func (l ObjectHeader) EmbeddedHeader() ObjectHeader { return l }

func (l ObjectHeader) isObject() {}
