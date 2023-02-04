package logic

import (
	. "github.com/xh-polaris/meowchat-system-rpc/constant"
	system "github.com/xh-polaris/meowchat-system-rpc/pb"
)

// 判断用户是否包含某个角色
func (l *AllowLogic) containsRole(userId, role string) bool {
	userRole, _ := l.svcCtx.SystemRPC.RetrieveUserRole(l.ctx, &system.RetrieveUserRoleReq{UserId: userId})
	if userRole == nil || userRole.Roles == nil {
		return false
	}

	for _, r := range userRole.Roles {
		if r.Type == role {
			return true
		}
	}

	return false
}

// 判断cid1的社区是不是cid2的社区的子社区
func (l *AllowLogic) subCommunityOf(cid1, cid2 string) bool {
	if cid1 == cid2 {
		return true
	}
	c1, _ := l.svcCtx.SystemRPC.RetrieveCommunity(l.ctx, &system.RetrieveCommunityReq{Id: cid1})
	return c1 != nil && c1.Community.ParentId == cid2
}

// 判断userId对应用户是否是超级管理员或是某个社区的管理员
func (l *AllowLogic) allowCommunityOrSuperAdmin(userId, communityId string) bool {
	userRole, err := l.svcCtx.SystemRPC.RetrieveUserRole(l.ctx, &system.RetrieveUserRoleReq{UserId: userId})
	if err != nil || userRole == nil || userRole.Roles == nil {
		return false
	}

	for _, r := range userRole.Roles {
		if r.Type == RoleSuperAdmin ||
			(r.Type == RoleCommunityAdmin && l.subCommunityOf(communityId, r.CommunityId)) {
			return true
		}
	}

	return false
}
