package routers

import (
	"fmt"
	"game-mining-server/configs"
	"game-mining-server/routers/api"
	"game-mining-server/routers/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func CorsConfig(env string) gin.HandlerFunc {
	corsConf := cors.Config{MaxAge: 12 * time.Hour}

	if env == configs.EnvDEV {
		corsConf.AllowMethods = []string{"GET", "POST", "DELETE", "OPTIONS", "PUT"}
		corsConf.AllowHeaders = []string{"Authorization", "Content-Type", "Upgrade", "Origin", "Connection", "Accept-Encoding", "Accept-Language", "Host"}
	} else {
		corsConf.AllowMethods = []string{"GET", "POST", "DELETE", "OPTIONS", "PUT"}
		corsConf.AllowHeaders = []string{"Authorization", "Content-Type", "Origin", "Connection", "Accept-Encoding", "Accept-Language", "Host"}
	}

	corsConf.AllowAllOrigins = true

	return cors.New(corsConf)
}

func InitAndRun(config *configs.BasicConfig) error {
	env := config.Env

	gin.SetMode(getGinMode(env))

	r := gin.New()
	_ = r.SetTrustedProxies(nil)

	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(CorsConfig(env))

	bindProxyApi(r)
	bindNextApi(r)
	bindWalletApi(r)

	bindUserApi(r, config.Version)
	bindTaskApi(r, config.Version)
	bindMomentApi(r, config.Version)

	return r.Run(fmt.Sprintf(":%d", config.Port))
}

func getGinMode(env string) string {
	switch env {
	case configs.EnvDEV:
		return gin.DebugMode
	case configs.EnvTEST:
		return gin.TestMode
	default:
		return gin.ReleaseMode
	}
}

func bindProxyApi(r *gin.Engine) {
	group := r.Group("/proxy")
	group.GET("/html", middleware.LimitIp480PerMinMiddleware(), api.ProxyGetHtml)
	group.Match([]string{http.MethodPost, http.MethodGet, http.MethodOptions}, "/req", middleware.LimitIp480PerMinMiddleware(), api.ProxyRequest)
}

func bindNextApi(r *gin.Engine) {
	group := r.Group("/_next")
	group.GET("/*any", middleware.LimitIp480PerMinMiddleware(), api.ProxyGetNextRes)

}

func bindWalletApi(r *gin.Engine) {
	group := r.Group("/wallet")
	group.GET("/price", middleware.LimitIp240PerMinMiddleware(), api.GetCoinPrice)
}

func bindUserApi(r *gin.Engine, version int) {
	group := r.Group(fmt.Sprintf("/api/%d/user", version))
	group.POST("/login", middleware.LimitIp30PerMinMiddleware(), api.Login)
	group.POST("/claim", middleware.LimitIp60PerMinMiddleware(), middleware.AuthMiddleware(false), api.CheckinClaim)
	group.GET("/invited", middleware.LimitIp120PerMinMiddleware(), middleware.AuthMiddleware(false), api.GetUserInvitedUserList)
	group.GET("/point", middleware.LimitIp120PerMinMiddleware(), middleware.AuthMiddleware(false), api.GetUserPoint)
	group.GET("/leaderboard", middleware.LimitIp120PerMinMiddleware(), middleware.AuthMiddleware(false), api.GetLeaderboard)
}

func bindTaskApi(r *gin.Engine, version int) {
	group := r.Group(fmt.Sprintf("/api/%d/task", version))
	group.GET("/status", middleware.LimitIp120PerMinMiddleware(), middleware.AuthMiddleware(false), api.GetUserTaskStatus)
	group.POST("/claim", middleware.LimitIp60PerMinMiddleware(), middleware.AuthMiddleware(false), api.UserTaskClaim)
}

func bindMomentApi(r *gin.Engine, version int) {
	group := r.Group(fmt.Sprintf("/api/%d/moment", version))
	group.POST("/create", middleware.AuthMiddleware(false), api.CreateMoment)
	group.DELETE("/delete", middleware.AuthMiddleware(false), api.DeleteMoment)
	group.GET("/list", middleware.AuthMiddleware(false), api.GetLatestMoments)
	group.POST("/:id/comment", middleware.AuthMiddleware(false), api.AddComment)
	group.DELETE("/comment/:id", middleware.AuthMiddleware(false), api.DeleteComment)
	group.GET("/:id/comments", middleware.AuthMiddleware(false), api.GetCommentsForMoment)
	group.POST("/:id/like", middleware.AuthMiddleware(false), api.LikeMoment)
	group.DELETE("/:id/like", middleware.AuthMiddleware(false), api.RollbackLikeMoment)
	group.POST("/:id/reward", middleware.AuthMiddleware(false), api.RewardMoment)
}
