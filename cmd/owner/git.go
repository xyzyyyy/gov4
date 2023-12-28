package cmd

import (
	"fmt"
	x "mywabot/system"
	"os/exec"
	"strings"
)

func init() {
	x.NewCmd(&x.ICmd{
		Name:    "git",
		Cmd:     []string{"git"},
		Tags: "owner",
		Desc:    "Git Cli",
		Prefix:  false,
		IsOwner: true,
		Exec: func(sock *x.Nc, m *x.IMsg) {
			if strings.HasPrefix(m.Query, "up") {
				m.Reply("Update from github...")
				resp,
					err := exec.Command("bash", "-c", "git pull").Output()
				if err != nil {
					m.Reply(fmt.Sprintf("%v", err))
					return
				}
				m.Reply(string(resp))
			} else {
				ex := "git add . && git commit -m \"" + m.Query + "\" && git push"
				m.Reply("Mengunggah file ke GitHub...\nPesan commit: " + m.Query)
				cmd := exec.Command("bash", "-c", ex)
				output,
					err := cmd.Output()
				if err != nil {
					fmt.Println(err)
					m.Reply("Terjadi kesalahan saat melakukan git push")
					return
				}
				m.Reply(string(output))
			}
		},
	})
}
