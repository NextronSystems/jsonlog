package thorlog

type WindowsEvent struct {
	LogObjectHeader
	Event string `json:"event" textlog:"event"`
}

func (WindowsEvent) reportable() {}

const typeWindowsEvent = "event"

func init() { AddLogObjectType(typeWindowsEvent, &WindowsEvent{}) }

func NewWindowsEvent(event string) *WindowsEvent {
	return &WindowsEvent{
		LogObjectHeader: LogObjectHeader{
			Type: typeWindowsEvent,
		},
		Event: event,
	}
}

type WindowsMutex struct {
	LogObjectHeader

	Mutex string `json:"mutex" textlog:"mutex"`
}

func (WindowsMutex) reportable() {}

const typeWindowsMutex = "mutex"

func init() { AddLogObjectType(typeWindowsMutex, &WindowsMutex{}) }

func NewWindowsMutex(mutex string) *WindowsMutex {
	return &WindowsMutex{
		LogObjectHeader: LogObjectHeader{
			Type: typeWindowsMutex,
		},
		Mutex: mutex,
	}
}
