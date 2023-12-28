package cmd

import (
        "fmt"
        x "mywabot/system"
        "net/url"
)

func init() {
        x.NewCmd(&x.ICmd{
                Name:   "ai",
                Cmd:    []string{"ai"},
                Tags:   "ai",
                IsQuery:  true,
                ValueQ: ".ai siapa kamu?",
                Exec: func(sock *x.Nc, m *x.IMsg) {
                        m.React("⏱️")
                        var res struct {
                                Answer string `json:"answer"`
                        }
                        err := x.GetResult("https://pro.amirull.dev/api/aiplus?apikey=adm&text="+url.QueryEscape(m.Query), &res)

                        if err != nil {
                                m.Reply(fmt.Sprint(err))
                                return
                        }
                        m.Reply(res.Answer)
                        m.React("✅")
                },
        })
}
