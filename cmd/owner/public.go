package cmd

import (
	"mywabot/config"
	x "mywabot/system"
)

func init(){
	x.NewCmd(&x.ICmd{
		Name: "public",
		Cmd: []string{"public"},
		Desc: "Set to public",
		Tags: "owner",
		Prefix: true,
		IsOwner: true,
		Exec: func(sock *x.Nc, m *x.IMsg){
			config.SetSelf(false)
			m.Reply("Update Public Mode!!")
		},
	})
}