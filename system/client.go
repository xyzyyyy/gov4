package system

/*
#include "opusreader.h"
*/

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"image/jpeg"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"strings"

	"github.com/amiruldev20/waSocket"

	waProto "github.com/amiruldev20/waSocket/binary/proto"
	"github.com/amiruldev20/waSocket/types"
	"github.com/amiruldev20/waSocket/types/events"

	"google.golang.org/protobuf/proto"
)

func NewClient(client *waSocket.Client) *Nc {
	return &Nc{
		WA: client,
	}
}

/* send text */
func (client *Nc) SendText(from types.JID, txt string, opts *waProto.ContextInfo) (waSocket.SendResponse, error) {
	ok, er := client.WA.SendMessage(context.Background(), from, &waProto.Message{
		ExtendedTextMessage: &waProto.ExtendedTextMessage{
			Text:        proto.String(txt),
			ContextInfo: opts,
		},
	})
	if er != nil {
		return waSocket.SendResponse{}, er
	}
	return ok, nil
}

func (client *Nc) SendWithNewsLestter(from types.JID, text string, newjid string, newserver int32, name string, opts *waProto.ContextInfo) (waSocket.SendResponse, error) {
	ok, er := client.SendText(from, text, &waProto.ContextInfo{
		ForwardedNewsletterMessageInfo: &waProto.ForwardedNewsletterMessageInfo{
			NewsletterJid:     proto.String(newjid),
			NewsletterName:    proto.String(name),
			ServerMessageId:   proto.Int32(newserver),
			ContentType:       waProto.ForwardedNewsletterMessageInfo_UPDATE.Enum(),
			AccessibilityText: proto.String(""),
		},
		IsForwarded:   proto.Bool(true),
		StanzaId:      opts.StanzaId,
		Participant:   opts.Participant,
		QuotedMessage: opts.QuotedMessage,
	})

	if er != nil {
		return waSocket.SendResponse{}, er
	}
	return ok, nil
}

/* sticker path */
func (client *Nc) StickerPath(from types.JID, path string, m IMsg) {
	stickerBuff,
		err := os.Open(path)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer stickerBuff.Close()

	stickerInfo,
		_ := stickerBuff.Stat()
	var stickerSize = stickerInfo.Size()
	stickerBytes := make([]byte, stickerSize)

	stickerBuffer := bufio.NewReader(stickerBuff)
	_,
		err = stickerBuffer.Read(stickerBytes)

	resp,
		err := client.WA.Upload(context.Background(), stickerBytes, waSocket.MediaImage)

	if err != nil {
		fmt.Println(err)
		return
	}

	_,
		err = client.WA.SendMessage(context.Background(), from, &waProto.Message{
		StickerMessage: &waProto.StickerMessage{
			Mimetype: proto.String("image/webp"),

			Url:           &resp.URL,
			DirectPath:    &resp.DirectPath,
			MediaKey:      resp.MediaKey,
			FileEncSha256: resp.FileEncSHA256,
			FileSha256:    resp.FileSHA256,
			FileLength:    &resp.FileLength,
			ContextInfo: &waProto.ContextInfo{
				StanzaId:      &m.ID,
				Participant:   proto.String(m.Sender.String()),
				QuotedMessage: m.Msg,
				Expiration:    &m.Exp,
			},
		},
	})

	if err != nil {
		fmt.Println(err)
		return
	}
}

