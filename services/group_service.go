package services

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/liyufei916/footballPaul/database"
	"github.com/liyufei916/footballPaul/models"
	"gorm.io/gorm"
)

const (
	inviteCodeCharset = "ABCDEFGHJKMNPQRSTUVWXYZ23456789" // exclude 0,O,1,I
	inviteCodeLen     = 6
)

var (
	ErrGroupNotFound      = errors.New("group not found")
	ErrNotGroupMember     = errors.New("user is not a member of this group")
	ErrAlreadyInGroup     = errors.New("user is already in this group")
	ErrNotGroupOwner      = errors.New("only the owner can perform this action")
	ErrNotGroupAdmin      = errors.New("only admins can perform this action")
	ErrInvalidInviteCode  = errors.New("invalid invite code")
	ErrCannotLeaveAsOwner = errors.New("owner cannot leave the group, transfer ownership or delete the group first")
	ErrCompetitionExists  = errors.New("competition is already being tracked in this group")
	ErrCompetitionNotFound = errors.New("competition not found in this group")
)

type GroupService struct{}

func NewGroupService() *GroupService {
	return &GroupService{}
}

// GenerateInviteCode generates a unique 6-character invite code
func (s *GroupService) GenerateInviteCode() (string, error) {
	result := make([]byte, inviteCodeLen)
	charsetLen := big.NewInt(int64(len(inviteCodeCharset)))
	for i := 0; i < inviteCodeLen; i++ {
		n, err := rand.Int(rand.Reader, charsetLen)
		if err != nil {
			return "", fmt.Errorf("failed to generate invite code: %w", err)
		}
		result[i] = inviteCodeCharset[n.Int64()]
	}

	code := string(result)

	// Ensure uniqueness
	var count int64
	database.DB.Model(&models.Group{}).Where("invite_code = ?", code).Count(&count)
	if count > 0 {
		return s.GenerateInviteCode() // regenerate
	}

	return code, nil
}

// CreateGroup creates a new group with the given name and creator as owner
func (s *GroupService) CreateGroup(name string, ownerID uint) (*models.Group, error) {
	inviteCode, err := s.GenerateInviteCode()
	if err != nil {
		return nil, err
	}

	group := &models.Group{
		Name:       name,
		InviteCode: inviteCode,
		OwnerID:    ownerID,
	}

	if err := database.DB.Create(group).Error; err != nil {
		return nil, fmt.Errorf("failed to create group: %w", err)
	}

	// Add creator as admin member
	member := &models.GroupMember{
		GroupID:  group.ID,
		UserID:   ownerID,
		Role:     "admin",
		JoinedAt: time.Now(),
	}
	if err := database.DB.Create(member).Error; err != nil {
		return nil, fmt.Errorf("failed to add owner as member: %w", err)
	}

	return group, nil
}

// GetGroupsByUserID returns all groups a user belongs to
func (s *GroupService) GetGroupsByUserID(userID uint) ([]models.GroupResponse, error) {
	var members []models.GroupMember
	if err := database.DB.Preload("Group").
		Where("user_id = ?", userID).
		Find(&members).Error; err != nil {
		return nil, err
	}

	responses := make([]models.GroupResponse, 0, len(members))
	for _, m := range members {
		g := m.Group

		// Count members and competitions
		var memberCount int64
		database.DB.Model(&models.GroupMember{}).Where("group_id = ?", g.ID).Count(&memberCount)

		var compCount int64
		database.DB.Model(&models.GroupCompetition{}).Where("group_id = ?", g.ID).Count(&compCount)

		resp := models.GroupResponse{
			ID:               g.ID,
			Name:             g.Name,
			OwnerID:          g.OwnerID,
			MemberCount:      int(memberCount),
			CompetitionCount: int(compCount),
			Role:             m.Role,
			JoinedAt:         m.JoinedAt,
			CreatedAt:        g.CreatedAt,
		}
		responses = append(responses, resp)
	}

	return responses, nil
}

// GetGroupByID returns a group by its ID
func (s *GroupService) GetGroupByID(groupID uint) (*models.Group, error) {
	var group models.Group
	if err := database.DB.First(&group, groupID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrGroupNotFound
		}
		return nil, err
	}
	return &group, nil
}

// GetGroupByInviteCode returns a group by its invite code
func (s *GroupService) GetGroupByInviteCode(code string) (*models.Group, error) {
	var group models.Group
	if err := database.DB.Where("invite_code = ?", code).First(&group).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrInvalidInviteCode
		}
		return nil, err
	}
	return &group, nil
}

