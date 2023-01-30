package svc

import (
	"github.com/xh-polaris/meowchat-collection-rpc/collectionrpc"
	"github.com/xh-polaris/meowchat-comment-rpc/commentrpc"
	"github.com/xh-polaris/meowchat-moment-rpc/momentrpc"
	"github.com/xh-polaris/meowchat-post-rpc/postrpc"
	"github.com/xh-polaris/meowchat-system-rpc/systemrpc"
	"github.com/zeromicro/go-zero/zrpc"
	"meowchat-authorization-rpc/internal/config"
)

type ServiceContext struct {
	Config        config.Config
	CollectionRPC collectionrpc.CollectionRpc
	MomentRPC     momentrpc.MomentRpc
	SystemRPC     systemrpc.SystemRpc
	CommentRPC    commentrpc.CommentRpc
	PostRPC       postrpc.PostRpc
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:        c,
		CollectionRPC: collectionrpc.NewCollectionRpc(zrpc.MustNewClient(c.CollectionRPC)),
		MomentRPC:     momentrpc.NewMomentRpc(zrpc.MustNewClient(c.MomentRPC)),
		SystemRPC:     systemrpc.NewSystemRpc(zrpc.MustNewClient(c.SystemRPC)),
		CommentRPC:    commentrpc.NewCommentRpc(zrpc.MustNewClient(c.CommentRPC)),
		PostRPC:       postrpc.NewPostRpc(zrpc.MustNewClient(c.PostRPC)),
	}
}
