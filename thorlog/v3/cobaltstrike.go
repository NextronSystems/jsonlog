package thorlog

type BeaconConfig struct {
	Type             string `json:"type" textlog:"type"`
	C2               string `json:"c2" textlog:"c2"`
	Port             string `json:"port" textlog:"port"`
	SpawnTo          string `json:"spawn_to" textlog:"spawnto"`
	InjectionProcess string `json:"injection_process" textlog:"injection_process"`
	Pipename         string `json:"pipe_name" textlog:"pipename"`
	UserAgent        string `json:"user_agent" textlog:"user_agent"`
	Proxy            string `json:"proxy" textlog:"proxy"`

	// FullConfig is the full configuration of the beacon.
	// For now, it is filled with strings only until we refactor the parsing module.
	FullConfig map[string]any `json:"full_config" textlog:"-"`

	// CipherParameters contains information about how the beacon is hidden in the file.
	CipherParameters CipherParameters `json:"cipher_parameters" textlog:"cipher_parameters,expand,omitempty"`
}

type CipherParameters struct {
	XafEncoded        bool       `json:"xaf_encoded" textlog:"xaf_encoded"`
	XafEncodingAnchor int64      `json:"xaf_encoding_anchor" textlog:"xaf_encoding_anchor,omitempty"`
	XorKey            byte       `json:"xor_key" textlog:"xor_key"`
	BeaconOffset      uint64     `json:"beacon_offset" textlog:"beacon_offset"`
	BeaconLength      uint64     `json:"beacon_length" textlog:"beacon_length"`
	BlockStart        FirstBytes `json:"block_start" textlog:"block_start"`
	PairwiseSwapped   bool       `json:"pairwise_swapped" textlog:"pairwise_swapped"`
}
