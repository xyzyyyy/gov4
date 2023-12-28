package system

import (
	"fmt"
	"regexp"
	"strings"
)

var lists []ICmd

func NewCmd(cmd *ICmd) {
	lists = append(lists, *cmd)
}

func GetList() []ICmd {
	return lists
}

func Get(c *Nc, m *IMsg) {
    
    // Executing listeners
    for _, listen := range ListenerStore {
        listen(c.WA, m)
    }
    
    
	var prefix string
	pattern := regexp.MustCompile(`[?!.#]`)
	for _, f := range pattern.FindAllString(m.Prefix, -1) {
		prefix = f
	}
	for _, cmd := range lists {
		if cmd.After != nil {
			cmd.After(c, m)
		}
		re := regexp.MustCompile(`^` + cmd.Name + `$`)
		if reg := len(re.FindAllString(strings.ReplaceAll(m.Prefix, prefix, ""), -1)) > 0; reg {
			var cmdWithPref bool
			var cmdWithoutPref bool
			if cmd.Prefix && (prefix != "" && strings.HasPrefix(m.Prefix, prefix)) {
				cmdWithPref = true
			} else {
				cmdWithPref = false
			}

			if !cmd.Prefix {
				cmdWithoutPref = true
			} else {
				cmdWithoutPref = false
			}

			if !cmdWithPref && !cmdWithoutPref {
				continue
			}

			//Checking
			if cmd.IsOwner && !m.IsOwner {
				m.Reply("Fitur ini hanya untuk owner!!")
				continue
			}

			if cmd.IsMedia && m.Media == nil {
				m.Reply("Silahkan reply / input media!!")
				continue
			}

			if cmd.IsQuery && m.Query == "" {
				txt := fmt.Sprintf("Silahkan masukan query\n%s", cmd.ValueQ)
				m.Reply(txt)
				continue
			}

			if cmd.IsGroup && !m.IsGroup {
				m.Reply("Fitur ini hanya dapat digunakan didalam grup!!")
				continue
			}

			if (m.IsGroup && cmd.IsAdmin) && !m.IsAdmin {
				m.Reply("Fitur ini hanya untuk admin grup!!")
				continue
			}

			if m.IsBotAdmin && cmd.IsBotAdm {
				m.Reply("Untuk menggunakan fitur ini, bot harus menjadi admin!!")
				continue
			}

			if cmd.IsWait {
				m.Reply("Permintaan sedang diproses..")
			}

			cmd.Exec(c, m)
		}
	}
}