/* send image */
func (client *Nc) SendImage(from types.JID, data any, caption string, m IMsg) (waSocket.SendResponse, error) {
	proc, err := ToByte(data)
	if err != nil {
		fmt.Println("Failed proccess get bytes!!")
	}
	pathimg := "tmp/" + GenerateID() + ".jpg"
	img,
		err := jpeg.Decode(bytes.NewReader(proc))
	if err != nil {
		fmt.Println(err)
	}
	width := img.Bounds().Dx()
	height := img.Bounds().Dy()
	Resize(proc, width, height, pathimg)
	imgThumb, err := os.ReadFile(pathimg)
	if err != nil {
		fmt.Println("Error read file:", err)
	}

	uploaded, err := client.WA.Upload(context.Background(), proc, waSocket.MediaImage)
	if err != nil {
		fmt.Printf("Failed to upload file: %v\n", err)
		return waSocket.SendResponse{}, err
	}
	resultImg := &waProto.Message{
		ImageMessage: &waProto.ImageMessage{
			Url:           proto.String(uploaded.URL),
			DirectPath:    proto.String(uploaded.DirectPath),
			MediaKey:      uploaded.MediaKey,
			Caption:       proto.String(caption),
			Mimetype:      proto.String(http.DetectContentType(proc)),
			JpegThumbnail: imgThumb,
			FileEncSha256: uploaded.FileEncSHA256,
			FileSha256:    uploaded.FileSHA256,
			FileLength:    proto.Uint64(uint64(len(proc))),
			ContextInfo: &waProto.ContextInfo{
				StanzaId:      &m.ID,
				Participant:   proto.String(m.Sender.String()),
				QuotedMessage: m.Msg,
				Expiration:    &m.Exp,
			},
		},
	}
	ok, _ := client.WA.SendMessage(context.Background(), from, resultImg)
	os.Remove(pathimg)
	return ok, nil
}

/* send image for m */
func (client *Nc) Mimg(from types.JID, data any, caption string, m *events.Message) (waSocket.SendResponse, error) {
	proc, err := ToByte(data)
	if err != nil {
		fmt.Println("Failed proccess get bytes!!")
	}
	pathimg := "tmp/" + GenerateID() + ".jpg"
	img,
		err := jpeg.Decode(bytes.NewReader(proc))
	if err != nil {
		fmt.Println(err)
	}
	width := img.Bounds().Dx()
	height := img.Bounds().Dy()
	Resize(proc, width, height, pathimg)
	imgThumb, err := os.ReadFile(pathimg)
	if err != nil {
		fmt.Println("Error read file:", err)
	}

	uploaded, err := client.WA.Upload(context.Background(), proc, waSocket.MediaImage)
	if err != nil {
		fmt.Printf("Failed to upload file: %v\n", err)
		return waSocket.SendResponse{}, err
	}
	var Expiration uint32

	if m.Message.GetExtendedTextMessage() != nil {
		Expiration = m.Message.GetExtendedTextMessage().GetContextInfo().GetExpiration()
	} else {
		Expiration = uint32(0)
	}
	resultImg := &waProto.Message{
		ImageMessage: &waProto.ImageMessage{
			Url:           proto.String(uploaded.URL),
			DirectPath:    proto.String(uploaded.DirectPath),
			MediaKey:      uploaded.MediaKey,
			Caption:       proto.String(caption),
			Mimetype:      proto.String(http.DetectContentType(proc)),
			JpegThumbnail: imgThumb,
			FileEncSha256: uploaded.FileEncSHA256,
			FileSha256:    uploaded.FileSHA256,
			FileLength:    proto.Uint64(uint64(len(proc))),
			ContextInfo: &waProto.ContextInfo{
				StanzaId:      &m.Info.ID,
				Participant:   proto.String(m.Info.Sender.String()),
				QuotedMessage: m.Message,
				Expiration:    &Expiration,
			},
		},
	}
	ok, _ := client.WA.SendMessage(context.Background(), from, resultImg)
	os.Remove(pathimg)
	return ok, nil
}

