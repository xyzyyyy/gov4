package cmd

import (
	"fmt"
	x "mywabot/system"
	"sort"
	"strings"
)

type item struct {
	Name   []string
	Prefix bool
}

type tagSlice []string

func (t tagSlice) Len() int {
	return len(t)
}

func (t tagSlice) Less(i int, j int) bool {
	return t[i] < t[j]
}

func (t tagSlice) Swap(i int, j int) {
	t[i], t[j] = t[j], t[i]
}

func menu(client *x.Nc, m *x.IMsg) {
	var str string
	str += fmt.Sprintf("Hai, *%s* 👋 I am an automatic system (WhatsApp Bot) that can help to do things, search and get data or information only via WhatsApp.", m.PushName)
	var tags map[string][]item
	for _, list := range x.GetList() {
		if tags == nil {
			tags = make(map[string][]item)
		}
		if _, ok := tags[list.Tags]; !ok {
			tags[list.Tags] = []item{}
		}
		tags[list.Tags] = append(tags[list.Tags], item{Name: list.Cmd, Prefix: list.Prefix})
	}

	var keys tagSlice
	for key := range tags {
		keys = append(keys, key)
	}

	sort.Sort(keys)

	for _, key := range keys {
		str += fmt.Sprintf("*%s MENU*\n", strings.ToUpper(key))
		for _, e := range tags[key] {
			var prefix string
			if e.Prefix {
				prefix = m.Prefix[:1]
			} else {
				prefix = ""
			}
			for _, nm := range e.Name {
				str += fmt.Sprintf("▹ %s%s\n", prefix, nm)
			}
		}
		str += "\n"
	}
	txt := str + "\n© Whatsapp Bot Hosting 2023"
	m.Reply(txt)
}

func init() {
	x.NewCmd(&x.ICmd{
		Name:   "menu",
		Cmd:    []string{"menu"},
		Tags:   "main",
		Prefix: true,
		Exec:   menu,
	})
}
