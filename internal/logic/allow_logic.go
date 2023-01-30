package logic

import (
	"context"
	. "meowchat-authorization-rpc/constant"

	"meowchat-authorization-rpc/internal/svc"
	"meowchat-authorization-rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type AllowLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAllowLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AllowLogic {
	return &AllowLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

var policies = map[string]func(*AllowLogic, *pb.AllowReq) bool{
	ObjectCommunity: (*AllowLogic).allowCommunity,
	ObjectNews:      (*AllowLogic).allowNews,
	ObjectNotice:    (*AllowLogic).allowNotice,
	ObjectPost:      (*AllowLogic).allowPost,
	ObjectCat:       (*AllowLogic).allowCat,
	ObjectMoment:    (*AllowLogic).allowMoment,
	ObjectComment:   (*AllowLogic).allowComment,
}

func (l *AllowLogic) Allow(in *pb.AllowReq) (*pb.AllowResp, error) {
	p := policies[in.Object]
	return &pb.AllowResp{
		Allow: p != nil && p(l, in),
	}, nil
}
