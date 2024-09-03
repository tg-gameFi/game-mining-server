package api

import (
	"game-mining-server/app"
	"game-mining-server/caches"
	"game-mining-server/configs"
	"game-mining-server/dbs"
	"game-mining-server/entities"
	"game-mining-server/routers/middleware"
	"game-mining-server/utils"
	"github.com/gin-gonic/gin"
	initdata "github.com/telegram-mini-apps/init-data-golang"
	"gorm.io/gorm"
	"log"
	"net/http"
	"time"
)

type UserLoginRes struct {
	User    *dbs.User    `json:"user"`
	IsNew   bool         `json:"isNew"`
	Session string       `json:"session"`
	Secret  string       `json:"secret"`
	Checkin *dbs.Checkin `json:"checkin,omitempty"`
}

type UserInvitedUserListRes struct {
	Users []*dbs.User `json:"users"`
	Total int64       `json:"total"`
}

type UserLeaderBoardRes struct {
	Users []*dbs.PointWithUser `json:"users"`
	Total int64                `json:"total"`
}

// Login
// @Tags User
// @Router /user/login [post]
// @Summary User request to log in, return user info and session
// @description User login in by Telegram mini app initData
// @Accept json
// @Success 200 {object} UserLoginRes
func Login(c *gin.Context) {
	// valid request body
	var params entities.UserLoginParam
	if e0 := c.ShouldBindJSON(&params); e0 != nil {
		c.JSON(http.StatusBadRequest, entities.ResFailed(entities.ErrInvalidParams, e0.Error()))
		return
	}

	// only valid on PROD env
	if app.Config().Basic.Env == configs.EnvPROD && app.Bot() != nil {
		if e1 := initdata.Validate(params.InitDataRaw, app.Bot().Token(), 24*time.Hour); e1 != nil {
			c.JSON(http.StatusBadRequest, entities.ResFailed(entities.ErrInvalidInitData, e1.Error()))
			return
		}
	}

	userInitData, e2 := initdata.Parse(params.InitDataRaw)
	if e2 != nil {
		c.JSON(http.StatusBadRequest, entities.ResFailed(entities.ErrParseInitDataFailed, e2.Error()))
		return
	}
	uid := userInitData.User.ID
	user, isNew, checkin, e4 := userLoginAndCheckin(app.DB(), uid, params.Referral, params.RandPoint, &dbs.User{
		Id:                uid,
		Username:          userInitData.User.Username,
		IsPremium:         userInitData.User.IsPremium,
		LanguageCode:      userInitData.User.LanguageCode,
		ReferralCode:      utils.GenReferralCode(uid),
		RewardPoints:      200,
		LastPointsRefresh: time.Now().UTC(),
	})
	if e4 != nil {
		c.JSON(http.StatusInternalServerError, entities.ResFailed(entities.ErrInternalDBInsertFailed, e4.Error()))
		return
	}

	// refresh rewardPoints when user login
	e5 := app.DB().RefreshUserRewardPoints(uid)
	if e5 != nil {
		log.Printf("RefreshUserRewardPoints error %v\n", e5)
	}

	// clear user cache if login success
	app.Cache().Delete(caches.GenUserCacheKey(uid))

	basicConfig := app.Config().Basic
	session, e6 := utils.GenSession(uid, basicConfig.SessionExpiresSec, basicConfig.SessionEncryptKey)
	if e6 != nil {
		c.JSON(http.StatusInternalServerError, entities.ResFailed(entities.ErrGenUserSessionFailed, e5.Error()))
		return
	}

	if e7 := app.Cache().SetString(caches.GenUserSessionCacheKey(uid), session, basicConfig.SessionExpiresSec); e7 != nil {
		c.JSON(http.StatusInternalServerError, entities.ResFailed(entities.ErrInternalDBInsertFailed, e6.Error()))
	} else {
		c.JSON(http.StatusOK, entities.ResSuccess(&UserLoginRes{
			User:    user,
			IsNew:   isNew,
			Session: session,
			Secret:  utils.GenUserSecret(user.Id),
			Checkin: checkin,
		}))
	}
}

