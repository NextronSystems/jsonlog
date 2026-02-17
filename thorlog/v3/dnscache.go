package thorlog

type DnsCacheEntry struct {
	LogObjectHeader
	Host string `json:"host" textlog:"entry"`
	IP   string `json:"ip" textlog:"ip"`
}

func (DnsCacheEntry) observed() {}

const typeDnsCacheEntry = "DNS cache entry"

func init() { AddLogObjectType(typeDnsCacheEntry, &DnsCacheEntry{}) }

func NewDnsCacheEntry(host string, ip string) *DnsCacheEntry {
	return &DnsCacheEntry{
		LogObjectHeader: LogObjectHeader{
			Type: typeDnsCacheEntry,
		},
		Host: host,
		IP:   ip,
	}
}
