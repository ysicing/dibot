package dibot

import (
	"github.com/imroc/req/v3"
	"github.com/ysicing/dibot/workwx"
)

type DiBot interface {
	Debug(bool)
}

func NewDiBot(t, webhook string) DiBot {
	client := req.C()
	switch t {
	case "workwx":
		return workwx.WeBot{Client: client, WebhookURL: webhook}
	default:
		return workwx.WeBot{Client: client, WebhookURL: webhook}
	}
}
