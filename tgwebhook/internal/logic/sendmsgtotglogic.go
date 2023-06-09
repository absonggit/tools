package logic

import (
	"context"
	"log"
	"tgwebhook/internal/svc"
	"tgwebhook/internal/types"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/zeromicro/go-zero/core/logx"
)

type SendmsgtotgLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSendmsgtotgLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendmsgtotgLogic {
	return &SendmsgtotgLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SendmsgtotgLogic) Sendmsgtotg(req *types.TgSendReq) (resp *types.TgSendResp, err error) {
	bot, err := tgbotapi.NewBotAPI(req.Token)
	if err != nil {
		log.Panic(err)
	}
	resp = new(types.TgSendResp)
	msg := tgbotapi.NewMessage(req.Chatid, req.Text)
	bot.Send(msg)
	if err != nil {
		log.Panic(err)
	}
	resp.Status = "true"
	resp.Message = req.Text
	return
}
