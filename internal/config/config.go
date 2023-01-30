package config

import "github.com/zeromicro/go-zero/zrpc"

type Config struct {
	zrpc.RpcServerConf
	CollectionRPC zrpc.RpcClientConf
	MomentRPC     zrpc.RpcClientConf
	SystemRPC     zrpc.RpcClientConf
	CommentRPC    zrpc.RpcClientConf
	PostRPC       zrpc.RpcClientConf
}
