package main

import (
	"context"
	"fmt"
	"net/http"

	_ "mywabot/cmd"
	_ "mywabot/cmd/ai"
	_ "mywabot/cmd/convert"
	_ "mywabot/cmd/main"
	_ "mywabot/cmd/owner"
	"mywabot/config"
	"mywabot/system"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/amiruldev20/waSocket"
	waProto "github.com/amiruldev20/waSocket/binary/proto"
	"github.com/amiruldev20/waSocket/store"
	"github.com/amiruldev20/waSocket/store/sqlstore"
	"github.com/amiruldev20/waSocket/types"
	"github.com/amiruldev20/waSocket/types/events"
	waLog "github.com/amiruldev20/waSocket/util/log"
	_ "github.com/mattn/go-sqlite3"
	"github.com/mdp/qrterminal"
	"google.golang.org/protobuf/proto"
)

func init() {
	store.DeviceProps.PlatformType = waProto.DeviceProps_EDGE.Enum()
	store.DeviceProps.Os = proto.String(config.Name)
}

func main() {
	go startWeb()
	startBot()
}
func startWeb() {
	/* web server */
	port := os.Getenv("PORT")
	if port == "" {
		port = "1337" // Port default jika tidak ada yang disetel
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "waSocket Bot Connected")
	})

	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
	/* end web server */
}
func startBot() {
	dbLog := waLog.Stdout("Database", "ERROR", true)
	container, err := sqlstore.New("sqlite3", "file:mywabot.db?_foreign_keys=on", dbLog)
	if err != nil {
		panic(err)
	}
	deviceStore, err := container.GetFirstDevice()
	if err != nil {
		panic(err)
	}
	clientLog := waLog.Stdout("Client", "ERROR", true)
	client := waSocket.NewClient(deviceStore, clientLog)
	handler := registerHandler(client)
	fmt.Println("Connecting to waSocket...")
	client.AddEventHandler(handler)

	client.PrePairCallback = func(jid types.JID, platform, businessName string) bool {
		fmt.Println("Connected to waSocket!!")
		return true
	}

	if client.Store.ID == nil {
		// No ID stored, new login
		// Switch Mode
		switch int(questLogin()) {
		case 1:

			if err := client.Connect(); err != nil {
				panic(err)
			}

			code, err := client.PairPhone(config.Bot, true, waSocket.PairClientChrome, "Chrome (Linux)")
			if err != nil {
				panic(err)
			}

			fmt.Println("Kode Login: " + code)
			break
		case 2:
			qrChan, _ := client.GetQRChannel(context.Background())
			if err := client.Connect(); err != nil {
				panic(err)
			}
			for evt := range qrChan {
				switch string(evt.Event) {
				case "code":
					qrterminal.GenerateHalfBlock(evt.Code, qrterminal.L, os.Stdout)
					fmt.Println("Scan Qrnya!!")
					break
				}
			}
			break
		default:
			panic("Pilih apa?")
		}
	} else {
		// Already logged in, just connect
		if err := client.Connect(); err != nil {
			panic(err)
		}
		fmt.Println("Connected to waSocket!!")
		own, _ := types.ParseJID("62882008211320@s.whatsapp.net")
		client.SendMessage(context.Background(), own, &waProto.Message{Conversation: proto.String("Bot Connected!!")})

	}

	// Listen to Ctrl+C (you can also do something else that prevents the program from exiting)
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	client.Disconnect()
}

func registerHandler(client *waSocket.Client) func(evt interface{}) {
	return func(evt interface{}) {
		switch v := evt.(type) {
		case *events.Message:
			sock := system.NewClient(client)
			m := system.NewSmsg(v, sock)
			if time.Since(m.Timestamp).Seconds() > 15 {
				return
			}

			if config.Self && !m.IsOwner {
				return
			}

			go system.Get(sock, m)
			return
		case *events.LoggedOut:
			con := evt.(*events.LoggedOut)
			if !con.OnConnect {
				fmt.Println(client.Store.ID.User)
				fmt.Println("LogOut Reason : " + con.Reason.String())
				panic("Log Out")
			}
			break
		case *events.Connected, *events.PushNameSetting:
			if len(client.Store.PushName) == 0 {
				return
			}
			client.SendPresence(types.PresenceAvailable)
		}
	}
}

func questLogin() int {
	fmt.Println("Silahlan Pilih Opsi Login :")
	fmt.Println("1. Pairing Code")
	fmt.Println("2. Qr")
	fmt.Print("Pilih : ")
	var input int
	_, err := fmt.Scanln(&input)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 0
	}

	return input
}
