package structs

type Button struct {
	Action struct {
		Type    string `json:"type"`
		Payload string `json:"payload"`
		Label   string `json:"label"`
	} `json:"action"`
	Color string `json:"color"`
}

type Keyboard struct {
	OneTime bool       `json:"one_time"`
	Buttons [][]Button `json:"buttons"`
}

type Message struct {
	ID          int
	Date        int
	Out         int
	UserID      int   `json:"user_id"`
	ChatID      int   `json:"chat_id"`
	PeerID      int64 `json:"peer_id"`
	ReadState   int   `json:"read_state"`
	ConvMsgID   int   `json:"conversation_message_id"`
	Title       string
	Body        string
	Action      string
	ActionMID   int `json:"action_mid"`
	Flags       int
	Timestamp   int64
	Payload     string
	FwdMessages []Message `json:"fwd_messages"`
}

// Messages - VK Messages
type Messages struct {
	Count int
	Items []*Message
}

type Reply struct {
	Msg      string
	Keyboard *Keyboard
}
