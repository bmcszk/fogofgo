package api

type Message struct {
	Type string `json:"type"`
}

type BasicMessage struct {
	Type    string `json:"type"`
	Payload string `json:"payload"`
}
