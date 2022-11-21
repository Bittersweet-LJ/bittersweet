package routes

import (
	"bittersweet/controller"
	"bittersweet/logger"
	"bittersweet/middlewares"
	"net/http"
	"time"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "bittersweet/docs" // 千万不要忘了导入把你上一步生成的docs

	"github.com/gin-gonic/gin"
)

func Setup(mode string) *gin.Engine {

	if mode == gin.ReleaseMode {
		gin.SetMode(gin.ReleaseMode) //gin 设置为发布模式
	}
	r := gin.New()
	r.Use(logger.GinLogger(), logger.GinRecovery(true), middlewares.RateLimitMiddleware(1*time.Second, 1))

	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	// 注册swagger api相关路由
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	v1 := r.Group("/api/v1")

	//注册业务路由
	v1.POST("/signup", controller.SignUpHandler)
	//登录
	v1.POST("/login", controller.LoginHandler)

	v1.GET("/community", controller.CommunityHandler)
	v1.GET("/community/:id", controller.CommunityDetailHandler)
	v1.GET("/post/:id", controller.GetPostDetailHandler)
	v1.GET("/posts", controller.GetPostListHandler)
	//根据时间或分数获取帖子列表
	v1.GET("/posts2", controller.GetPostListHandler2)

	v1.Use(middlewares.JWTAuthMiddleware()) // 应用JWT认证中间件
	{
		v1.POST("/post", controller.CreatePostHandler)
		//帖子投票
		v1.POST("/vote", controller.PostVoteHandler)
	}

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"msg": "404",
		})
	})
	return r
}
