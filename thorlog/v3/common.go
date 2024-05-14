package thorlog

import "github.com/NextronSystems/jsonlog/thorlog/common"

type (
	LogEventMetadata = common.LogEventMetadata
	LogLevel         = common.LogLevel
	Event            = common.Event
	Version          = common.Version
)

const (
	JsonV1 = common.JsonV1
	JsonV2 = common.JsonV2
	JsonV3 = common.JsonV3

	Error   = common.Error
	Alert   = common.Alert
	Warning = common.Warning
	Notice  = common.Notice
	Info    = common.Info
	Debug   = common.Debug
)

const currentVersion = "v3.0.0"
