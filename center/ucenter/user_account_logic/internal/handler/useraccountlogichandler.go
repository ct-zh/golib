package handler

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"user_account_logic/internal/logic"
	"user_account_logic/internal/svc"
	"user_account_logic/internal/types"
)

func User_account_logicHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.Request
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := logic.NewUser_account_logicLogic(r.Context(), svcCtx)
		resp, err := l.User_account_logic(&req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
