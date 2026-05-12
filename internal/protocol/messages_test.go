package protocol

import (
	"encoding/json"
	"testing"
)

func TestAuthRoundTrip(t *testing.T) {
	in := ClientEnvelope{Type: TypeAuth, Auth: &AuthPayload{InitData: "hello"}}
	b, err := json.Marshal(in)
	if err != nil {
		t.Fatal(err)
	}
	var out ClientEnvelope
	if err := json.Unmarshal(b, &out); err != nil {
		t.Fatal(err)
	}
	if out.Type != TypeAuth || out.Auth == nil || out.Auth.InitData != "hello" {
		t.Fatalf("round trip: %+v", out)
	}
}

func TestSnapshotMarshalsExpectedShape(t *testing.T) {
	s := ServerEnvelope{Type: TypeSnapshot, Snapshot: &SnapshotPayload{
		TServer: 1, AckSeq: 7,
		MalletA: V2{X: 100, Y: 200}, MalletB: V2{X: 700, Y: 200},
		Puck:    PuckSnap{X: 400, Y: 200, VX: 10, VY: 0},
		Score:   Score{A: 0, B: 0},
	}}
	b, err := json.Marshal(s)
	if err != nil {
		t.Fatal(err)
	}
	if !contains(string(b), `"type":"SNAPSHOT"`) || !contains(string(b), `"ackSeq":7`) {
		t.Fatalf("shape: %s", b)
	}
}

func contains(s, sub string) bool { return len(s) >= len(sub) && indexOf(s, sub) >= 0 }
func indexOf(s, sub string) int {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return i
		}
	}
	return -1
}
