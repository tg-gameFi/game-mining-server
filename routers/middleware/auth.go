package middleware

import (
	"game-mining-server/app"
	"game-mining-server/caches"
	"game-mining-server/configs"
	"game-mining-server/dbs"
	"game-mining-server/entities"
	"game-mining-server/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"time"
)

func parseAuthUser(authHeader string) (*dbs.User, int, string) {
	split := strings.Split(authHeader, " ")
	// token format: Bearer eyJhbGciOiJIUzI1NiIsInR..., split and take at index 1
	if authHeader == "" || len(split) != 2 || split[0] != "Bearer" {
		return nil, entities.ErrInvalidAuthHeader, "bad auth header format"
	}
	session, e0 := utils.ParseSession(split[1], app.Config().Basic.SessionEncryptKey)
	if e0 != nil {
		return nil, entities.ErrInvalidAuthHeader, e0.Error()
	}

	if time.Now().Unix() >= session.IssuedAt+int64(session.ExpiresSec*1000) {
		return nil, entities.ErrUserAuthExpired, "session expired"
	}

	user, err := caches.UserFindByIdCached(app.Cache(), app.DB(), session.Uid, app.Config().Basic.SessionExpiresSec)
	if err != nil {
		return nil, entities.ErrUserNotFound, "user not found"
	}

	return user, entities.Ok, ""
}

// AuthMiddleware allowPublic true to allow request with no auth header, but auth header is provided, valid it
func AuthMiddleware(allowPublic bool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		// allow public and auth header is empty, go through as a public request
		if allowPublic && (authHeader == "" || len(authHeader) == 0) {
			ctx.Set(configs.CurUser, nil)
			ctx.Next()
		} else if user, code, msg := parseAuthUser(authHeader); code != entities.Ok || user == nil {
			ctx.JSON(http.StatusUnauthorized, entities.ResFailed(code, msg))
			ctx.Abort()
		} else {
			ctx.Set(configs.CurUser, *user)
			ctx.Next()
		}
	}
}

func CurrentRequestUser(c *gin.Context) *dbs.User {
	userData, exist := c.Get(configs.CurUser)
	if userData == nil || !exist {
		return nil
	}
	user := userData.(dbs.User)
	return &user
}

func CheckUserAndJsonParams[T any](c *gin.Context) (*dbs.User, *T) {
	user := CurrentRequestUser(c)
	if user == nil {
		c.JSON(http.StatusUnauthorized, entities.ResFailed(entities.ErrUserNotFound, "unauthorized"))
		return nil, nil
	}
	var params T
	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, entities.ResFailed(entities.ErrInvalidParams, err.Error()))
		return nil, nil
	}
	return user, &params
}

func CheckUserAndQueryParams[T any](c *gin.Context) (*dbs.User, *T) {
	user := CurrentRequestUser(c)
	if user == nil {
		c.JSON(http.StatusUnauthorized, entities.ResFailed(entities.ErrUserNotFound, "unauthorized"))
		return nil, nil
	}
	var params T
	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, entities.ResFailed(entities.ErrInvalidParams, err.Error()))
		return nil, nil
	}
	return user, &params
}
