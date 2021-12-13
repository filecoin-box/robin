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
			app = notify.NewNotify(&notify.Config{
				notify.Platform(m.config.Notify.Platform),
				m.config.Notify.Webhook,
			})
		}
	}

}
