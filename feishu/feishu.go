package feishu

import (
	"fmt"

	"github.com/imroc/req/v3"
)

type FxBot struct {
	Client *req.Client
	Config interface{}
}

type Config struct {
	WebhookURL string
}

type Message struct {
	MsgType string  `json:"msg_type"`
	Content Content `json:"content"`
}

type Content struct {
	Text     string   `json:"text,omitempty"`
	Post     PostBody `json:"post,omitempty"`
	ImageKey string   `json:"image_key,omitempty"`
	Card     Card     `json:"card,omitempty"`
}

type PostBody struct {
	ZHCN PostBodyContents `json:"zh_cn"`
}

type PostBodyContents struct {
	Title   string            `json:"title"`
	Content []PostBodyContent `json:"content"`
}

type PostBodyContent struct {
	Tag    string `json:"tag"`
	Text   string `json:"text,omitempty"`
	Href   string `json:"href,omitempty"`
	UserID string `json:"user_id,omitempty"`
}

// https://open.feishu.cn/tool/cardbuilder?from=custom_bot_doc

type Card struct {
	Config   CardConfig    `json:"config,omitempty"`
	Header   CardHeader    `json:"header,omitempty"`
	Elements []CardElement `json:"elements,omitempty"`
}

type CardConfig struct {
	WideScreenMode bool `json:"wide_screen_mode"`
	EnableForward  bool `json:"enable_forward"`
}

type CardHeader struct {
	Template string   `json:"template"`
	Title    CardText `json:"title"`
}

type CardElement struct {
	Tag     string       `json:"tag"`
	Fields  []CardField  `json:"fields,omitempty"`
	Text    CardText     `json:"text,omitempty"`
	Actions []CardAction `json:"action,omitempty"`
}

type CardField struct {
	IsShort bool     `json:"is_short"`
	Text    CardText `json:"text"`
}

type CardText struct {
	Content string `json:"content"`
	Tag     string `json:"tag"`
}

type CardAction struct {
	Tag  string   `json:"tag"`
	Text CardText `json:"text,omitempty"`
	URL  string   `json:"url,omitempty"`
	Type string   `json:"type,omitempty"`
}

type Response struct {
	Errcode int    `json:"errcode"`
	Errmsg  string `json:"errmsg"`
}

func (fx FxBot) Send(msg *Message) (resp *Response, err error) {
	resp = &Response{}
	r, err := fx.Client.R().
		SetBodyJsonMarshal(msg).
		EnableDumpWithoutRequest().
		SetResult(resp).
		Post("") // TODO webhook url
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

func (fx FxBot) Debug(debug bool) {
	if debug {
		fx.Client.EnableDumpAll().EnableDebugLog().EnableTraceAll()
	} else {
		fx.Client.DisableDebugLog().DisableDumpAll().DisableTraceAll()
	}
}
