package cmd

import (
	"mywabot/config"
	x "mywabot/system"
)

func init(){
	x.NewCmd(&x.ICmd{
		Name: "self",
		Cmd: []string{"self"},
		Desc: "Set to self",
		Tags: "owner",
		Prefix: true,
		IsOwner: true,
		Exec: func(sock *x.Nc, m *x.IMsg){
			config.SetSelf(true)
			m.Reply("Update Self Mode!!")
		},
	})
}