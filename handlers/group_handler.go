package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/liyufei916/footballPaul/models"
	"github.com/liyufei916/footballPaul/services"
)

type GroupHandler struct {
	groupService            *services.GroupService
	leaderboardService      *services.GroupLeaderboardService
}

func NewGroupHandler() *GroupHandler {
	return &GroupHandler{
		groupService:       services.NewGroupService(),
		leaderboardService: services.NewGroupLeaderboardService(),
	}
}

// CreateGroup creates a new group
// POST /api/groups
func (h *GroupHandler) CreateGroup(c *gin.Context) {
	userID, _ := c.Get("userID")

	var req struct {
		Name string `json:"name" binding:"required,min=1,max=50"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "组名必填，长度1-50字符"})
		return
	}

	group, err := h.groupService.CreateGroup(req.Name, userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建组队失败"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"group":   group.ToResponse(),
		"message": "组队创建成功",
	})
}

// GetMyGroups returns all groups the current user belongs to
// GET /api/groups
func (h *GroupHandler) GetMyGroups(c *gin.Context) {
	userID, _ := c.Get("userID")

	groups, err := h.groupService.GetGroupsByUserID(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取组队列表失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"groups": groups})
}

// GetGroup returns a group's details
// GET /api/groups/:id
func (h *GroupHandler) GetGroup(c *gin.Context) {
	userID, _ := c.Get("userID")
	groupID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的组队ID"})
		return
	}

	if !h.groupService.IsGroupMember(uint(groupID), userID.(uint)) {
		c.JSON(http.StatusForbidden, gin.H{"error": "你不是该组成员，无法查看"})
		return
	}

	group, err := h.groupService.GetGroupByID(uint(groupID))
	if err != nil {
		if err == services.ErrGroupNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "组队不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取组队信息失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"group": group.ToResponse()})
}

// JoinGroup joins a group via invite code
// POST /api/groups/join
func (h *GroupHandler) JoinGroup(c *gin.Context) {
	userID, _ := c.Get("userID")

	var req struct {
		InviteCode string `json:"invite_code" binding:"required,len=6"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "邀请码必填，6位字符"})
		return
	}

	group, err := h.groupService.JoinGroup(req.InviteCode, userID.(uint))
	if err != nil {
		if err == services.ErrInvalidInviteCode {
			c.JSON(http.StatusBadRequest, gin.H{"error": "邀请码无效"})
			return
		}
		if err == services.ErrAlreadyInGroup {
			c.JSON(http.StatusConflict, gin.H{"error": "你已经在该组中"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "加入组队失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "加入成功",
		"group":   group.ToResponse(),
	})
}

// LeaveGroup leaves a group
// DELETE /api/groups/:id/leave
func (h *GroupHandler) LeaveGroup(c *gin.Context) {
	userID, _ := c.Get("userID")
	groupID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的组队ID"})
		return
	}

	err = h.groupService.LeaveGroup(uint(groupID), userID.(uint))
	if err != nil {
		if err == services.ErrGroupNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "组队不存在"})
			return
		}
		if err == services.ErrCannotLeaveAsOwner {
			c.JSON(http.StatusBadRequest, gin.H{"error": "组长无法离开，请先转让组长或解散该组"})
			return
		}
		if err == services.ErrNotGroupMember {
			c.JSON(http.StatusForbidden, gin.H{"error": "你不是该组成员"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "离开组队失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "已离开该组"})
}

// DeleteGroup deletes a group (owner only)
// DELETE /api/groups/:id
func (h *GroupHandler) DeleteGroup(c *gin.Context) {
	userID, _ := c.Get("userID")
	groupID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的组队ID"})
		return
	}

	err = h.groupService.DeleteGroup(uint(groupID), userID.(uint))
	if err != nil {
		if err == services.ErrGroupNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "组队不存在"})
			return
		}
		if err == services.ErrNotGroupOwner {
			c.JSON(http.StatusForbidden, gin.H{"error": "只有组长可以解散该组"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "解散组队失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "组队已解散"})
}

// GetMembers returns all members of a group
// GET /api/groups/:id/members
func (h *GroupHandler) GetMembers(c *gin.Context) {
	userID, _ := c.Get("userID")
	groupID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的组队ID"})
		return
	}

	if !h.groupService.IsGroupMember(uint(groupID), userID.(uint)) {
		c.JSON(http.StatusForbidden, gin.H{"error": "你不是该组成员"})
		return
	}

	members, err := h.groupService.GetGroupMembers(uint(groupID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取成员列表失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"members": members})
}

// GetCompetitions returns all competitions tracked by a group
// GET /api/groups/:id/competitions
func (h *GroupHandler) GetCompetitions(c *gin.Context) {
	userID, _ := c.Get("userID")
	groupID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的组队ID"})
		return
	}

	if !h.groupService.IsGroupMember(uint(groupID), userID.(uint)) {
		c.JSON(http.StatusForbidden, gin.H{"error": "你不是该组成员"})
		return
	}

	competitions, err := h.groupService.GetGroupCompetitions(uint(groupID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取赛事列表失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"competitions": competitions})
}

