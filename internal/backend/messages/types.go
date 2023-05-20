package messages

type BasicMessage struct {
	MsgType string `json:"msgtype"`
}

type ConnectMessage struct {
	MsgType  string `json:"msgtype"`
	User     string `json:"user"`
	Password string `json:"password"`
}

type GetMatchStatus struct {
	MsgType string `json:"msgtype"`
	MatchID string `json:"id"`
}

// type RespMatchStatus struct {
// 	MsgType   string    `json:"msgtype"`
// 	MatchData MatchData `json:"value"`
// }
