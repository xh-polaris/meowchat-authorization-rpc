package logic

import (
	"context"
	. "github.com/golang/mock/gomock"
	. "github.com/smartystreets/goconvey/convey"
	. "github.com/xh-polaris/meowchat-system-rpc/constant"
	"github.com/xh-polaris/meowchat-system-rpc/pb"
	. "meowchat-authorization-rpc/constant"
	"meowchat-authorization-rpc/internal/config"
	"meowchat-authorization-rpc/internal/logic/mock"
	"meowchat-authorization-rpc/internal/svc"
	pb2 "meowchat-authorization-rpc/pb"
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
		allow := l.allowCommunity(&pb2.AllowReq{
			Object: ObjectCommunity,
			Action: ActionRead,
		})
		So(allow, ShouldBeTrue)
	})

	Convey("允许超级管理员", t, func() {
		mockSystemRpc.EXPECT().RetrieveUserRole(Any(), Any(), Any()).Return(&pb.RetrieveUserRoleResp{
			Roles: []*pb.Role{
				{
					Type: RoleSuperAdmin,
				},
			},
		}, nil)
		allow := l.allowCommunity(&pb2.AllowReq{
			Object: ObjectCommunity,
			Action: ActionWrite,
		})
		So(allow, ShouldBeTrue)
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
		allow := l.allowCommunity(&pb2.AllowReq{
			Object:   ObjectCommunity,
			ObjectId: "TestCommId",
			Action:   ActionWrite,
		})
		So(allow, ShouldBeTrue)
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
		allow := l.allowCommunity(&pb2.AllowReq{
			Object:   ObjectCommunity,
			ObjectId: "TestCommId2",
			Action:   ActionWrite,
		})
		So(allow, ShouldBeFalse)
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
		allow := l.allowCommunity(&pb2.AllowReq{
			Object:   ObjectCommunity,
			ObjectId: "Child",
			Action:   ActionWrite,
		})
		So(allow, ShouldBeTrue)
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
		allow := l.allowCommunity(&pb2.AllowReq{
			Object:   ObjectCommunity,
			ObjectId: "Child",
			Action:   ActionWrite,
		})
		So(allow, ShouldBeTrue)
	})

}
