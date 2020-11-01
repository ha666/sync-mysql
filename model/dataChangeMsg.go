package model

type DataChangeMsg struct {
	Database string                 `json:"database"`
	Table    string                 `json:"table"`
	Type     string                 `json:"type"`
	Ts       int64                  `json:"ts"`
	Xid      int64                  `json:"xid"`
	Commit   bool                   `json:"commit"`
	Data     map[string]interface{} `json:"data"`
	Old      map[string]interface{} `json:"old"`
}