// JoinGroup adds a user to a group via invite code
func (s *GroupService) JoinGroup(inviteCode string, userID uint) (*models.Group, error) {
	group, err := s.GetGroupByInviteCode(inviteCode)
	if err != nil {
		return nil, err
	}

	// Check if already a member
	var existing models.GroupMember
	err = database.DB.Where("group_id = ? AND user_id = ?", group.ID, userID).First(&existing).Error
	if err == nil {
		return nil, ErrAlreadyInGroup
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	member := &models.GroupMember{
		GroupID:  group.ID,
		UserID:   userID,
		Role:     "member",
		JoinedAt: time.Now(),
	}
	if err := database.DB.Create(member).Error; err != nil {
		return nil, fmt.Errorf("failed to join group: %w", err)
	}

	return group, nil
}

// LeaveGroup removes a user from a group
func (s *GroupService) LeaveGroup(groupID, userID uint) error {
	group, err := s.GetGroupByID(groupID)
	if err != nil {
		return err
	}

	if group.OwnerID == userID {
		return ErrCannotLeaveAsOwner
	}

	result := database.DB.Where("group_id = ? AND user_id = ?", groupID, userID).Delete(&models.GroupMember{})
	if result.RowsAffected == 0 {
		return ErrNotGroupMember
	}

	return nil
}

// DeleteGroup deletes a group (only owner can do this)
func (s *GroupService) DeleteGroup(groupID, userID uint) error {
	group, err := s.GetGroupByID(groupID)
	if err != nil {
		return err
	}

	if group.OwnerID != userID {
		return ErrNotGroupOwner
	}

	return database.DB.Transaction(func(tx *gorm.DB) error {
		// Delete group competitions
		if err := tx.Where("group_id = ?", groupID).Delete(&models.GroupCompetition{}).Error; err != nil {
			return err
		}
		// Delete group members
		if err := tx.Where("group_id = ?", groupID).Delete(&models.GroupMember{}).Error; err != nil {
			return err
		}
		// Delete group
		if err := tx.Delete(&models.Group{}, groupID).Error; err != nil {
			return err
		}
		return nil
	})
}

// IsGroupMember checks if a user is a member of a group
func (s *GroupService) IsGroupMember(groupID, userID uint) bool {
	var count int64
	database.DB.Model(&models.GroupMember{}).Where("group_id = ? AND user_id = ?", groupID, userID).Count(&count)
	return count > 0
}

// IsGroupAdmin checks if a user is an admin of a group
func (s *GroupService) IsGroupAdmin(groupID, userID uint) bool {
	var count int64
	database.DB.Model(&models.GroupMember{}).Where("group_id = ? AND user_id = ? AND role = ?", groupID, userID, "admin").Count(&count)
	return count > 0
}

// IsGroupOwner checks if a user is the owner of a group
func (s *GroupService) IsGroupOwner(groupID, userID uint) bool {
	group, err := s.GetGroupByID(groupID)
	if err != nil {
		return false
	}
	return group.OwnerID == userID
}

// GetMemberRole returns the role of a user in a group
func (s *GroupService) GetMemberRole(groupID, userID uint) (string, error) {
	var member models.GroupMember
	err := database.DB.Where("group_id = ? AND user_id = ?", groupID, userID).First(&member).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", ErrNotGroupMember
		}
		return "", err
	}
	return member.Role, nil
}

// GetGroupMembers returns all members of a group
func (s *GroupService) GetGroupMembers(groupID uint) ([]models.GroupMemberResponse, error) {
	var members []models.GroupMember
	if err := database.DB.Preload("User").Where("group_id = ?", groupID).Find(&members).Error; err != nil {
		return nil, err
	}

	responses := make([]models.GroupMemberResponse, 0, len(members))
	for _, m := range members {
		responses = append(responses, models.GroupMemberResponse{
			UserID:   m.UserID,
			Username: m.User.Username,
			Role:     m.Role,
			JoinedAt: m.JoinedAt,
		})
	}

	return responses, nil
}

// GetGroupCompetitions returns all competitions tracked by a group
func (s *GroupService) GetGroupCompetitions(groupID uint) ([]models.Competition, error) {
	var groupComps []models.GroupCompetition
	if err := database.DB.Preload("Competition").Where("group_id = ?", groupID).Find(&groupComps).Error; err != nil {
		return nil, err
	}

	competitions := make([]models.Competition, 0, len(groupComps))
	for _, gc := range groupComps {
		competitions = append(competitions, gc.Competition)
	}

	return competitions, nil
}

// AddCompetitionToGroup adds a competition to a group
func (s *GroupService) AddCompetitionToGroup(groupID, competitionID uint) error {
	// Verify competition exists
	var comp models.Competition
	if err := database.DB.First(&comp, competitionID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrCompetitionNotFound
		}
		return err
	}

	// Check if already added
	var existing models.GroupCompetition
	err := database.DB.Where("group_id = ? AND competition_id = ?", groupID, competitionID).First(&existing).Error
	if err == nil {
		return ErrCompetitionExists
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	gc := &models.GroupCompetition{
		GroupID:       groupID,
		CompetitionID: competitionID,
		CreatedAt:     time.Now(),
	}

	return database.DB.Create(gc).Error
}

// RemoveCompetitionFromGroup removes a competition from a group
func (s *GroupService) RemoveCompetitionFromGroup(groupID, competitionID uint) error {
	result := database.DB.Where("group_id = ? AND competition_id = ?", groupID, competitionID).Delete(&models.GroupCompetition{})
	if result.RowsAffected == 0 {
		return ErrCompetitionNotFound
	}
	return result.Error
}

// TransferOwnership transfers group ownership to another member
func (s *GroupService) TransferOwnership(groupID, currentUserID, newOwnerID uint) error {
	group, err := s.GetGroupByID(groupID)
	if err != nil {
		return err
	}

	if group.OwnerID != currentUserID {
		return ErrNotGroupOwner
	}

	// Verify new owner is a member
	if !s.IsGroupMember(groupID, newOwnerID) {
		return ErrNotGroupMember
	}

	return database.DB.Transaction(func(tx *gorm.DB) error {
		// Update group owner
		if err := tx.Model(&models.Group{}).Where("id = ?", groupID).Update("owner_id", newOwnerID).Error; err != nil {
			return err
		}
		// Update roles: new owner becomes admin, old owner stays admin
		if err := tx.Model(&models.GroupMember{}).
			Where("group_id = ? AND user_id = ?", groupID, newOwnerID).
			Update("role", "admin").Error; err != nil {
			return err
		}
		return nil
	})
}
