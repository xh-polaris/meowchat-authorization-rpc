package logic

import (
	"context"
	. "github.com/golang/mock/gomock"
	. "github.com/smartystreets/goconvey/convey"
	. "github.com/xh-polaris/meowchat-authorization-rpc/constant"
	"github.com/xh-polaris/meowchat-authorization-rpc/internal/config"
	"github.com/xh-polaris/meowchat-authorization-rpc/internal/logic/mock"
	"github.com/xh-polaris/meowchat-authorization-rpc/internal/svc"
	pb2 "github.com/xh-polaris/meowchat-authorization-rpc/pb"
	pb4 "github.com/xh-polaris/meowchat-comment-rpc/pb"
	pb5 "github.com/xh-polaris/meowchat-moment-rpc/pb"
	pb3 "github.com/xh-polaris/meowchat-post-rpc/pb"
	. "github.com/xh-polaris/meowchat-system-rpc/constant"
	"github.com/xh-polaris/meowchat-system-rpc/pb"
	"testing"
	_ "unsafe"
)

func TestAllowLogic_Allow_Community(t *testing.T) {
	ctrl := NewController(t)
	defer ctrl.Finish()

	mockCollectionRpc := mock.NewMockCollectionRpc(ctrl)
	mockMomentRpc := mock.NewMockMomentRpc(ctrl)
	mockSystemRpc := mock.NewMockSystemRpc(ctrl)
	mockCommentRpc := mock.NewMockCommentRpc(ctrl)
	mockPostRpc := mock.NewMockPostRpc(ctrl)

	svcCtx := &svc.ServiceContext{
		Config:        config.Config{},
		CollectionRPC: mockCollectionRpc,
		MomentRPC:     mockMomentRpc,
		SystemRPC:     mockSystemRpc,
		CommentRPC:    mockCommentRpc,
		PostRPC:       mockPostRpc,
	}
	l := NewAllowLogic(context.Background(), svcCtx)

	Convey("允许读", t, func() {
		allow, _ := l.Allow(&pb2.AllowReq{
			Object: ObjectCommunity,
			Action: ActionRead,
		})
		So(allow.Allow, ShouldBeTrue)
	})

	Convey("不允许普通用户写", t, func() {
		mockSystemRpc.EXPECT().RetrieveUserRole(Any(), Any(), Any()).Return(&pb.RetrieveUserRoleResp{
			Roles: []*pb.Role{
				{
					Type: RoleUser,
				},
			},
		}, nil)
		allow, _ := l.Allow(&pb2.AllowReq{
			Object: ObjectCommunity,
			Action: ActionWrite,
		})
		So(allow.Allow, ShouldBeFalse)
	})

	Convey("允许超级管理员", t, func() {
		mockSystemRpc.EXPECT().RetrieveUserRole(Any(), Any(), Any()).Return(&pb.RetrieveUserRoleResp{
			Roles: []*pb.Role{
				{
					Type: RoleSuperAdmin,
				},
			},
		}, nil)
		allow, _ := l.Allow(&pb2.AllowReq{
			Object: ObjectCommunity,
			Action: ActionWrite,
		})
		So(allow.Allow, ShouldBeTrue)
	})

	Convey("允许社区管理员", t, func() {
		mockSystemRpc.EXPECT().RetrieveUserRole(Any(), Any(), Any()).Return(&pb.RetrieveUserRoleResp{
			Roles: []*pb.Role{
				{
					Type:        RoleCommunityAdmin,
					CommunityId: "TestCommId",
				},
			},
		}, nil)
		allow, _ := l.Allow(&pb2.AllowReq{
			Object:   ObjectCommunity,
			ObjectId: "TestCommId",
			Action:   ActionWrite,
		})
		So(allow.Allow, ShouldBeTrue)
	})

	Convey("社区管理员ID不符", t, func() {
		mockSystemRpc.EXPECT().RetrieveUserRole(Any(), Any()).Return(&pb.RetrieveUserRoleResp{
			Roles: []*pb.Role{
				{
					Type:        RoleCommunityAdmin,
					CommunityId: "TestCommId",
				},
			},
		}, nil)
		mockSystemRpc.EXPECT().RetrieveCommunity(Any(), Any()).Return(&pb.RetrieveCommunityResp{
			Community: &pb.Community{
				Id: "TestCommId2",
			},
		}, nil)
		allow, _ := l.Allow(&pb2.AllowReq{
			Object:   ObjectCommunity,
			ObjectId: "TestCommId2",
			Action:   ActionWrite,
		})
		So(allow.Allow, ShouldBeFalse)
	})

	Convey("允许社区管理员子社区", t, func() {
		mockSystemRpc.EXPECT().RetrieveUserRole(Any(), Any()).Return(&pb.RetrieveUserRoleResp{
			Roles: []*pb.Role{
				{
					Type:        RoleCommunityAdmin,
					CommunityId: "Parent",
				},
			},
		}, nil)
		mockSystemRpc.EXPECT().RetrieveCommunity(Any(), Any()).Return(&pb.RetrieveCommunityResp{
			Community: &pb.Community{
				Id:       "Child",
				ParentId: "Parent",
			},
		}, nil)
		allow, _ := l.Allow(&pb2.AllowReq{
			Object:   ObjectCommunity,
			ObjectId: "Child",
			Action:   ActionWrite,
		})
		So(allow.Allow, ShouldBeTrue)
	})

	Convey("多社区管理员测试", t, func() {
		mockSystemRpc.EXPECT().RetrieveUserRole(Any(), Any()).Return(&pb.RetrieveUserRoleResp{
			Roles: []*pb.Role{
				{
					Type:        RoleCommunityAdmin,
					CommunityId: "AnyParent",
				},
				{
					Type:        RoleCommunityAdmin,
					CommunityId: "Parent",
				},
			},
		}, nil)
		mockSystemRpc.EXPECT().RetrieveCommunity(Any(), Any()).Return(&pb.RetrieveCommunityResp{
			Community: &pb.Community{
				Id:       "AnyParent",
				ParentId: "zz",
			},
		}, nil)
		mockSystemRpc.EXPECT().RetrieveCommunity(Any(), Any()).Return(&pb.RetrieveCommunityResp{
			Community: &pb.Community{
				Id:       "Child",
				ParentId: "Parent",
			},
		}, nil)
		allow, _ := l.Allow(&pb2.AllowReq{
			Object:   ObjectCommunity,
			ObjectId: "Child",
			Action:   ActionWrite,
		})
		So(allow.Allow, ShouldBeTrue)
	})
}

