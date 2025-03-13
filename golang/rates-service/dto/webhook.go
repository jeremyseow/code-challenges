package dto

type Payload struct {
	Event string `json:"event"`
}

type WebHookRequest struct {
	URL  string  `json:"url"`
	ID   string  `json:"id"`
	Data Payload `json:"data"`
}
