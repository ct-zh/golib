package handler

import (
	"net/http"

	"user_profile_logic/internal/logic"
	"user_profile_logic/internal/svc"
	"user_profile_logic/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetProfile(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetRsp
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		p := logic.NewProfile(r.Context(), ctx)
		resp, err := p.GetProfile(req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}

func MGetProfile(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		var req types.MGetRsp
		if err := httpx.ParsePath(request, &req); err != nil {
			httpx.Error(writer, err)
			return
		}

		p := logic.NewProfile(request.Context(), ctx)
		resp, err := p.MGetProfile(req)
		if err != nil {
			httpx.Error(writer, err)
		} else {
			httpx.OkJson(writer, resp)
		}
	}
}

func Update(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uid := struct {
			Uid int64 `json:"uid"`
		}{}
		if err := httpx.ParsePath(r, &uid); err != nil {
			httpx.Error(w, err)
			return
		}
		body := make(map[string]string)
		if err := httpx.ParseJsonBody(r, &body); err != nil {
			httpx.Error(w, err)
			return
		}

		p := logic.NewProfile(r.Context(), ctx)
		err := p.UpdateProfile(uid.Uid, body)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.Ok(w)
		}
	}
}