// UserFindOrCreate Try to find user by id, if not found, create user with newUser data, update username if changed
// return user, isNewUser, checkin, error
func userLoginAndCheckin(dbService *dbs.Service, uid int64, referral string, randPoint int64, newUser *dbs.User) (user *dbs.User, isNew bool, checkin *dbs.Checkin, err error) {
	err = dbService.DBInstance.Transaction(func(tx *gorm.DB) error {
		if referral != "" {
			var referralUser *dbs.User
			if e := tx.Where("users.referral_code = ?", referral).First(&referralUser).Error; e == nil && referralUser != nil {
				if referralUser.Id != uid { // cannot referral self
					newUser.ReferralUid = referralUser.Id
				}
			}
		}

		if result := tx.Where("users.id = ?", uid).Attrs(newUser).FirstOrCreate(&user); result.Error != nil {
			return result.Error
		} else {
			isNew = result.RowsAffected > 0
		}

		// new user should add a random point
		if isNew {
			if _, e0 := dbService.PointClaimTask(tx, uid, randPoint); e0 != nil {
				return e0
			}
		}

		// get or create a checkin res, if already claimed in recent duration, just return nil
		// otherwise return latest unclaimed checkin or create new checkin
		checkinRes, e1 := dbService.CheckinGetLatestCheckin(tx, uid, app.Config().Basic)
		if e1 != nil {
			return e1
		}
		checkin = checkinRes

		var changed = false
		if user.Username != newUser.Username {
			user.Username = newUser.Username
			changed = true
		}
		if user.LanguageCode != newUser.LanguageCode {
			user.LanguageCode = newUser.LanguageCode
			changed = true
		}
		if user.IsPremium != newUser.IsPremium {
			user.IsPremium = newUser.IsPremium
			changed = true
		}
		if changed {
			return tx.Save(user).Error
		} else {
			return nil
		}
	})
	return
}

// CheckinClaim
// @Tags User
// @Router /user/claim [post]
// @Summary Current user request claim daily checkin
// @description Current user request claim daily checkin
func CheckinClaim(c *gin.Context) {
	user, params := middleware.CheckUserAndJsonParams[entities.UserCheckinClaimParam](c)
	if user == nil || params == nil {
		return
	}

	var point *dbs.Point
	e0 := app.DB().DBInstance.Transaction(func(tx *gorm.DB) error {
		var checkin dbs.Checkin
		// find unclaimed checkin for specified user
		if e1 := tx.Where("id = ? AND uid = ? AND status = ?", params.CheckinId, user.Id, configs.CheckinStatusUnclaimed).First(&checkin).Error; e1 != nil {
			return e1
		}
		if _point, e2 := app.DB().PointClaimTask(tx, checkin.Uid, checkin.RewardPoint); e2 != nil {
			return e2
		} else {
			point = _point
		}
		checkin.Status = configs.CheckinStatusClaimed
		return tx.Save(&checkin).Error
	})
	if e0 != nil {
		c.JSON(http.StatusInternalServerError, entities.ResFailed(entities.ErrInternalDBUpdateFailed, e0.Error()))
	} else {
		c.JSON(http.StatusOK, entities.ResSuccess(*point))
	}
}

// GetUserPoint
// @Tags User
// @Router /user/point [get]
// @Summary Get current user's point
// @description Get current user's point
func GetUserPoint(c *gin.Context) {
	user := middleware.CurrentRequestUser(c)
	if user == nil {
		c.JSON(http.StatusUnauthorized, entities.ResFailed(entities.ErrUserNotFound, "unauthorized"))
		return
	}
	point, e0 := app.DB().PointFindByUid(user.Id)
	if e0 != nil {
		c.JSON(http.StatusInternalServerError, entities.ResFailed(entities.ErrInternalDBQueryFailed, e0.Error()))
	} else {
		c.JSON(http.StatusOK, entities.ResSuccess(point))
	}
}

// GetUserInvitedUserList
// @Tags User
// @Router /user/invited [get]
// @Summary Get current user's invited user list
// @description Get current user's invited user list
func GetUserInvitedUserList(c *gin.Context) {
	user, params := middleware.CheckUserAndQueryParams[entities.InvitedUserListParam](c)
	if user == nil || params == nil {
		return
	}
	users, total, e0 := app.DB().UserFindInvitedUserList(user.Id, params)
	if e0 != nil {
		c.JSON(http.StatusBadRequest, entities.ResFailed(entities.ErrInternalDBQueryFailed, e0.Error()))
	} else {
		c.JSON(http.StatusOK, entities.ResSuccess(&UserInvitedUserListRes{
			Users: users,
			Total: total,
		}))
	}
}

// GetLeaderboard
// @Tags User
// @Router /user/leaderboard [get]
// @Summary Get a list of user that sort by point
// @description Get a list of user that sort by point
func GetLeaderboard(c *gin.Context) {
	user, params := middleware.CheckUserAndQueryParams[entities.LeaderBoardParam](c)
	if user == nil || params == nil {
		return
	}

	users, total, e := app.DB().PointGetLeaderBoardWithUser(params)
	if e != nil {
		c.JSON(http.StatusBadRequest, entities.ResFailed(entities.ErrInternalDBQueryFailed, e.Error()))
	} else {
		c.JSON(http.StatusOK, entities.ResSuccess(&UserLeaderBoardRes{
			Users: users,
			Total: total,
		}))
	}
}
