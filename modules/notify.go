package modules

import (
	logging "github.com/ipfs/go-log/v2"
	"github.com/luluup777/robin/notify"
)

var robinLog = logging.Logger("robin")

func (m *Monitor) robin() {
	app := notify.NewNotify(&notify.Config{
		notify.Platform(m.config.Notify.Platform),
		m.config.Notify.Webhook,
	})

	for {
		select {
		case msg := <-m.notify:
			go func() {
				err := app.Send(msg)
				if err != nil {
					robinLog.Errorw("robin send msg", "err", err)
				}
			}()
		case <-m.notifyConfChange:
			newApp := notify.NewNotify(&notify.Config{
				notify.Platform(m.config.Notify.Platform),
				m.config.Notify.Webhook,
			})
			err := newApp.Send("test notify")
			if err != nil {
				robinLog.Warnw("notify test fail, no update", "platform", m.config.Notify.Platform, "webhook", m.config.Notify.Webhook, "err", err)
			} else {
				robinLog.Infow("update notify config", "platform", m.config.Notify.Platform, "webhook", m.config.Notify.Webhook)
				app = newApp
			}
		}
	}
}
