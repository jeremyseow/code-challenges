package event

type ClientRequestAny struct {
	DeviceModel string      `json:"device_model"`
	OsType      string      `json:"os_type"`
	Events      []*EventAny `json:"events"`
}

type EventAny struct {
	EventName      string         `json:"event_name"`
	EventTimestamp uint64         `json:"event_timestamp"`
	ParamsAny      map[string]any `json:"params_any"`
}
