package logic

import (
	. "github.com/xh-polaris/meowchat-authorization-rpc/constant"
	"github.com/xh-polaris/meowchat-authorization-rpc/pb"
	cat "github.com/xh-polaris/meowchat-collection-rpc/pb"
	comment "github.com/xh-polaris/meowchat-comment-rpc/pb"
	moment "github.com/xh-polaris/meowchat-moment-rpc/pb"
	post "github.com/xh-polaris/meowchat-post-rpc/pb"
	. "github.com/xh-polaris/meowchat-system-rpc/constant"
	system "github.com/xh-polaris/meowchat-system-rpc/pb"
)

// 社区权限
//  允许读，允许超级管理员、对应社区的管理员写
func (l *AllowLogic) allowCommunity(in *pb.AllowReq) bool {
	if in.Action == ActionRead {
		return true
	}

	return l.allowCommunityOrSuperAdmin(in.UserId, in.ObjectId)
}

// 通知权限
//  允许读，允许超级管理员、对应社区的管理员写
func (l *AllowLogic) allowNotice(in *pb.AllowReq) bool {
	if in.Action == ActionRead {
		return true
	}

	notice, _ := l.svcCtx.SystemRPC.RetrieveNotice(l.ctx, &system.RetrieveNoticeReq{Id: in.ObjectId})
	if notice == nil || notice.Notice == nil {
		return false
	}

	return l.allowCommunityOrSuperAdmin(in.UserId, notice.Notice.CommunityId)
}

// 轮播图权限
//  允许读，允许超级管理员、对应社区的管理员写
func (l *AllowLogic) allowNews(in *pb.AllowReq) bool {
	if in.Action == ActionRead {
		return true
	}

	news, _ := l.svcCtx.SystemRPC.RetrieveNews(l.ctx, &system.RetrieveNewsReq{Id: in.ObjectId})
	if news == nil || news.News == nil {
		return false
	}

	return l.allowCommunityOrSuperAdmin(in.UserId, news.News.CommunityId)
}

// 帖子权限
//  允许读，允许超级管理员、帖子发布者写
func (l *AllowLogic) allowPost(in *pb.AllowReq) bool {
	if in.Action == ActionRead || l.containsRole(in.UserId, RoleSuperAdmin) {
		return true
	}

	p, _ := l.svcCtx.PostRPC.RetrievePost(l.ctx, &post.RetrievePostReq{PostId: in.ObjectId})
	if p == nil || p.Post == nil {
		return false
	}

	return p.Post.UserId == in.UserId
}

// 猫咪信息权限
//  允许读，允许超级管理员、对应社区的管理员写
func (l *AllowLogic) allowCat(in *pb.AllowReq) bool {
	if in.Action == ActionRead {
		return true
	}

	c, _ := l.svcCtx.CollectionRPC.RetrieveCat(l.ctx, &cat.RetrieveCatReq{CatId: in.ObjectId})
	if c == nil || c.Cat == nil {
		return false
	}

	return l.allowCommunityOrSuperAdmin(in.UserId, c.Cat.CommunityId)
}

// 动态权限
//  允许读，允许超级管理员、对应社区的管理员、动态发布者写
func (l *AllowLogic) allowMoment(in *pb.AllowReq) bool {
	if in.Action == ActionRead {
		return true
	}

	m, _ := l.svcCtx.MomentRPC.RetrieveMoment(l.ctx, &moment.RetrieveMomentReq{MomentId: in.ObjectId})
	if m == nil || m.Moment == nil {
		return false
	}

	// 允许操作自己的moment
	if m.Moment.UserId == in.UserId {
		return true
	}

	return l.allowCommunityOrSuperAdmin(in.UserId, m.Moment.CommunityId)
}

// 评论权限
//  允许读，允许超级管理员、评论发布者写
func (l *AllowLogic) allowComment(in *pb.AllowReq) bool {
	if in.Action == ActionRead || l.containsRole(in.UserId, RoleSuperAdmin) {
		return true
	}

	c, _ := l.svcCtx.CommentRPC.RetrieveCommentById(l.ctx, &comment.RetrieveCommentByIdRequest{Id: in.ObjectId})
	if c == nil || c.Comment == nil {
		return false
	}

	// 允许操作自己的comment
	if c.Comment.AuthorId == in.UserId {
		return true
	}

	// 如果对评论从属对象有权限，对其下所有评论也有权限
	allowParentReq := &pb.AllowReq{
		UserId:   in.UserId,
		ObjectId: c.Comment.ParentId,
		Action:   in.Action,
	}
	switch c.Comment.Type {
	case ObjectMoment:
		return l.allowMoment(allowParentReq)
	case ObjectPost:
		return l.allowPost(allowParentReq)
	}

	return false
}