func TestAllowLogic_Allow_Notice(t *testing.T) {
	ctrl := NewController(t)
	defer ctrl.Finish()

	mockCollectionRpc := mock.NewMockCollectionRpc(ctrl)
	mockMomentRpc := mock.NewMockMomentRpc(ctrl)
	mockSystemRpc := mock.NewMockSystemRpc(ctrl)
	mockCommentRpc := mock.NewMockCommentRpc(ctrl)
	mockPostRpc := mock.NewMockPostRpc(ctrl)

	svcCtx := &svc.ServiceContext{
		Config:        config.Config{},
		CollectionRPC: mockCollectionRpc,
		MomentRPC:     mockMomentRpc,
		SystemRPC:     mockSystemRpc,
		CommentRPC:    mockCommentRpc,
		PostRPC:       mockPostRpc,
	}
	l := NewAllowLogic(context.Background(), svcCtx)

	Convey("允许读", t, func() {
		allow, _ := l.Allow(&pb2.AllowReq{
			Object: ObjectNotice,
			Action: ActionRead,
		})
		So(allow.Allow, ShouldBeTrue)
	})

	Convey("不允许普通用户写", t, func() {
		mockSystemRpc.EXPECT().RetrieveUserRole(Any(), Any(), Any()).Return(&pb.RetrieveUserRoleResp{
			Roles: []*pb.Role{
				{
					Type: RoleUser,
				},
			},
		}, nil)
		allow, _ := l.Allow(&pb2.AllowReq{
			Object: ObjectCommunity,
			Action: ActionWrite,
		})
		So(allow.Allow, ShouldBeFalse)
	})

	Convey("允许超级管理员", t, func() {
		mockSystemRpc.EXPECT().RetrieveNotice(Any(), Any()).Return(&pb.RetrieveNoticeResp{Notice: &pb.Notice{}}, nil)
		mockSystemRpc.EXPECT().RetrieveUserRole(Any(), Any(), Any()).Return(&pb.RetrieveUserRoleResp{
			Roles: []*pb.Role{
				{
					Type: RoleSuperAdmin,
				},
			},
		}, nil)
		allow, _ := l.Allow(&pb2.AllowReq{
			Object: ObjectNotice,
			Action: ActionWrite,
		})
		So(allow.Allow, ShouldBeTrue)
	})

	Convey("允许社区管理员", t, func() {
		mockSystemRpc.EXPECT().RetrieveNotice(Any(), Any()).Return(&pb.RetrieveNoticeResp{
			Notice: &pb.Notice{
				Id:          "NoticeId",
				CommunityId: "TestCommId",
			},
		}, nil)
		mockSystemRpc.EXPECT().RetrieveUserRole(Any(), Any(), Any()).Return(&pb.RetrieveUserRoleResp{
			Roles: []*pb.Role{
				{
					Type:        RoleCommunityAdmin,
					CommunityId: "TestCommId",
				},
			},
		}, nil)
		allow, _ := l.Allow(&pb2.AllowReq{
			Object:   ObjectNotice,
			ObjectId: "NoticeId",
			Action:   ActionWrite,
		})
		So(allow.Allow, ShouldBeTrue)
	})

	Convey("社区管理员ID不符", t, func() {
		mockSystemRpc.EXPECT().RetrieveNotice(Any(), Any()).Return(&pb.RetrieveNoticeResp{
			Notice: &pb.Notice{
				Id:          "NoticeId",
				CommunityId: "TestCommId2",
			},
		}, nil)
		mockSystemRpc.EXPECT().RetrieveUserRole(Any(), Any()).Return(&pb.RetrieveUserRoleResp{
			Roles: []*pb.Role{
				{
					Type:        RoleCommunityAdmin,
					CommunityId: "TestCommId",
				},
			},
		}, nil)
		mockSystemRpc.EXPECT().RetrieveCommunity(Any(), Any()).Return(&pb.RetrieveCommunityResp{
			Community: &pb.Community{
				Id: "TestCommId2",
			},
		}, nil)
		allow, _ := l.Allow(&pb2.AllowReq{
			Object:   ObjectNotice,
			ObjectId: "NoticeId",
			Action:   ActionWrite,
		})
		So(allow.Allow, ShouldBeFalse)
	})

	Convey("允许社区管理员子社区", t, func() {
		mockSystemRpc.EXPECT().RetrieveNotice(Any(), Any()).Return(&pb.RetrieveNoticeResp{
			Notice: &pb.Notice{
				Id:          "NoticeId",
				CommunityId: "Child",
			},
		}, nil)
		mockSystemRpc.EXPECT().RetrieveUserRole(Any(), Any()).Return(&pb.RetrieveUserRoleResp{
			Roles: []*pb.Role{
				{
					Type:        RoleCommunityAdmin,
					CommunityId: "Parent",
				},
			},
		}, nil)
		mockSystemRpc.EXPECT().RetrieveCommunity(Any(), Any()).Return(&pb.RetrieveCommunityResp{
			Community: &pb.Community{
				Id:       "Child",
				ParentId: "Parent",
			},
		}, nil)
		allow, _ := l.Allow(&pb2.AllowReq{
			Object:   ObjectNotice,
			ObjectId: "NoticeId",
			Action:   ActionWrite,
		})
		So(allow.Allow, ShouldBeTrue)
	})

	Convey("多社区管理员测试", t, func() {
		mockSystemRpc.EXPECT().RetrieveNotice(Any(), Any()).Return(&pb.RetrieveNoticeResp{
			Notice: &pb.Notice{
				Id:          "NoticeId",
				CommunityId: "Child",
			},
		}, nil)
		mockSystemRpc.EXPECT().RetrieveUserRole(Any(), Any()).Return(&pb.RetrieveUserRoleResp{
			Roles: []*pb.Role{
				{
					Type:        RoleCommunityAdmin,
					CommunityId: "AnyParent",
				},
				{
					Type:        RoleCommunityAdmin,
					CommunityId: "Parent",
				},
			},
		}, nil)
		mockSystemRpc.EXPECT().RetrieveCommunity(Any(), Any()).Return(&pb.RetrieveCommunityResp{
			Community: &pb.Community{
				Id:       "AnyParent",
				ParentId: "zz",
			},
		}, nil)
		mockSystemRpc.EXPECT().RetrieveCommunity(Any(), Any()).Return(&pb.RetrieveCommunityResp{
			Community: &pb.Community{
				Id:       "Child",
				ParentId: "Parent",
			},
		}, nil)
		allow, _ := l.Allow(&pb2.AllowReq{
			Object:   ObjectNotice,
			ObjectId: "NoticeId",
			Action:   ActionWrite,
		})
		So(allow.Allow, ShouldBeTrue)
	})
}

