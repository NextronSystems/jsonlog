package thorlog

type BeaconConfig struct {
	Type             string `json:"type" textlog:"type"`
	C2               string `json:"c2" textlog:"c2"`
	Port             string `json:"port" textlog:"port"`
	SpawnTo          string `json:"spawnto" textlog:"spawnto"`
	InjectionProcess string `json:"injection_process" textlog:"injection_process"`
	Pipename         string `json:"pipename" textlog:"pipename"`
	UserAgent        string `json:"user_agent" textlog:"user_agent"`
	Proxy            string `json:"proxy" textlog:"proxy"`
}
