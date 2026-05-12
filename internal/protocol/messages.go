package protocol

const (
	TypeAuth       = "AUTH"
	TypeInput      = "INPUT"
	TypePing       = "PING"
	TypeAuthOK     = "AUTH_OK"
	TypeAuthFail   = "AUTH_FAIL"
	TypeMatchState = "MATCH_STATE"
	TypeSnapshot   = "SNAPSHOT"
	TypeGoal       = "GOAL"
	TypeMatchEnd   = "MATCH_END"
	TypePong       = "PONG"
)

type V2 struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

type Score struct {
	A int `json:"a"`
	B int `json:"b"`
}

type PuckSnap struct {
	X  float64 `json:"x"`
	Y  float64 `json:"y"`
	VX float64 `json:"vx"`
	VY float64 `json:"vy"`
}

type AuthPayload struct {
	InitData string `json:"initData"`
}

type InputPayload struct {
	Seq          int   `json:"seq"`
	TClient      int64 `json:"tClient"`
	MalletTarget V2    `json:"malletTarget"`
}

type PingPayload struct {
	TClient int64 `json:"tClient"`
}

type ClientEnvelope struct {
	Type  string        `json:"type"`
	Auth  *AuthPayload  `json:"auth,omitempty"`
	Input *InputPayload `json:"input,omitempty"`
	Ping  *PingPayload  `json:"ping,omitempty"`
}

type Opponent struct {
	Username   string `json:"username"`
	TelegramID int64  `json:"telegramId"`
}

type AuthOKPayload struct {
	PlayerSlot string   `json:"playerSlot"`
	Opponent   Opponent `json:"opponent"`
}

type AuthFailPayload struct {
	Reason string `json:"reason"`
}

type MatchStatePayload struct {
	Phase           string `json:"phase"`
	CountdownMs     *int   `json:"countdownMs,omitempty"`
	DurationLeftMs  *int   `json:"durationLeftMs,omitempty"`
}

type SnapshotPayload struct {
	TServer int64    `json:"tServer"`
	AckSeq  int      `json:"ackSeq"`
	MalletA V2       `json:"malletA"`
	MalletB V2       `json:"malletB"`
	Puck    PuckSnap `json:"puck"`
	Score   Score    `json:"score"`
}

type GoalPayload struct {
	Scorer string `json:"scorer"`
	Score  Score  `json:"score"`
}

type MatchEndPayload struct {
	WinnerUserID *string `json:"winnerUserId"`
	Reason       string  `json:"reason"`
	FinalScore   Score   `json:"finalScore"`
}

type PongPayload struct {
	TClient int64 `json:"tClient"`
	TServer int64 `json:"tServer"`
}

type ServerEnvelope struct {
	Type       string             `json:"type"`
	AuthOK     *AuthOKPayload     `json:"authOk,omitempty"`
	AuthFail   *AuthFailPayload   `json:"authFail,omitempty"`
	MatchState *MatchStatePayload `json:"matchState,omitempty"`
	Snapshot   *SnapshotPayload   `json:"snapshot,omitempty"`
	Goal       *GoalPayload       `json:"goal,omitempty"`
	MatchEnd   *MatchEndPayload   `json:"matchEnd,omitempty"`
	Pong       *PongPayload       `json:"pong,omitempty"`
}
