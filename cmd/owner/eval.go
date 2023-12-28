package cmd

import (
	"encoding/json"
	"mywabot/config"
	x "mywabot/system"

	"github.com/robertkrimen/otto"
)

func init() {
	x.NewCmd(&x.ICmd{
		Name:    `=>`,
		Cmd:     []string{"=>"},
		Tags:    "owner",
		Prefix:  false,
		IsOwner: true,
		Exec: func(client *x.Nc, m *x.IMsg) {
			vm := otto.New()
			vm.Set("Name", config.Name)
			vm.Set("M", m)

			h, err := vm.Run(m.Query)
			if err != nil {
				m.Reply(err.Error())
				return
			}

			if h.IsObject() {
				var data interface{}
				h, _ := vm.Run("JSON.stringify(" + m.Query + ")")
				json.Unmarshal([]byte(h.String()), &data)
				pe, _ := json.MarshalIndent(data, "", "  ")
				m.Reply(string(pe))
			} else {
				m.Reply(h.String())
			}
		},
	})
}
