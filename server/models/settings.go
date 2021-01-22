package models

//Settings for the server
type Settings struct {
	DefaultReply             int    `json:"default_reply"`
	PartialMockServerAddress string `json:"partial_mock_server_address"`
}
