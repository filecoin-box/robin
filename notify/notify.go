package notify

import (
	"github.com/luluup777/robin/notify/lark"
	"github.com/luluup777/robin/notify/slack"
	"golang.org/x/xerrors"
)

type Platform string

const (
	Slack Platform = "slack"
	Lark  Platform = "lark"
)

type Notify struct {
	c *Config
}

type Config struct {
	Platform Platform
	Webhook  string
}

func NewNotify(c *Config) *Notify {
	return &Notify{
		c: c,
	}
}

func (n *Notify) Send(msg string) error {
	switch n.c.Platform {
	case Slack:
		return slack.NewMessage(msg).Send(n.c.Webhook)
	case Lark:
		return lark.NewMessage("text", msg).Send(n.c.Webhook)
	default:
		return xerrors.New("not supported")
	}
}
