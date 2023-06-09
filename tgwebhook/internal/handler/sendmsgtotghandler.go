package handler

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"tgwebhook/internal/logic"
	"tgwebhook/internal/svc"
	"tgwebhook/internal/types"
)

func sendmsgtotgHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.TgSendReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewSendmsgtotgLogic(r.Context(), svcCtx)
		resp, err := l.Sendmsgtotg(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
