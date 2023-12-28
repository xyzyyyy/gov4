package cmd

import (
	"context"
	x "mywabot/system"

	waProto "github.com/amiruldev20/waSocket/binary/proto"
	"github.com/amiruldev20/waSocket/types"
	"google.golang.org/protobuf/proto"
)

func init(){
x.NewCmd(&x.ICmd{
Name: "upsw",
Cmd: []string{"upsw"},
Tags: "owner",
IsOwner: true,
Exec: func(sock *x.Nc, m *x.IMsg){

if m.Msg.GetExtendedTextMessage() != nil {
// sw text
hx := "fbfcfa"
hxx := "3d6e75"
white, _ := x.HextoUint32(hx)

bg, _ := x.HextoUint32(hxx)
sock.WA.SendMessage(context.Background(), types.StatusBroadcastJID, &waProto.Message{
ExtendedTextMessage: &waProto.ExtendedTextMessage{
Text:           proto.String(m.Query),
TextArgb:       proto.Uint32(white),                    BackgroundArgb: proto.Uint32(bg),
}, 
})
m.React("✅")

} else {
m.Reply("Input text or reply media")
}
},
})
}