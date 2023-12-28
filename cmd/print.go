package cmd

import (
        "fmt"
        x "mywabot/system"

        "github.com/amiruldev20/waSocket"
        "github.com/fatih/color"
)

func init() {

        x.ListenerAdd(
                func(c *waSocket.Client, m *x.IMsg) {

                        if m.IsGroup {
                                get, err := c.GetGroupInfo(m.From)

                                if err != nil {
                                        fmt.Println(err)
                                        return
                                }
                                color.Yellow("\n--------------------------------------------")
                                color.Cyan("CHAT: %s", m.From.String())
                                color.Cyan("GC NAME: %s", get.Name)
                                color.Yellow("JID: %s", m.Sender.ToNonAD())
                                color.Yellow("NAME: %s", m.PushName)
                                color.Green("TYPE: %s", m.Type)
                                color.HiGreen("ID: %s", m.ID)
                                color.Cyan("Device ID: %d", m.Sender.Device)
                                color.Green("Message:\n%s", m.Text)
                                color.Yellow("--------------------------------------------")
                        } else {
                                magenta := color.New(color.FgGreen).SprintFunc()
                                color.Yellow("\n--------------------------------------------")
                                fmt.Printf("%s\n", magenta(fmt.Sprintf("JID: %s", m.Sender.String())))
                                fmt.Printf("%s\n", magenta(fmt.Sprintf("NAME: %s", m.PushName)))
                                fmt.Printf("%s\n", magenta(fmt.Sprintf("TYPE: %s", m.Type)))
                                fmt.Printf("%s\n", magenta(fmt.Sprintf("MESSAGE:\n%s", m.Text)))
                                color.Yellow("--------------------------------------------")
                        }
                },
        )
}
