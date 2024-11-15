package thorlog

type BeaconConfig struct {
	Type             string `json:"beacon_type" textlog:"beacon_type"`
	C2               string `json:"c2" textlog:"beacon_c2"`
	Port             string `json:"port" textlog:"beacon_port"`
	SpawnTo          string `json:"spawnto" textlog:"beacon_spawnto"`
	InjectionProcess string `json:"injection_process" textlog:"beacon_injection_process"`
	Pipename         string `json:"pipename" textlog:"beacon_pipename"`
	UserAgent        string `json:"user_agent" textlog:"beacon_user_agent"`
	Proxy            string `json:"proxy" textlog:"beacon_proxy"`
}
