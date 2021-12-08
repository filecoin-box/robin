package slack

import (
	"encoding/json"
	"errors"
	"github.com/imroc/req"
	"golang.org/x/xerrors"
)

const (
	SlackUrl = "https://slack.com/api/chat.postMessage"
)

type message struct {
	Token   string `json:"token"`
	Channel string `json:"channel"`
	Text    string `json:"text"`
}

type response struct {
	Ok    bool   `json:"ok"`
	Error string `json:"error"`
}

func NewMessage(token, channel, text string) *message {
	return &message{
		Token:   token,
		Channel: channel,
		Text:    text,
	}
}

func (m *message) Send() error {
	if m.Token == "" {
		return errors.New("token does not exist")
	}
	if m.Channel == "" {
		return errors.New("channel does not exist")
	}
	if m.Text == "" {
		return errors.New("message does not exist")
	}

	msg, _ := json.Marshal(*m)
	params := &req.Param{}
	err := json.Unmarshal(msg, params)
	if err != nil {
		return err
	}

	resp, err := req.Post(SlackUrl, *params)
	if err != nil {
		return err
	}

	r := new(response)
	err = resp.ToJSON(r)
	if err != nil {
		return err
	}
	if !r.Ok {
		return xerrors.New(r.Error)
	}

	return nil
}
