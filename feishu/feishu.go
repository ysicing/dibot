package feishu

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"

	"github.com/imroc/req/v3"
)

type FxBot struct {
	Client     *req.Client
	Config     interface{}
	webhookURL string
	uploadURL  string
}

type TextMessage struct {
	Content             string   `json:"content"`
	MentionedList       []string `json:"mentioned_list,omitempty"`
	MentionedMobileList []string `json:"mentioned_mobile_list,omitempty"`
}

type MarkdownMessage struct {
	Content string `json:"content"`
}

type Message struct {
	Msgtype  string           `json:"msgtype"`
	Text     *TextMessage     `json:"text,omitempty"`
	Markdown *MarkdownMessage `json:"markdown,omitempty"`
	File     *FileMessage     `json:"file,omitempty"`
}

type Response struct {
	Errcode int    `json:"errcode"`
	Errmsg  string `json:"errmsg"`
}

type UploadResponse struct {
	Response
	Type      string `json:"type"`
	MediaID   string `json:"media_id"`
	CreatedAt string `json:"created_at"`
}

type FileMessage struct {
	MediaID string `json:"media_id"`
}

func (fx FxBot) getUploadURL() string {
	if fx.uploadURL != "" {
		return fx.uploadURL
	}
	fx.uploadURL = strings.ReplaceAll(fx.webhookURL, "webhook/send", "webhook/upload_media")
	return fx.uploadURL
}

func (fx FxBot) Send(msg *Message) (resp *Response, err error) {
	resp = &Response{}
	r, err := fx.Client.R().
		SetBodyJsonMarshal(msg).
		EnableDumpWithoutRequest().
		SetResult(resp).
		Post(fx.webhookURL)
	if err != nil {
		return
	}
	if !r.IsSuccess() {
		err = fmt.Errorf("bad response:\n%s", r.Dump())
		return
	}
	if resp.Errcode != 0 {
		err = fmt.Errorf(resp.Errmsg)
	}
	return
}

func (fx FxBot) SendFileContent(filename string, content []byte) (resp *Response, err error) {
	upload, err := fx.Upload(filename, content)
	if err != nil {
		return
	}
	file := &FileMessage{
		MediaID: upload.MediaID,
	}
	return fx.Send(&Message{
		Msgtype: "file",
		File:    file,
	})
}

func (fx FxBot) Upload(filename string, data []byte) (resp *UploadResponse, err error) {
	resp = &UploadResponse{}
	cd := new(req.ContentDisposition)
	cd.Add("filelength", strconv.Itoa(len(data)))
	r, err := fx.Client.R().
		SetFileUpload(req.FileUpload{
			ParamName:               "media",
			FileName:                filename,
			File:                    bytes.NewReader(data),
			ExtraContentDisposition: cd,
		}).EnableDumpWithoutRequest().
		SetQueryParam("type", "file").
		SetResult(resp).
		Post(fx.getUploadURL())
	if err != nil {
		return
	}
	if !r.IsSuccess() {
		err = fmt.Errorf("bad response:\n%s", r.Dump())
		return
	}
	if resp.Errcode != 0 {
		err = fmt.Errorf(resp.Errmsg)
	}
	return
}

func (fx FxBot) SendMarkdownContent(markdown string) (resp *Response, err error) {
	return fx.SendMarkdown(&MarkdownMessage{
		Content: markdown,
	})
}

func (fx FxBot) SendMarkdown(markdown *MarkdownMessage) (resp *Response, err error) {
	msg := &Message{Msgtype: "markdown", Markdown: markdown}
	return fx.Send(msg)
}

func (fx FxBot) SendText(text *TextMessage) (resp *Response, err error) {
	msg := &Message{Msgtype: "text", Text: text}
	return fx.Send(msg)
}

func (fx FxBot) SendTextContent(text string) (resp *Response, err error) {
	msg := &TextMessage{
		Content: text,
	}
	return fx.SendText(msg)
}

func (fx FxBot) Debug(debug bool) {
	if debug {
		fx.Client.EnableDumpAll().EnableDebugLog().EnableTraceAll()
	} else {
		fx.Client.DisableDebugLog().DisableDumpAll().DisableTraceAll()
	}
}
