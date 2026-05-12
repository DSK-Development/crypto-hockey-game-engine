package account

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Client struct {
	base    string
	token   string
	http    *http.Client
	retries int
}

func New(base, token string) *Client {
	return &Client{base: base, token: token, http: &http.Client{Timeout: 5 * time.Second}, retries: 3}
}

type AuthResult struct {
	AccessToken string `json:"accessToken"`
	UserID      string `json:"userId"`
	TelegramID  int64  `json:"telegramId"`
	Username    string `json:"username"`
}

func (c *Client) AuthTelegram(ctx context.Context, initData string) (AuthResult, error) {
	var out AuthResult
	err := c.do(ctx, http.MethodPost, "/auth/telegram",
		map[string]string{"initData": initData}, &out, false)
	return out, err
}

type Participant struct {
	UserID    string `json:"userId"`
	Placement int    `json:"placement"`
}

type MatchResult struct {
	MatchID      string        `json:"matchId"`
	Stake        int64         `json:"stakeStars"`
	Participants []Participant `json:"participants"`
	Reason       string        `json:"reason"`
}

func (c *Client) PostMatchResult(ctx context.Context, r MatchResult) error {
	return c.do(ctx, http.MethodPost, "/internal/match-result", r, nil, true)
}

func (c *Client) do(ctx context.Context, method, path string, body, out any, retry bool) error {
	var lastErr error
	attempts := 1
	if retry {
		attempts = c.retries
	}
	for i := 0; i < attempts; i++ {
		err := c.once(ctx, method, path, body, out)
		if err == nil {
			return nil
		}
		lastErr = err
		if !retry {
			return err
		}
		// backoff: 100ms, 200ms, 400ms
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(time.Duration(100*(1<<i)) * time.Millisecond):
		}
	}
	return lastErr
}

func (c *Client) once(ctx context.Context, method, path string, body, out any) error {
	var reqBody io.Reader
	if body != nil {
		buf, err := json.Marshal(body)
		if err != nil {
			return err
		}
		reqBody = bytes.NewReader(buf)
	}
	req, err := http.NewRequestWithContext(ctx, method, c.base+path, reqBody)
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
	if res.StatusCode >= 500 {
		return fmt.Errorf("account %s %s: %d", method, path, res.StatusCode)
	}
	if res.StatusCode >= 400 {
		buf, _ := io.ReadAll(res.Body)
		return fmt.Errorf("account %s %s: %d: %s", method, path, res.StatusCode, buf)
	}
	if out != nil && res.StatusCode != 204 {
		return json.NewDecoder(res.Body).Decode(out)
	}
	return nil
}
