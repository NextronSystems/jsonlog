package thorlog

type DnsCacheEntry struct {
	LogObjectHeader
	Host string `json:"host" textlog:"entry"`
	IP   string `json:"ip" textlog:"ip"`
}

const typeDnsCacheEntry = "DNSCache entry"

func init() { AddLogObjectType(typeDnsCacheEntry, &DnsCacheEntry{}) }

func NewDnsCacheEntry(host string, ip string) *DnsCacheEntry {
	return &DnsCacheEntry{
		LogObjectHeader: LogObjectHeader{
			Type:    typeDnsCacheEntry,
			Summary: host,
		},
		Host: host,
		IP:   ip,
	}
}
