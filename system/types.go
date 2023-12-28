package system

import (
	"time"

	"github.com/amiruldev20/waSocket"
	"github.com/amiruldev20/waSocket/binary/proto"
	waProto "github.com/amiruldev20/waSocket/binary/proto"
	"github.com/amiruldev20/waSocket/types"
	"github.com/amiruldev20/waSocket/types/events"
)

type Nc struct {
	WA *waSocket.Client
}

type Event struct {
	Client  *waSocket.Client
	Message *events.Message
	Context *proto.ContextInfo
	Pattern string
	Args    []string
	Text    string
}

type ICmd struct {
	Name     string
	Cmd      []string
	Desc     string
	Tags     string
	Prefix   bool
	IsOwner  bool
	IsMedia  bool
	IsQuery  bool
	ValueQ   string
	IsGroup  bool
	IsAdmin  bool
	IsBotAdm bool
	IsWait   bool
	After    func(client *Nc, m *IMsg)
	Exec     func(client *Nc, m *IMsg)
}

type IMsg struct {
	AUpdate         string
	Exp             uint32
	Timestamp       time.Time
	From            types.JID
	IsBot           bool
	Sender          types.JID
	PushName        string
	ID              string
	Type            string
	IsOwner         bool
	IsGroup         bool
	Query           string
	Text            string
	Prefix          string
	IsImage         bool
	IsVideo         bool
	IsQuotedImage   bool
	IsQuotedVideo   bool
	IsQuotedSticker bool
	IsAdmin         bool
	IsBotAdmin      bool
	Media           waSocket.DownloadableMessage
	Msg             *waProto.Message
	Quoted          *waProto.ContextInfo
	React           func(text string)
	Reply           func(text string)
	SendImage       func(url string, cap string)
}
