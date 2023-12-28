package cmd

import (
	x "mywabot/system"
)

func init() {
	x.NewCmd(&x.ICmd{
		Name:   "sc",
		Cmd:    []string{"sc"},
		Tags:   "main",
		Prefix: true,
		Exec: func(sock *x.Nc, m *x.IMsg) {
			m.Reply(`*waSocket Bot Info*

Script ini dibuat menggunakan bahasa GoLang

script ini di deploy pada hosting
dan dimanage melalui website.

*Thanks To:*
Vnia (Referensi serialize)
Justshorsuop (Plugins system)

Source Code:
https://github.com/amiruldev20/wabothost

Library:
https://github.com/amiruldev20/waSocket`)
		},
	})
}
