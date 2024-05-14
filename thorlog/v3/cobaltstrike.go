package thorlog

import (
	"github.com/NextronSystems/jsonlog"
)

type BeaconConfig struct {
	jsonlog.ObjectHeader

	Type             string `json:"type" textlog:"beacon_type"`
	C2               string `json:"c2" textlog:"beacon_c2"`
	Port             string `json:"port" textlog:"beacon_port"`
	SpawnTo          string `json:"spawnto" textlog:"beacon_spawnto"`
	InjectionProcess string `json:"injection_process" textlog:"beacon_injection_process"`
	Pipename         string `json:"pipename" textlog:"beacon_pipename"`
	UserAgent        string `json:"user_agent" textlog:"beacon_user_agent"`
	Proxy            string `json:"proxy" textlog:"beacon_proxy"`
}

const typeBeaconConfig = "CobaltStrike Beacon configuration"

func init() { AddLogObjectType(typeBeaconConfig, &BeaconConfig{}) }

func NewBeaconConfig() *BeaconConfig {
	return &BeaconConfig{
		ObjectHeader: jsonlog.ObjectHeader{
			Type: typeBeaconConfig,
		},
	}
}
