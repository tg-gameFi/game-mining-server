package dbs

import (
	"fmt"
	"game-mining-server/entities"
	"game-mining-server/utils"
	"gorm.io/gorm"
)

type Point struct {
	Uid                   int64 `gorm:"primaryKey;type:bigint" json:"uid"`        // point user id, unique
	CreatedAt             int64 `gorm:"autoCreateTime:milli" json:"createdAt"`    // created ts: 1670400478555
	UpdatedAt             int64 `gorm:"autoUpdateTime:milli" json:"-"`            // updated ts: 1670400478555
	LastClaimedPointValue int64 `gorm:"type:bigint" json:"lastClaimedPointValue"` // last time claimed point value
	TotalWalletPointValue int64 `gorm:"type:bigint" json:"totalWalletPointValue"` // total wallet related point value, no need record for last time
	LastInvitePointLevel  int64 `gorm:"type:bigint" json:"lastInvitePointLevel"`  // last time user claimed for invite level
	TotalInvitePointValue int64 `gorm:"type:bigint" json:"totalInvitePointValue"` // total invite point value
	TotalPointValue       int64 `gorm:"type:bigint" json:"totalPointValue"`       // total point value, totalPointValue = claimedPointValue * count + totalWalletPointValue + totalInvitePointValue
}

var pointWithUserQueryFields = `
points.uid AS uid,
points.created_at AS created_at, 
points.total_point_value AS total_point_value,
users.username AS username
`

type PointWithUser struct {
	Point
	Username string `json:"username,omitempty"`
}

func (u *Point) TableName() string {
	return "points"
}

// PointFindByUid find a user's point by uid
func (s *Service) PointFindByUid(uid int64) (*Point, error) {
	var point Point
	if e := s.DBInstance.Where("uid = ?", uid).First(&point).Error; e != nil {
		return nil, e
	} else {
		return &point, nil
	}
}

// PointClaimTask Claim a point value for specified user, if not exists, create it
func (s *Service) PointClaimTask(db *gorm.DB, uid int64, claimed int64) (*Point, error) {
	var point Point
	if e0 := db.Where(Point{Uid: uid}).Attrs(&Point{Uid: uid}).FirstOrCreate(&point).Error; e0 != nil {
		return nil, e0
	}
	point.LastClaimedPointValue = claimed
	point.TotalPointValue = point.TotalPointValue + claimed
	if e5 := db.Save(&point).Error; e5 != nil {
		return nil, e5
	} else {
		return &point, nil
	}
}

func (s *Service) PointClaimForWallet(uid int64, walletPoint int64) (*Point, error) {
	var point Point
	e := s.DBInstance.Transaction(func(tx *gorm.DB) error {
		if e0 := tx.Where(Point{Uid: uid}).Attrs(&Point{Uid: uid}).FirstOrCreate(&point).Error; e0 != nil {
			return e0
		}
		point.TotalWalletPointValue = point.TotalWalletPointValue + walletPoint
		point.TotalPointValue = point.TotalPointValue + walletPoint
		return tx.Save(&point).Error
	})
	return &point, e
}

func (s *Service) PointClaimForInvite(uid int64, level int64) (*Point, error) {
	var point Point
	e := s.DBInstance.Transaction(func(tx *gorm.DB) error {
		if e0 := tx.Where(Point{Uid: uid}).Attrs(&Point{Uid: uid}).FirstOrCreate(&point).Error; e0 != nil {
			return e0
		}

		if level <= point.LastInvitePointLevel {
			// already claimed
			return fmt.Errorf("already claimed")
		}

		var inviteCount int64
		if e1 := s.DBInstance.Model(&User{}).Where("referral_uid = ?", uid).Count(&inviteCount).Error; e1 != nil {
			return e1
		}

		invitePoint := utils.CalPointForInvite(point.LastInvitePointLevel, level, inviteCount)
		if invitePoint <= 0 {
			return fmt.Errorf("not enough invites: %d or not supported level: %d", inviteCount, level)
		}

		point.LastInvitePointLevel = level
		point.TotalInvitePointValue = point.TotalInvitePointValue + invitePoint
		point.TotalPointValue = point.TotalPointValue + invitePoint
		return tx.Save(&point).Error
	})
	return &point, e
}

func (s *Service) PointGetLeaderBoardWithUser(params *entities.LeaderBoardParam) ([]*PointWithUser, int64, error) {
	var points []*PointWithUser
	var total int64
	result := s.DBInstance.Model(&Point{}).Offset(-1).Limit(-1).Count(&total).
		Select(pointWithUserQueryFields).Joins("left join users on users.id = points.uid").
		Offset(params.Offset).Limit(params.Limit).Order("points.total_point_value desc").Scan(&points)

	if result.Error != nil {
		return nil, 0, result.Error
	} else {
		return points, total, nil
	}
}
