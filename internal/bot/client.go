package bot

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Client struct {
	base  string
	token string
	http  *http.Client
}

func New(base, token string) *Client {
	return &Client{base: base, token: token, http: &http.Client{Timeout: 3 * time.Second}}
}

type MatchResult struct {
	WinnerUserID *string `json:"winnerUserId"`
	Reason       string  `json:"reason"`
	ScoreA       int     `json:"scoreA"`
	ScoreB       int     `json:"scoreB"`
}

func (c *Client) NotifyResult(ctx context.Context, matchID string, r MatchResult) error {
	buf, err := json.Marshal(r)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("%s/matches/%s/result", c.base, matchID), bytes.NewReader(buf))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Service-Token", c.token)
	res, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode >= 400 {
		return fmt.Errorf("bot notify result: %d", res.StatusCode)
	}
	return nil
}
