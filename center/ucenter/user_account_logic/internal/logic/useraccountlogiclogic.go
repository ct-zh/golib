package logic

import (
	"context"

	"user_account_logic/internal/svc"
	"user_account_logic/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type User_account_logicLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUser_account_logicLogic(ctx context.Context, svcCtx *svc.ServiceContext) *User_account_logicLogic {
	return &User_account_logicLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *User_account_logicLogic) User_account_logic(req *types.Request) (resp *types.Response, err error) {
	// todo: add your logic here and delete this line

	return
}
