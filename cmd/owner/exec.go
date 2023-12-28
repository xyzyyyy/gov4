package cmd

import (
	"fmt"
	x "mywabot/system"
	"os/exec"
)

func init() {
	x.NewCmd(&x.ICmd{
		Name:    `\$`,
		Cmd:     []string{"$"},
		Tags:    "owner",
		Prefix:  false,
		IsOwner: true,
		Exec: func(client *x.Nc, m *x.IMsg) {
			out, err := exec.Command("bash", "-c", m.Query).Output()
			if err != nil {
				m.Reply(fmt.Sprintf("%v", err))
				return
			}
			m.Reply(string(out))
		},
	})
}
