package api

import (
	"game-mining-server/app"
	"game-mining-server/configs"
	"game-mining-server/dbs"
	"game-mining-server/entities"
	"game-mining-server/routers/middleware"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type UserTaskStatus struct {
	SocialTasks      []*dbs.Task `json:"socialTasks"`
	Point            *dbs.Point  `json:"point"`
	InvitedUserCount int64       `json:"invitedUserCount"`
}

// GetUserTaskStatus
// @Tags Task
// @Router /task/status [get]
// @Summary Get current user all tasks status and related info
// @description Get current user all tasks status and related info
func GetUserTaskStatus(c *gin.Context) {
	user := middleware.CurrentRequestUser(c)
	if user == nil {
		c.JSON(http.StatusUnauthorized, entities.ResFailed(entities.ErrUserNotFound, "unauthorized"))
		return
	}

	socialTasks, e0 := app.DB().TaskFindAllOrCreateSocial(user.Id)
	if e0 != nil {
		c.JSON(http.StatusBadRequest, entities.ResFailed(entities.ErrInternalDBQueryFailed, e0.Error()))
		return
	}

	point, e1 := app.DB().PointFindByUid(user.Id)
	if e1 != nil {
		c.JSON(http.StatusBadRequest, entities.ResFailed(entities.ErrInternalDBQueryFailed, e1.Error()))
		return
	}

	invitedCount, e2 := app.DB().UserCountInvitedUsers(user.Id)
	if e2 != nil {
		c.JSON(http.StatusBadRequest, entities.ResFailed(entities.ErrInternalDBQueryFailed, e2.Error()))
		return
	}

	c.JSON(http.StatusOK, entities.ResSuccess(&UserTaskStatus{
		SocialTasks:      socialTasks,
		Point:            point,
		InvitedUserCount: invitedCount,
	}))
}

// UserTaskClaim
// @Tags Task
// @Router /task/claim [post]
// @Summary Current user request to claim a task
// @description Current user request to claim a task
func UserTaskClaim(c *gin.Context) {
	user, params := middleware.CheckUserAndJsonParams[entities.UserTaskClaimParam](c)
	if user == nil || params == nil {
		return
	}

	var point *dbs.Point
	var err error
	if params.TaskGroup == configs.TaskGroupSocial {
		// social claim, key is taskId
		point, err = app.DB().TaskClaim(params.ClaimKey, user.Id, configs.TaskStatusClaimable, configs.TaskStatusClaimed)
	} else if params.TaskGroup == configs.TaskGroupWallet {
		// wallet claim, hash is not used for now
		point, err = app.DB().PointClaimForWallet(user.Id, configs.TaskWalletBaseRewardPoint)
	} else if params.TaskGroup == configs.TaskGroupInvite {
		// invite claim, key is level
		if level, e := strconv.ParseInt(params.ClaimKey, 10, 64); e != nil {
			c.JSON(http.StatusBadRequest, entities.ResFailed(entities.ErrInvalidParams, e.Error()))
		} else {
			point, err = app.DB().PointClaimForInvite(user.Id, level)
		}
	}
	if err != nil {
		c.JSON(http.StatusBadRequest, entities.ResFailed(entities.ErrInternalDBUpdateFailed, err.Error()))
	} else {
		c.JSON(http.StatusOK, entities.ResSuccess(point))
	}
}
