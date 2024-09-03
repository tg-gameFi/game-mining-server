package dbs

import (
	"fmt"
	"game-mining-server/entities"
	"gorm.io/gorm"
	"time"
)

type RewardLog struct {
	Id          int64     `gorm:"primaryKey"`
	MomentId    int64     `gorm:"not null"`
	FromUserId  int64     `gorm:"not null"`
	ToUserId    int64     `gorm:"not null"`
	Amount      int       `gorm:"not null"`
	Description string    `gorm:"column:descr"`
	CreatedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP"`
}

func (u *RewardLog) TableName() string {
	return "reward_logs"
}

func (s *Service) RewardMoment(params entities.RewardMomentParam) error {
	return s.DBInstance.Transaction(func(tx *gorm.DB) error {
		// Check if fromUserID has enough points
		var fromUser User
		if err := tx.First(&fromUser, params.FromUserId).Error; err != nil {
			return fmt.Errorf("user not found, %d", params.FromUserId)
		}

		if fromUser.RewardPoints < params.Amount {
			return fmt.Errorf("user dont have enough reward points, %d", params.FromUserId)
		}

		// Create reward log
		rewardLog := RewardLog{
			MomentId:    params.MomentId,
			FromUserId:  params.FromUserId,
			ToUserId:    params.ToUserId,
			Amount:      params.Amount,
			Description: "Moment reward",
		}
		if err := tx.Create(&rewardLog).Error; err != nil {
			return err
		}

		// Update fromUser's reward points
		if err := tx.Model(&User{}).Where("id = ?", params.FromUserId).UpdateColumn("reward_points", gorm.Expr("reward_points - ?", params.Amount)).Error; err != nil {
			return err
		}

		// Update toUser's reward points
		if err := tx.Model(&User{}).Where("id = ?", params.ToUserId).UpdateColumn("reward_points", gorm.Expr("reward_points + ?", params.Amount)).Error; err != nil {
			return err
		}

		// Update moment's rewards count (now as total amount)
		if err := tx.Model(&Moment{}).Where("id = ?", params.MomentId).UpdateColumn("rewards_amount", gorm.Expr("rewards_amount + ?", params.Amount)).Error; err != nil {
			return err
		}

		return nil
	})
}

func (s *Service) RecordDailyRefresh(userID int64) error {
	rewardLog := RewardLog{
		MomentId:    0,
		FromUserId:  0,
		ToUserId:    userID,
		Amount:      200,
		Description: "Daily reward points refresh",
	}
	return s.DBInstance.Create(&rewardLog).Error
}