// AddCompetition adds a competition to a group
// POST /api/groups/:id/competitions
func (h *GroupHandler) AddCompetition(c *gin.Context) {
	userID, _ := c.Get("userID")
	groupID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的组队ID"})
		return
	}

	if !h.groupService.IsGroupAdmin(uint(groupID), userID.(uint)) {
		c.JSON(http.StatusForbidden, gin.H{"error": "只有管理员可以添加赛事"})
		return
	}

	var req struct {
		CompetitionID uint `json:"competition_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "competition_id 必填"})
		return
	}

	err = h.groupService.AddCompetitionToGroup(uint(groupID), req.CompetitionID)
	if err != nil {
		if err == services.ErrCompetitionNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "赛事不存在"})
			return
		}
		if err == services.ErrCompetitionExists {
			c.JSON(http.StatusConflict, gin.H{"error": "该赛事已在组内"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "添加赛事失败"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "赛事已添加到本组"})
}

// RemoveCompetition removes a competition from a group
// DELETE /api/groups/:id/competitions/:competitionId
func (h *GroupHandler) RemoveCompetition(c *gin.Context) {
	userID, _ := c.Get("userID")
	groupID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的组队ID"})
		return
	}
	competitionID, err := strconv.ParseUint(c.Param("competitionId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的赛事ID"})
		return
	}

	if !h.groupService.IsGroupAdmin(uint(groupID), userID.(uint)) {
		c.JSON(http.StatusForbidden, gin.H{"error": "只有管理员可以移除赛事"})
		return
	}

	err = h.groupService.RemoveCompetitionFromGroup(uint(groupID), uint(competitionID))
	if err != nil {
		if err == services.ErrCompetitionNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "赛事不在该组内"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "移除赛事失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "赛事已从本组移除"})
}

// GetLeaderboard returns the leaderboard for a specific competition within a group
// GET /api/groups/:id/leaderboard/:competitionId
func (h *GroupHandler) GetLeaderboard(c *gin.Context) {
	userID, _ := c.Get("userID")
	groupID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的组队ID"})
		return
	}
	competitionID, err := strconv.ParseUint(c.Param("competitionId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的赛事ID"})
		return
	}

	limit := 50
	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 {
			limit = parsed
		}
	}

	if !h.groupService.IsGroupMember(uint(groupID), userID.(uint)) {
		c.JSON(http.StatusForbidden, gin.H{"error": "你不是该组成员"})
		return
	}

	entries, err := h.leaderboardService.GetGroupLeaderboard(uint(groupID), uint(competitionID), limit)
	if err != nil {
		if err == services.ErrCompetitionNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "该赛事不在组内追踪列表中"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取排行榜失败"})
		return
	}

	// Get competition info
	var comp models.Competition
	comps, _ := h.groupService.GetGroupCompetitions(uint(groupID))
	for _, groupComp := range comps {
		if groupComp.ID == uint(competitionID) {
			comp = groupComp
			break
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"competition": comp,
		"leaderboard": entries,
	})
}

// GetGroupCompetitionPredictions returns all predictions from group members for a competition
// GET /api/groups/:id/competitions/:competitionId/predictions
func (h *GroupHandler) GetGroupCompetitionPredictions(c *gin.Context) {
	userID, _ := c.Get("userID")
	groupID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的组队ID"})
		return
	}
	competitionID, err := strconv.ParseUint(c.Param("competitionId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的赛事ID"})
		return
	}

	if !h.groupService.IsGroupMember(uint(groupID), userID.(uint)) {
		c.JSON(http.StatusForbidden, gin.H{"error": "你不是该组成员"})
		return
	}

	predictions, err := h.groupService.GetGroupCompetitionPredictions(uint(groupID), uint(competitionID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取预测列表失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"predictions": predictions,
	})
}

// TransferOwnership transfers group ownership to another member
// PUT /api/groups/:id/transfer-owner
func (h *GroupHandler) TransferOwnership(c *gin.Context) {
	userID, _ := c.Get("userID")
	groupID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的组队ID"})
		return
	}

	var req struct {
		NewOwnerID uint `json:"new_owner_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "new_owner_id 必填"})
		return
	}

	err = h.groupService.TransferOwnership(uint(groupID), userID.(uint), req.NewOwnerID)
	if err != nil {
		if err == services.ErrGroupNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "组队不存在"})
			return
		}
		if err == services.ErrNotGroupOwner {
			c.JSON(http.StatusForbidden, gin.H{"error": "只有组长可以转让"})
			return
		}
		if err == services.ErrNotGroupMember {
			c.JSON(http.StatusBadRequest, gin.H{"error": "新组长必须是组成员"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "转让失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "组长已转让"})
}
