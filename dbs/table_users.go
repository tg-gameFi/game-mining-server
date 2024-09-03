package dbs

import (
	"game-mining-server/entities"
	"time"
)

type User struct {
	Id           int64  `redis:"id" gorm:"primaryKey;type:bigint" json:"id"`                // user id, generate by Telegram, which means Telegram user id
	CreatedAt    int64  `redis:"ct" gorm:"autoCreateTime:milli" json:"createdAt"`           // created ts: 1670400478555
	UpdatedAt    int64  `redis:"ut" gorm:"autoUpdateTime:milli" json:"-"`                   // updated ts: 1670400478555
	Username     string `redis:"un" gorm:"type:varchar(255)" json:"username"`               // username in Telegram
	IsPremium    bool   `redis:"ip" gorm:"type:bool" json:"isPremium"`                      // is premium user
	ReferralCode string `redis:"rc" gorm:"type:varchar(255)" json:"referralCode"`           // current user's referral code, generate when user first login
	LanguageCode string `redis:"lc" gorm:"type:varchar(255)" json:"languageCode,omitempty"` // user language code
	ReferralUid  int64  `redis:"ru" gorm:"type:bigint" json:"referralUid,omitempty"`        // current user referral user id, which who has referral current user
	// for moments
	Moments  []Moment  `gorm:"foreignKey:UserId"`
	Comments []Comment `gorm:"foreignKey:UserId"`
	Likes    []Like    `gorm:"foreignKey:UserId"`
	// reward points per user everyday
	RewardPoints      int       `gorm:"default:200"`
	LastPointsRefresh time.Time `gorm:"default:CURRENT_TIMESTAMP"`
}

func (u *User) TableName() string {
	return "users"
}

// UserInsert add a user info
func (s *Service) UserInsert(user *User) error {
	return s.DBInstance.Create(&user).Error
}

// UserFindById find user by id
func (s *Service) UserFindById(id int64) (*User, error) {
	var user User
	if e := s.DBInstance.Where("id = ?", id).First(&user).Error; e != nil {
		return nil, e
	} else {
		return &user, nil
	}
}

// UserFindByReferralCode find user by referral code
func (s *Service) UserFindByReferralCode(referralCode string) (*User, error) {
	var user User
	if e := s.DBInstance.Where("referral_code = ?", referralCode).First(&user).Error; e != nil {
		return nil, e
	} else {
		return &user, nil
	}
}

// UserUpdateFields Update user fields
func (s *Service) UserUpdateFields(id int64, updated map[string]interface{}) error {
	return s.DBInstance.Model(&User{}).Where("id = ?", id).Updates(updated).Error
}

// UserFindInvitedUserList return a user's invited user list
func (s *Service) UserFindInvitedUserList(uid int64, params *entities.InvitedUserListParam) ([]*User, int64, error) {
	var users []*User
	var total int64
	result := s.DBInstance.Model(&User{}).Select("id", "created_at", "username").Where("referral_uid = ?", uid).
		Offset(-1).Limit(-1).Count(&total).
		Offset(params.Offset).Limit(params.Limit).Order("created_at desc").
		Find(&users)
	if result.Error != nil {
		return nil, 0, result.Error
	} else {
		return users, total, nil
	}
}

// UserCountInvitedUsers return a count of user's invited users
func (s *Service) UserCountInvitedUsers(uid int64) (int64, error) {
	var total int64
	result := s.DBInstance.Model(&User{}).Where("referral_uid = ?", uid).Offset(-1).Limit(-1).Count(&total)
	if result.Error != nil {
		return 0, result.Error
	} else {
		return total, nil
	}
}

// RefreshUserRewardPoints refresh user's reward points everyday
func (s *Service) RefreshUserRewardPoints(userID int64) error {
	var user User
	if err := s.DBInstance.First(&user, userID).Error; err != nil {
		return err
	}

	// get current UTC time
	currentDate := time.Now().UTC().Truncate(24 * time.Hour)
	lastRefreshDate := user.LastPointsRefresh.UTC().Truncate(24 * time.Hour)

	// If it has been a day since the last refresh time
	if currentDate.After(lastRefreshDate) {
		// update user reward points
		err := s.DBInstance.Model(&user).Updates(map[string]interface{}{
			"RewardPoints":      200,
			"LastPointsRefresh": time.Now().UTC(),
		}).Error
		if err != nil {
			return err
		}

		// record reward points change
		return s.RecordDailyRefresh(userID)
	}

	return nil
}