/* send video */
func (client *Nc) SendVideo(from types.JID, media any, caption string, m IMsg) (waSocket.SendResponse, error) {
	data, err := ToByte(media)
	if err != nil {
		fmt.Println("Failed proccess get bytes!!")
	}
	pathvid := "tmp/" + GenerateID() + ".mp4"
	paththumb := GenerateID() + ".jpg"
	GetThumbnail(data, 150, pathvid, paththumb)

	imgThumb, err := os.ReadFile("tmp/0" + paththumb)
	if err != nil {
		fmt.Println("Error read file:", err)
	}
	uploaded, err := client.WA.Upload(context.Background(), data, waSocket.MediaVideo)
	if err != nil {
		fmt.Printf("Failed to upload file: %v\n", err)
		return waSocket.SendResponse{}, err
	}
	// thumb
	up, err := client.WA.Upload(context.Background(), imgThumb, waSocket.MediaImage)
	if err != nil {
		fmt.Println("err upload thumb vid")
	}
	resultVideo := &waProto.Message{
		VideoMessage: &waProto.VideoMessage{
			Url:                 proto.String(uploaded.URL),
			DirectPath:          proto.String(uploaded.DirectPath),
			MediaKey:            uploaded.MediaKey,
			Caption:             proto.String(caption),
			Mimetype:            proto.String(http.DetectContentType(data)),
			ThumbnailDirectPath: &up.DirectPath,
			ThumbnailSha256:     up.FileSHA256,
			ThumbnailEncSha256:  up.FileEncSHA256,
			JpegThumbnail:       imgThumb,
			FileEncSha256:       uploaded.FileEncSHA256,
			FileSha256:          uploaded.FileSHA256,
			FileLength:          proto.Uint64(uint64(len(data))),
			ContextInfo: &waProto.ContextInfo{
				StanzaId:      &m.ID,
				Participant:   proto.String(m.Sender.String()),
				QuotedMessage: m.Msg,
				Expiration:    &m.Exp,
			},
		},
	}
	ok, er := client.WA.SendMessage(context.Background(), from, resultVideo)
	os.Remove("tmp/0" + paththumb)
	os.Remove(pathvid)
	os.Remove(paththumb)
	if er != nil {
		return waSocket.SendResponse{}, er
	}
	return ok, nil
}

/* send document */
func (client *Nc) SendDocument(from types.JID, media any, fileName string, caption string, m IMsg) (waSocket.SendResponse, error) {
	data, err := ToByte(media)
	if err != nil {
		fmt.Println("Failed proccess get bytes!!")
	}
	uploaded, err := client.WA.Upload(context.Background(), data, waSocket.MediaDocument)
	if err != nil {
		fmt.Printf("Failed to upload file: %v\n", err)
		return waSocket.SendResponse{}, err
	}

	// thumb
	conjp := "./tmp/thum.jpg"
	err = os.WriteFile(conjp, data, 0644)
	if err != nil {
		fmt.Println("Failed saved thumb")

	}
	imgThumb, err := os.ReadFile(conjp)

	resultDoc := &waProto.Message{
		DocumentMessage: &waProto.DocumentMessage{
			Url:           proto.String(uploaded.URL),
			DirectPath:    proto.String(uploaded.DirectPath),
			MediaKey:      uploaded.MediaKey,
			FileName:      proto.String(fileName),
			Caption:       proto.String(caption),
			Mimetype:      proto.String(http.DetectContentType(data)),
			FileEncSha256: uploaded.FileEncSHA256,
			FileSha256:    uploaded.FileSHA256,
			FileLength:    proto.Uint64(uint64(len(data))),
			JpegThumbnail: imgThumb,
			ContextInfo: &waProto.ContextInfo{
				StanzaId:      &m.ID,
				Participant:   proto.String(m.Sender.String()),
				QuotedMessage: m.Msg,
				Expiration:    &m.Exp,
			},
		},
	}
	ok, er := client.WA.SendMessage(context.Background(), from, resultDoc)
	os.Remove(conjp)
	if er != nil {
		return waSocket.SendResponse{}, er
	}
	return ok, nil
}

