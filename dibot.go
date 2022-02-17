package dibot

import (
	"github.com/imroc/req/v3"
	"github.com/ysicing/dibot/feishu"
	"github.com/ysicing/dibot/workwx"
)

type DiBot interface {
	Debug(bool)
}

func NewDiBot(t string, config interface{}) DiBot {
	client := req.C()
	switch t {
	case "workwx":
		return workwx.WeBot{Client: client, Config: config}
	case "feishu":
		return feishu.FxBot{Client: client, Config: config}
	default:
		return feishu.FxBot{Client: client, Config: config}
	}
}
