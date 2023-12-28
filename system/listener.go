package system

import (
    "github.com/amiruldev20/waSocket"
)

var ListenerStore = []func(*waSocket.Client, *IMsg){}

func ListenerAdd(f func(*waSocket.Client, *IMsg)) {
	ListenerStore = append(ListenerStore, f)
}