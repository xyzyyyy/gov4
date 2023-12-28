package system

import (
	"context"
	"mywabot/config"
	"strings"

	"github.com/amiruldev20/waSocket"
	waProto "github.com/amiruldev20/waSocket/binary/proto"
	"github.com/amiruldev20/waSocket/types/events"
	"google.golang.org/protobuf/proto"
)

func NewSmsg(mess *events.Message, sock *Nc) *IMsg {

	var command string
	var media waSocket.DownloadableMessage
	var owner []string
	var isOwner = false
	botNum, _ := sock.ParseJID(sock.WA.Store.ID.User)
	quotedMsg := mess.Message.GetExtendedTextMessage().GetContextInfo().GetQuotedMessage()

	var owners []string
	var owns []string
	owners = config.Owner
	for _, owner := range owners {
		waNumber := owner + "@s.whatsapp.net"
		owns = append(owns, waNumber)
	}
	owner = append(owns, botNum.String())

	for _, own := range owner {
		if own == mess.Info.Sender.ToNonAD().String() {
			isOwner = true
		}
	}

	if pe := mess.Message.GetExtendedTextMessage().GetText(); pe != "" {
		command = pe
	} else if pe := mess.Message.GetImageMessage().GetCaption(); pe != "" {
		command = pe
	} else if pe := mess.Message.GetVideoMessage().GetCaption(); pe != "" {
		command = pe
	} else if pe := mess.Message.GetConversation(); pe != "" {
		command = pe
	}

	if quotedMsg != nil && (quotedMsg.ImageMessage != nil || quotedMsg.VideoMessage != nil || quotedMsg.StickerMessage != nil) {
		if msg := quotedMsg.GetImageMessage(); msg != nil {
			media = msg
		} else if msg := quotedMsg.GetVideoMessage(); msg != nil {
			media = msg
		} else if msg := quotedMsg.GetStickerMessage(); msg != nil {
			media = msg
		}
	} else if mess.Message != nil && (mess.Message.ImageMessage != nil || mess.Message.VideoMessage != nil) {
		if msg := mess.Message.GetImageMessage(); msg != nil {
			media = msg
		} else if msg := mess.Message.GetVideoMessage(); msg != nil {
			media = msg
		}
	} else {
		media = nil
	}

	return &IMsg{
		AUpdate:   "https://github.com/amiruldev20/waSocket-bot",
		Timestamp: mess.Info.Timestamp,
		From:      mess.Info.Chat,
		Sender:    mess.Info.Sender,
		PushName:  mess.Info.PushName,
		ID:        mess.Info.ID,
		IsOwner:   isOwner,
		IsBot:     mess.Info.IsFromMe,
		IsGroup:   mess.Info.IsGroup,
		Query:     strings.Join(strings.Split(command, " ")[1:], ` `),
		Text:      command,
		Type: mess.Info.Type,
		Prefix:    strings.ToLower(strings.Split(command, " ")[0]),
		Media:     media,
		Msg:       mess.Message,
		IsImage: func() bool {
			if mess.Message.GetImageMessage() != nil {
				return true
			} else {
				return false
			}
		}(),
		IsAdmin: func() bool {
			if !mess.Info.IsGroup {
				return false
			}
			admin, err := sock.FetchGroupAdmin(mess.Info.Chat)
			if err != nil {
				return false
			}
			for _, v := range admin {
				if v == mess.Info.Sender.String() {
					return true
				}
			}
			return false
		}(),
		IsBotAdmin: func() bool {
			if !mess.Info.IsGroup {
				return false
			}
			admin, err := sock.FetchGroupAdmin(mess.Info.Chat)
			if err != nil {
				return false
			}
			for _, v := range admin {
				if v == botNum.String() {
					return true
				}
			}
			return false
		}(),
		Quoted: mess.Message.GetExtendedTextMessage().GetContextInfo(),
		IsQuotedImage: func() bool {
			if quotedMsg.GetImageMessage() != nil {
				return true
			} else {
				return false
			}
		}(),
		IsQuotedVideo: func() bool {
			if quotedMsg.GetVideoMessage() != nil {
				return true
			} else {
				return false
			}
		}(),
		IsQuotedSticker: func() bool {
			if quotedMsg.GetStickerMessage() != nil {
				return true
			} else {
				return false
			}
		}(),
		Exp: func() uint32 {
			var Expiration uint32
			if mess.Message.GetExtendedTextMessage() != nil {
				Expiration = mess.Message.GetExtendedTextMessage().GetContextInfo().GetExpiration()
			} else {
				Expiration = uint32(0)
			}
			return Expiration
		}(),
		React: func(r string) {
			var reaction = sock.WA.BuildReaction(mess.Info.Chat, mess.Info.Sender, mess.Info.ID, r)
			var x = []waSocket.SendRequestExtra{}
			sock.WA.SendMessage(context.Background(), mess.Info.Chat, reaction, x...)
		},
		Reply: func(text string) {
			var Expiration uint32

			if mess.Message.GetExtendedTextMessage() != nil {
				Expiration = mess.Message.GetExtendedTextMessage().GetContextInfo().GetExpiration()
			} else {
				Expiration = uint32(0)
			}

			sock.SendText(mess.Info.Chat, text, &waProto.ContextInfo{
				StanzaId:      &mess.Info.ID,
				Participant:   proto.String(mess.Info.Sender.String()),
				QuotedMessage: mess.Message,
				Expiration:    &Expiration,
			})

		},
		SendImage: func(url string, cap string){
			sock.Mimg(mess.Info.Chat, url, cap, mess)
		},
	}
}
