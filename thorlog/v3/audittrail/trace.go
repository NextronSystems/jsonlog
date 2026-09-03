package audittrail

import (
	"time"

	"github.com/NextronSystems/jsonlog/thorlog/common"
)

// Trace describes an entry in the audit trail log.
// Both the objects that THOR observed during a scan (AuditRecord) and the messages
// that it printed while doing so (AuditMessage) are traces.
type Trace interface {
	Timestamps() map[string]time.Time
	Version() common.Version
}
