package logic

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"user_profile_logic/internal/svc"
	"user_profile_logic/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type Profile struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewProfile(ctx context.Context, svcCtx *svc.ServiceContext) *Profile {
	return &Profile{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (p *Profile) GetProfile(req types.GetRsp) (resp types.GetReply, err error) {
	resp.Data = types.Profile{Uid: 111}
	return
}

func (p *Profile) MGetProfile(req types.MGetRsp) (resp types.MGetReply, err error) {
	uidSlice := strings.Split(req.Uids, ",")
	if len(uidSlice) == 0 {
		return resp, fmt.Errorf("invalid params")
	}
	resp.Data = make([]types.Profile, 0)
	for _, s := range uidSlice {
		uid, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			continue
		}
		resp.Data = append(resp.Data, types.Profile{Uid: uid})
	}

	return
}

func (p *Profile) UpdateProfile(uid int64, data map[string]string) error {
	return nil
}
