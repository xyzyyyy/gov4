package config

import "sync"

var (
        mu    sync.Mutex
        Name  = "MyWa-GO ryhardev.my.id"
        Login = "code"
        Bot   = "447389672815"
        Owner = []string{"62882008211320", "447389672815"}
        Self  = false
        LolSite = "https://api.lolhuman.xyz/"
        LolKey = "5f38494f3555283d0446abdf"
)

func SetName(newName string) {
        mu.Lock()
        defer mu.Unlock()
        Name = newName
}

func SetSelf(new bool) {
        mu.Lock()
        defer mu.Unlock()
        Self = new
  }
