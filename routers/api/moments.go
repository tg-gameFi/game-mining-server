package api

import (
	"fmt"
	"game-mining-server/app"
	"game-mining-server/entities"
	"github.com/gin-gonic/gin"
	"net/http"
)

func CreateMoment(c *gin.Context) {
	var params entities.CreateMomentParam
	if e0 := c.ShouldBindJSON(&params); e0 != nil {
		c.JSON(http.StatusBadRequest, entities.ResFailed(entities.ErrInvalidParams, e0.Error()))
		return
	}

	momentId, err := app.DB().CreateMoment(params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, entities.ResFailed(entities.ErrInternalDBInsertFailed, fmt.Sprintf("Failed to create moment, %s", err.Error())))
		return
	}

	c.JSON(http.StatusOK, entities.ResSuccess(gin.H{
		"message":   "Moment created successfully",
		"moment_id": momentId,
	}))
}

func DeleteMoment(c *gin.Context) {
	var params entities.DeleteMomentParam
	if e0 := c.ShouldBindJSON(&params); e0 != nil {
		c.JSON(http.StatusBadRequest, entities.ResFailed(entities.ErrInvalidParams, e0.Error()))
		return
	}

	err := app.DB().DeleteMoment(params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, entities.ResFailed(entities.ErrInternalDBInsertFailed, fmt.Sprintf("Failed to delete moment, %s", err.Error())))
		return
	}

	c.JSON(http.StatusOK, entities.ResSuccess("Moment deleted successfully"))
}

func GetLatestMoments(c *gin.Context) {
	var params entities.GetLatestMomentsParam
	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, entities.ResFailed(entities.ErrInvalidParams, err.Error()))
		return
	}

	moments, err := app.DB().GetLatestMoments(params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, entities.ResFailed(entities.ErrInternalDBInsertFailed, fmt.Sprintf("Failed to get moments, %s", err.Error())))
		return
	}

	c.JSON(http.StatusOK, entities.ResSuccess(moments))
}

func AddComment(c *gin.Context) {
	var params entities.AddCommentParam
	if e0 := c.ShouldBindJSON(&params); e0 != nil {
		c.JSON(http.StatusBadRequest, entities.ResFailed(entities.ErrInvalidParams, e0.Error()))
		return
	}

	commentId, err := app.DB().CommentMoment(params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, entities.ResFailed(entities.ErrInternalDBInsertFailed, fmt.Sprintf("Failed to add comment, %s", err.Error())))
		return
	}

	c.JSON(http.StatusOK, entities.ResSuccess(gin.H{
		"message":    "comment added successfully",
		"comment_id": commentId,
	}))
}

func DeleteComment(c *gin.Context) {
	var params entities.DeleteCommentParam
	if e0 := c.ShouldBindJSON(&params); e0 != nil {
		c.JSON(http.StatusBadRequest, entities.ResFailed(entities.ErrInvalidParams, e0.Error()))
		return
	}

	err := app.DB().DeleteComment(params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, entities.ResFailed(entities.ErrInternalDBInsertFailed, fmt.Sprintf("Failed to delete comment, %s", err.Error())))
		return
	}

	c.JSON(http.StatusOK, entities.ResSuccess("Comment deleted successfully"))
}

func GetCommentsForMoment(c *gin.Context) {
	var params entities.GetCommentsForMomentParam
	if e0 := c.ShouldBindQuery(&params); e0 != nil {
		c.JSON(http.StatusBadRequest, entities.ResFailed(entities.ErrInvalidParams, e0.Error()))
		return
	}

	comments, err := app.DB().GetCommentsForMoment(params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, entities.ResFailed(entities.ErrInternalDBInsertFailed, "Failed to get comments"))
		return
	}

	c.JSON(http.StatusOK, entities.ResSuccess(comments))
}

func LikeMoment(c *gin.Context) {
	var params entities.LikeMomentParam
	if e0 := c.ShouldBindJSON(&params); e0 != nil {
		c.JSON(http.StatusBadRequest, entities.ResFailed(entities.ErrInvalidParams, e0.Error()))
		return
	}

	err := app.DB().LikeMoment(params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, entities.ResFailed(entities.ErrInternalDBInsertFailed, fmt.Sprintf("Failed to like moment, %s", err.Error())))
		return
	}

	c.JSON(http.StatusOK, entities.ResSuccess("Moment liked successfully"))
}

func RollbackLikeMoment(c *gin.Context) {
	var params entities.LikeMomentParam
	if e0 := c.ShouldBindJSON(&params); e0 != nil {
		c.JSON(http.StatusBadRequest, entities.ResFailed(entities.ErrInvalidParams, e0.Error()))
		return
	}

	err := app.DB().RollbackLikeMoment(params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, entities.ResFailed(entities.ErrInternalDBInsertFailed, fmt.Sprintf("Failed to rollback like moment, %s", err.Error())))
		return
	}

	c.JSON(http.StatusOK, entities.ResSuccess("Moment unliked successfully"))
}

func RewardMoment(c *gin.Context) {
	var params entities.RewardMomentParam
	if e0 := c.ShouldBindJSON(&params); e0 != nil {
		c.JSON(http.StatusBadRequest, entities.ResFailed(entities.ErrInvalidParams, e0.Error()))
		return
	}

	err := app.DB().RewardMoment(params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, entities.ResFailed(entities.ErrInternalDBInsertFailed, fmt.Sprintf("Failed to reward moment, %s", err.Error())))
		return
	}

	c.JSON(http.StatusOK, entities.ResSuccess("Moment rewarded successfully"))
}