/* send audio */
func (client *Nc) SendAudio(from types.JID, media any, ptt bool, m IMsg) (waSocket.SendResponse, error) {
	data, err := ToByte(media)
	if err != nil {
		fmt.Println("Failed proccess get bytes!!")
	}
	uploaded, err := client.WA.Upload(context.Background(), data, waSocket.MediaAudio)
	if err != nil {
		fmt.Printf("Failed to upload file: %v\n", err)
		return waSocket.SendResponse{}, err
	}
	/*
	   waveform := make([]byte, C.WAVEFORM_SAMPLES_COUNT)
	   	for i := range c_waveform {
	   		waveform[i] = byte(c_waveform[i]) // convert while copying
	   	}*/
	resultAu := &waProto.Message{
		AudioMessage: &waProto.AudioMessage{
			Url:           proto.String(uploaded.URL),
			DirectPath:    proto.String(uploaded.DirectPath),
			MediaKey:      uploaded.MediaKey,
			Mimetype:      proto.String(http.DetectContentType(data)),
			FileEncSha256: uploaded.FileEncSHA256,
			FileSha256:    uploaded.FileSHA256,
			FileLength:    proto.Uint64(uint64(len(data))),
			Ptt:           proto.Bool(ptt),
			ContextInfo: &waProto.ContextInfo{
				StanzaId:      &m.ID,
				Participant:   proto.String(m.Sender.String()),
				QuotedMessage: m.Msg,
				Expiration:    &m.Exp,
			},
		},
	}
	ok, er := client.WA.SendMessage(context.Background(), from, resultAu)
	if er != nil {
		return waSocket.SendResponse{}, er
	}
	return ok, nil
}

/* kick member */
func (x *Nc) KickMember(from types.JID, jid types.JID) ([]types.GroupParticipant, error) {
	return x.WA.UpdateGroupParticipants(from, []types.JID{jid.ToNonAD()}, waSocket.ParticipantChangeRemove)
}

func (client *Nc) UploadImage(data []byte) (string, error) {
	bodyy := &bytes.Buffer{}
	writer := multipart.NewWriter(bodyy)
	part, _ := writer.CreateFormFile("file", "file")
	_, err := io.Copy(part, bytes.NewBuffer(data))
	if err != nil {
		return "", err
	}
	writer.Close()

	// Create request
	req, err := http.NewRequest("POST", "https://telegra.ph/upload", bodyy)
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Send request and handle response
	htt := &http.Client{}
	resp, err := htt.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("HTTP Error: %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var uploads []struct {
		Path string `json:"src"`
	}
	if err := json.Unmarshal(body, &uploads); err != nil {
		m := map[string]string{}
		if err := json.Unmarshal(data, &m); err != nil {
			return "", err
		}
		return "", fmt.Errorf("telegraph: %s", m["error"])
	}

	return "https://telegra.ph/" + uploads[0].Path, nil
}

func (client *Nc) ParseJID(arg string) (types.JID, bool) {
	if arg[0] == '+' {
		arg = arg[1:]
	}
	if !strings.ContainsRune(arg, '@') {
		return types.NewJID(arg, types.DefaultUserServer), true
	} else {
		recipient, err := types.ParseJID(arg)
		if err != nil {
			return recipient, false
		} else if recipient.User == "" {
			return recipient, false
		}
		return recipient, true
	}
}

func (client *Nc) FetchGroupAdmin(Jid types.JID) ([]string, error) {
	var Admin []string
	resp, err := client.WA.GetGroupInfo(Jid)
	if err != nil {
		return Admin, err
	} else {
		for _, group := range resp.Participants {
			if group.IsAdmin || group.IsSuperAdmin {
				Admin = append(Admin, group.JID.String())
			}
		}
	}
	return Admin, err
}

func (client *Nc) SendSticker(jid types.JID, data []byte, opts *waProto.ContextInfo) {
	uploaded, err := client.WA.Upload(context.Background(), data, waSocket.MediaImage)
	if err != nil {
		fmt.Printf("Failed to upload file: %v\n", err)
	}

	client.WA.SendMessage(context.Background(), jid, &waProto.Message{
		StickerMessage: &waProto.StickerMessage{
			Url:           proto.String(uploaded.URL),
			DirectPath:    proto.String(uploaded.DirectPath),
			MediaKey:      uploaded.MediaKey,
			Mimetype:      proto.String(http.DetectContentType(data)),
			FileEncSha256: uploaded.FileEncSHA256,
			FileSha256:    uploaded.FileSHA256,
			FileLength:    proto.Uint64(uint64(len(data))),
			ContextInfo:   opts,
		},
	})
}