func TestAllowLogic_Allow_Post(t *testing.T) {
	ctrl := NewController(t)
	defer ctrl.Finish()

	mockCollectionRpc := mock.NewMockCollectionRpc(ctrl)
	mockMomentRpc := mock.NewMockMomentRpc(ctrl)
	mockSystemRpc := mock.NewMockSystemRpc(ctrl)
	mockCommentRpc := mock.NewMockCommentRpc(ctrl)
	mockPostRpc := mock.NewMockPostRpc(ctrl)

	svcCtx := &svc.ServiceContext{
		Config:        config.Config{},
		CollectionRPC: mockCollectionRpc,
		MomentRPC:     mockMomentRpc,
		SystemRPC:     mockSystemRpc,
		CommentRPC:    mockCommentRpc,
		PostRPC:       mockPostRpc,
	}
	l := NewAllowLogic(context.Background(), svcCtx)

	Convey("允许读", t, func() {
		allow, _ := l.Allow(&pb2.AllowReq{
			Object: ObjectPost,
			Action: ActionRead,
		})
		So(allow.Allow, ShouldBeTrue)
	})

	Convey("不允许非创建人的普通用户写", t, func() {
		mockPostRpc.EXPECT().RetrievePost(Any(), Any()).Return(&pb3.RetrievePostResp{
			Post: &pb3.Post{
				Id:     "PostId",
				UserId: "AnotherPostUserId2",
			},
		}, nil)
		mockSystemRpc.EXPECT().RetrieveUserRole(Any(), Any()).Return(&pb.RetrieveUserRoleResp{}, nil)
		allow, _ := l.Allow(&pb2.AllowReq{
			UserId:   "PoserUserId",
			Object:   ObjectPost,
			ObjectId: "NoticeId",
			Action:   ActionWrite,
		})
		So(allow.Allow, ShouldBeFalse)
	})

	Convey("允许超级管理员", t, func() {
		mockSystemRpc.EXPECT().RetrieveUserRole(Any(), Any(), Any()).Return(&pb.RetrieveUserRoleResp{
			Roles: []*pb.Role{
				{
					Type: RoleSuperAdmin,
				},
			},
		}, nil)
		allow, _ := l.Allow(&pb2.AllowReq{
			Object: ObjectPost,
			Action: ActionWrite,
		})
		So(allow.Allow, ShouldBeTrue)
	})

	Convey("允许帖子发布者", t, func() {
		mockPostRpc.EXPECT().RetrievePost(Any(), Any()).Return(&pb3.RetrievePostResp{
			Post: &pb3.Post{
				Id:     "PostId",
				UserId: "PostUserId",
			},
		}, nil)
		mockSystemRpc.EXPECT().RetrieveUserRole(Any(), Any(), Any()).Return(&pb.RetrieveUserRoleResp{}, nil)
		allow, _ := l.Allow(&pb2.AllowReq{
			UserId:   "PostUserId",
			Object:   ObjectPost,
			ObjectId: "PostId",
			Action:   ActionWrite,
		})
		So(allow.Allow, ShouldBeTrue)
	})

}

