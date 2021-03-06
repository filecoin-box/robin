package slack

import (
	"bytes"
	"encoding/json"
	"golang.org/x/xerrors"
	"net/http"
)

type message struct {
	Text string `json:"text"`
}

func NewMessage(text string) *message {
	return &message{
		Text: text,
	}
}

func (m *message) Send(webhook string) error {
	msgByte, err := json.Marshal(m)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", webhook, bytes.NewBuffer(msgByte))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		return nil
	}

	return xerrors.Errorf("%s", resp.Status)
}
