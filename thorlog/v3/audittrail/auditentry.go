package audittrail

import (
	"time"

	"github.com/NextronSystems/jsonlog/thorlog/common"
)

// AuditEntry describes an entry in the audit trail log.
type AuditEntry interface {
	Timestamps() map[string]time.Time
	Version() common.Version
}