func TestAllowLogic_Allow_Comment(t *testing.T) {
	ctrl := NewController(t)
	defer ctrl.Finish()

	mockCollectionRpc := mock.NewMockCollectionRpc(ctrl)
	mockMomentRpc := mock.NewMockMomentRpc(ctrl)
	mockSystemRpc := mock.NewMockSystemRpc(ctrl)
	mockCommentRpc := mock.NewMockCommentRpc(ctrl)
	mockPostRpc := mock.NewMockPostRpc(ctrl)

	svcCtx := &svc.ServiceContext{
		Config:        config.Config{},
		CollectionRPC: mockCollectionRpc,
		MomentRPC:     mockMomentRpc,
		SystemRPC:     mockSystemRpc,
		CommentRPC:    mockCommentRpc,
		PostRPC:       mockPostRpc,
	}
	l := NewAllowLogic(context.Background(), svcCtx)

	Convey("允许读", t, func() {
		allow, _ := l.Allow(&pb2.AllowReq{
			Object: ObjectComment,
			Action: ActionRead,
		})
		So(allow.Allow, ShouldBeTrue)
	})

	Convey("不允许非创建人的普通用户写", t, func() {
		mockCommentRpc.EXPECT().RetrieveCommentById(Any(), Any()).Return(&pb4.RetrieveCommentByIdResponse{
			Comment: &pb4.Comment{
				Id:       "CommentId",
				AuthorId: "Another2CommentAuthorId",
			},
		}, nil)
		mockSystemRpc.EXPECT().RetrieveUserRole(Any(), Any()).Return(&pb.RetrieveUserRoleResp{}, nil)
		allow, _ := l.Allow(&pb2.AllowReq{
			UserId:   "CommentAuthorId",
			Object:   ObjectComment,
			ObjectId: "CommentId",
			Action:   ActionWrite,
		})
		So(allow.Allow, ShouldBeFalse)
	})

	Convey("允许超级管理员", t, func() {
		mockSystemRpc.EXPECT().RetrieveUserRole(Any(), Any(), Any()).Return(&pb.RetrieveUserRoleResp{
			Roles: []*pb.Role{
				{
					Type: RoleSuperAdmin,
				},
			},
		}, nil)
		allow, _ := l.Allow(&pb2.AllowReq{
			Object: ObjectComment,
			Action: ActionWrite,
		})
		So(allow.Allow, ShouldBeTrue)
	})

	Convey("允许评论发布者", t, func() {
		mockCommentRpc.EXPECT().RetrieveCommentById(Any(), Any()).Return(&pb4.RetrieveCommentByIdResponse{
			Comment: &pb4.Comment{
				Id:       "CommentId",
				AuthorId: "CommentAuthorId",
			},
		}, nil)
		mockSystemRpc.EXPECT().RetrieveUserRole(Any(), Any(), Any()).Return(&pb.RetrieveUserRoleResp{}, nil)
		allow, _ := l.Allow(&pb2.AllowReq{
			UserId:   "CommentAuthorId",
			Object:   ObjectComment,
			ObjectId: "CommentId",
			Action:   ActionWrite,
		})
		So(allow.Allow, ShouldBeTrue)
	})

	Convey("帖子发布者操作评论", t, func() {
		mockCommentRpc.EXPECT().RetrieveCommentById(Any(), Any()).Return(&pb4.RetrieveCommentByIdResponse{
			Comment: &pb4.Comment{
				Id:       "CommentId",
				Type:     "post",
				ParentId: "PostId",
				AuthorId: "Another2CommentAuthorId",
			},
		}, nil)
		mockPostRpc.EXPECT().RetrievePost(Any(), Any()).Return(&pb3.RetrievePostResp{
			Post: &pb3.Post{
				Id:     "PostId",
				UserId: "CommentAuthorId",
			},
		}, nil)
		mockSystemRpc.EXPECT().RetrieveUserRole(Any(), Any()).Times(2).Return(&pb.RetrieveUserRoleResp{}, nil)
		allow, _ := l.Allow(&pb2.AllowReq{
			UserId:   "CommentAuthorId",
			Object:   ObjectComment,
			ObjectId: "CommentId",
			Action:   ActionWrite,
		})
		So(allow.Allow, ShouldBeTrue)
	})

	Convey("非帖子发布者操作评论", t, func() {
		mockCommentRpc.EXPECT().RetrieveCommentById(Any(), Any()).Return(&pb4.RetrieveCommentByIdResponse{
			Comment: &pb4.Comment{
				Id:       "CommentId",
				Type:     "post",
				ParentId: "PostId",
				AuthorId: "Another2CommentAuthorId",
			},
		}, nil)
		mockPostRpc.EXPECT().RetrievePost(Any(), Any()).Return(&pb3.RetrievePostResp{
			Post: &pb3.Post{
				Id:     "PostId",
				UserId: "NotCommentAuthorId",
			},
		}, nil)
		mockSystemRpc.EXPECT().RetrieveUserRole(Any(), Any()).Times(2).Return(&pb.RetrieveUserRoleResp{}, nil)
		allow, _ := l.Allow(&pb2.AllowReq{
			UserId:   "CommentAuthorId",
			Object:   ObjectComment,
			ObjectId: "CommentId",
			Action:   ActionWrite,
		})
		So(allow.Allow, ShouldBeFalse)
	})

	Convey("动态发布者操作评论", t, func() {
		mockCommentRpc.EXPECT().RetrieveCommentById(Any(), Any()).Return(&pb4.RetrieveCommentByIdResponse{
			Comment: &pb4.Comment{
				Id:       "CommentId",
				Type:     "moment",
				ParentId: "MomentId",
				AuthorId: "=============",
			},
		}, nil)
		mockMomentRpc.EXPECT().RetrieveMoment(Any(), Any()).Return(&pb5.RetrieveMomentResp{
			Moment: &pb5.Moment{
				Id:     "MomentId",
				UserId: "CommentAuthorId",
			},
		}, nil)
		mockSystemRpc.EXPECT().RetrieveUserRole(Any(), Any()).Return(&pb.RetrieveUserRoleResp{}, nil)
		allow, _ := l.Allow(&pb2.AllowReq{
			UserId:   "CommentAuthorId",
			Object:   ObjectComment,
			ObjectId: "CommentId",
			Action:   ActionWrite,
		})
		So(allow.Allow, ShouldBeTrue)
	})

	Convey("非动态发布者操作评论", t, func() {
		mockCommentRpc.EXPECT().RetrieveCommentById(Any(), Any()).Return(&pb4.RetrieveCommentByIdResponse{
			Comment: &pb4.Comment{
				Id:       "CommentId",
				Type:     "moment",
				ParentId: "MomentId",
				AuthorId: "Another2CommentAuthorId",
			},
		}, nil)
		mockMomentRpc.EXPECT().RetrieveMoment(Any(), Any()).Return(&pb5.RetrieveMomentResp{
			Moment: &pb5.Moment{
				Id:     "MomentId",
				UserId: "NotMomentUserId",
			},
		}, nil)
		mockSystemRpc.EXPECT().RetrieveUserRole(Any(), Any()).Times(2).Return(&pb.RetrieveUserRoleResp{}, nil)
		allow, _ := l.Allow(&pb2.AllowReq{
			UserId:   "CommentAuthorId",
			Object:   ObjectComment,
			ObjectId: "CommentId",
			Action:   ActionWrite,
		})
		So(allow.Allow, ShouldBeFalse)
	})

	Convey("动态发布者的社区管理员操作评论", t, func() {
		mockCommentRpc.EXPECT().RetrieveCommentById(Any(), Any()).Return(&pb4.RetrieveCommentByIdResponse{
			Comment: &pb4.Comment{
				Id:       "CommentId",
				Type:     "moment",
				ParentId: "MomentId",
				AuthorId: "Another2CommentAuthorId",
			},
		}, nil)
		mockMomentRpc.EXPECT().RetrieveMoment(Any(), Any()).Return(&pb5.RetrieveMomentResp{
			Moment: &pb5.Moment{
				Id:          "MomentId",
				UserId:      "==============",
				CommunityId: "CommId",
			},
		}, nil)
		mockSystemRpc.EXPECT().RetrieveUserRole(Any(), Any()).Return(&pb.RetrieveUserRoleResp{}, nil)
		mockSystemRpc.EXPECT().RetrieveUserRole(Any(), Any()).Return(&pb.RetrieveUserRoleResp{
			Roles: []*pb.Role{
				{
					Type:        RoleCommunityAdmin,
					CommunityId: "CommId",
				},
			},
		}, nil)
		allow, _ := l.Allow(&pb2.AllowReq{
			UserId:   "CommentAuthorId",
			Object:   ObjectComment,
			ObjectId: "CommentId",
			Action:   ActionWrite,
		})
		So(allow.Allow, ShouldBeTrue)
	})

}
