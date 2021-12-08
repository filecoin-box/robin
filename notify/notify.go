package notify

import (
	"github.com/luluup777/robin/notify/slack"
	"golang.org/x/xerrors"
)

type Platform string

const (
	Slack Platform = "slack"
)

type Notify struct {
	c *Config
}

type Config struct {
	Platform Platform
	Token    string
	Channel  string
}

func NewNotify(c *Config) *Notify {
	return &Notify{
		c: c,
	}
}

func (n *Notify) Send(msg string) error {
	switch n.c.Platform {
	case Slack:
		return slack.NewMessage(n.c.Token, n.c.Channel, msg).Send()
	default:
		return xerrors.New("not supported")
	}
}
