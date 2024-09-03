package dbs

import (
	"errors"
	"fmt"
	"game-mining-server/entities"
	"gorm.io/gorm"
	"time"
)

type Comment struct {
	Id            int64     `gorm:"primaryKey;autoIncrement"`
	UserId        int64     `gorm:"not null"`
	MomentId      int64     `gorm:"not null"`
	Content       string    `gorm:"type:text"`
	ReplyToUserId *int64    `gorm:"default:null"` // 引用的用户ID，可以为空
	CreatedAt     time.Time `gorm:"default:CURRENT_TIMESTAMP"`
}

type Like struct {
	Id        int64     `gorm:"primaryKey"`
	MomentId  int64     `gorm:"not null"`
	UserId    int64     `gorm:"not null"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
}

func (u *Comment) TableName() string {
	return "comments"
}

func (u *Like) TableName() string {
	return "likes"
}

func (s *Service) CommentMoment(params entities.AddCommentParam) (int64, error) {
	comment := Comment{
		UserId:        params.UserId,
		MomentId:      params.MomentId,
		Content:       params.Content,
		ReplyToUserId: params.ReplyToUserId,
	}

	err := s.DBInstance.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&comment).Error; err != nil {
			return err
		}

		// update Moment cnt
		if err := tx.Model(&Moment{}).Where("id = ?", params.MomentId).UpdateColumn("comments_count", gorm.Expr("comments_count + ?", 1)).Error; err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return 0, err
	}
	return comment.Id, nil
}

// DeleteComment ...
func (s *Service) DeleteComment(params entities.DeleteCommentParam) error {
	return s.DBInstance.Transaction(func(tx *gorm.DB) error {
		var comment Comment
		if err := tx.First(&comment, params.CommentId).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("comment not found")
			}
			return err
		}

		var moment Moment
		if err := tx.First(&moment, comment.MomentId).Error; err != nil {
			return err
		}

		// Only the author of the comment or the author of the Moment can delete the comment
		if comment.UserId != params.UserId && moment.UserId != params.UserId {
			return fmt.Errorf("unauthorized to delete this comment")
		}

		if err := tx.Delete(&comment).Error; err != nil {
			return err
		}

		// update Moment cnt
		if err := tx.Model(&Moment{}).Where("id = ?", comment.MomentId).UpdateColumn("comments_count", gorm.Expr("comments_count - ?", 1)).Error; err != nil {
			return err
		}

		return nil
	})
}

// GetCommentsForMoment for `more comments` bottom, by limit and offset
func (s *Service) GetCommentsForMoment(params entities.GetCommentsForMomentParam) ([]Comment, error) {
	var comments []Comment
	err := s.DBInstance.
		Where("moment_id = ?", params.MomentId).
		Order("created_at DESC").
		Limit(params.Limit).
		Offset(params.Offset).
		Find(&comments).Error

	if err != nil {
		return nil, err
	}
	return comments, nil
}

func (s *Service) LikeMoment(params entities.LikeMomentParam) error {
	like := Like{
		UserId:   params.UserId,
		MomentId: params.MomentId,
	}
	if err := s.DBInstance.Create(&like).Error; err != nil {
		return err
	}

	// update Moment like cnt
	return s.DBInstance.Model(&Moment{}).Where("id = ?", params.MomentId).Update("likes_count", gorm.Expr("likes_count + ?", 1)).Error
}

func (s *Service) RollbackLikeMoment(params entities.LikeMomentParam) error {
	return s.DBInstance.Transaction(func(tx *gorm.DB) error {
		// Check if the like exists
		var like Like
		if err := tx.Where("user_id = ? AND moment_id = ?", params.UserId, params.MomentId).First(&like).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("like not found")
			}
			return err
		}

		// Delete the like
		if err := tx.Delete(&like).Error; err != nil {
			return err
		}

		// Update Moment like count
		var moment Moment
		if err := tx.First(&moment, params.MomentId).Error; err != nil {
			return err
		}

		// Update likes_count, ensuring it doesn't go below zero
		if moment.LikesCount > 0 {
			if err := tx.Model(&Moment{}).Where("id = ?", params.MomentId).UpdateColumn("likes_count", gorm.Expr("likes_count - ?", 1)).Error; err != nil {
				return err
			}
		}

		return nil
	})
}
