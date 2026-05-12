package ws

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/coder/websocket"
	"github.com/dskdev/crypto-hockey-game-engine/internal/protocol"
)

type Conn struct {
	c           *websocket.Conn
	slot        int
	userID      string
	telegramID  int64
	username    string
	closeOnce   sync.Once
	doneCh      chan struct{}
}

func newConn(c *websocket.Conn) *Conn {
	return &Conn{c: c, doneCh: make(chan struct{})}
}

func (c *Conn) Slot() int          { return c.slot }
func (c *Conn) UserID() string     { return c.userID }
func (c *Conn) Username() string   { return c.username }
func (c *Conn) TelegramID() int64  { return c.telegramID }

func (c *Conn) ReadJSON(ctx context.Context, v *protocol.ClientEnvelope) error {
	_, data, err := c.c.Read(ctx)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, v)
}

func (c *Conn) WriteJSON(ctx context.Context, v protocol.ServerEnvelope) error {
	buf, err := json.Marshal(v)
	if err != nil {
		return err
	}
	return c.c.Write(ctx, websocket.MessageText, buf)
}

func (c *Conn) Close() {
	c.closeOnce.Do(func() {
		_ = c.c.Close(websocket.StatusNormalClosure, "bye")
		close(c.doneCh)
	})
}

func (c *Conn) WaitClosed() { <-c.doneCh }

func (c *Conn) Done() <-chan struct{} { return c.doneCh }
