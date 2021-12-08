package modules

import "github.com/luluup777/robin/parse"

func (m *Monitor) watchConfig() {
	var confCh = make(chan parse.Config, 1)
	go parse.WatchConf(confCh)

	for {
		select {
		case c := <-confCh:
			oldConf := m.config
			m.config = c
			if oldConf.M != c.M {
				m.monitorConfChange <- struct{}{}
			}
			if oldConf.N != c.N {
				m.notifyConfChange <- struct{}{}
			}
		}
	}
}
