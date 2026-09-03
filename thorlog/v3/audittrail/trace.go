package audittrail

import (
	"time"

	"github.com/NextronSystems/jsonlog/thorlog/common"
)

// Trace describes an entry in the audit trail log.
type Trace interface {
	Timestamps() map[string]time.Time
	Version() common.Version
}
