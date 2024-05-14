package thorlog

type HostsFileEntry struct {
	LogObjectHeader
	Host string `json:"host" textlog:"host"`
	IP   string `json:"ip" textlog:"ip"`
}

const typeHostsFileEntry = "hosts file entry"

func init() { AddLogObjectType(typeHostsFileEntry, &HostsFileEntry{}) }

func NewHostsFileEntry(host string, ip string) *HostsFileEntry {
	return &HostsFileEntry{
		LogObjectHeader: LogObjectHeader{
			Type:    typeHostsFileEntry,
			Summary: host,
		},
		Host: host,
		IP:   ip,
	}
}
