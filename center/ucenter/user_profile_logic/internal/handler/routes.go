package handler

import (
	"net/http"

	"user_profile_logic/internal/svc"

	"github.com/zeromicro/go-zero/rest"
)

func RegisterHandlers(server *rest.Server, serverCtx *svc.ServiceContext) {
	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodGet,
				Path:    "/api_server/v1/users/profile/mget/:uids/:owner_id",
				Handler: MGetProfile(serverCtx),
			},
			{
				Method:  http.MethodGet,
				Path:    "/api_server/v1/users/profile/get/:uid/:id",
				Handler: GetProfile(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/api_server/v1/users/profile/update/:uid",
				Handler: Update(serverCtx),
			},
		},
	)
}
