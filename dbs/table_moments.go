package dbs

import (
	"errors"
	"fmt"
	"game-mining-server/entities"
	"gorm.io/gorm"
	"time"
)

type Moment struct {
	Id            int64  `gorm:"primaryKey;autoIncrement"`
	UserId        int64  `gorm:"not null"`
	Content       string `gorm:"type:text"`
	ImageURL      string
	LikesCount    int       `gorm:"default:0"`
	RewardsAmount int       `gorm:"default:0"`
	Comments      []Comment `gorm:"foreignKey:MomentId"`
	CreatedAt     time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt     time.Time `gorm:"column:updated_at"`
}

func (u *Moment) TableName() string {
	return "moments"
}

func (s *Service) CreateMoment(params entities.CreateMomentParam) (int64, error) {
	moment := Moment{
		UserId:   params.UserId,
		Content:  params.Content,
		ImageURL: params.ImageURL,
	}

	result := s.DBInstance.Create(&moment)
	if result.Error != nil {
		return 0, result.Error
	}
	return moment.Id, nil
}

func (s *Service) DeleteMoment(params entities.DeleteMomentParam) error {
	return s.DBInstance.Transaction(func(tx *gorm.DB) error {
		// Check if the moment exists and belongs to the user
		var moment Moment
		if err := tx.First(&moment, params.MomentId).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("moment not found")
			}
			return err
		}

		if moment.UserId != params.UserId {
			return fmt.Errorf("unauthorized to delete this moment")
		}

		// Delete comments
		if err := tx.Where("moment_id = ?", params.MomentId).Delete(&Comment{}).Error; err != nil {
			return err
		}

		// Delete likes
		if err := tx.Where("moment_id = ?", params.MomentId).Delete(&Like{}).Error; err != nil {
			return err
		}

		// Delete the moment
		if err := tx.Delete(&moment).Error; err != nil {
			return err
		}

		return nil
	})
}

// GetLatestMoments with Moments，include comments, likes and rewards
func (s *Service) GetLatestMoments(params entities.GetLatestMomentsParam) ([]Moment, error) {
	var moments []Moment
	err := s.DBInstance.
		Preload("Comments", func(db *gorm.DB) *gorm.DB {
			return db.Order("comments.created_at DESC").Limit(10) // load latest 10 comments
		}).
		Order("created_at DESC").
		Limit(params.Limit).
		Offset(params.Offset).
		Find(&moments).Error

	if err != nil {
		return nil, err
	}

	return moments, nil
}
