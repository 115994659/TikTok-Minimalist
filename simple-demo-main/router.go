package main

import (
	"github.com/RaymondCode/simple-demo/controller"
	"github.com/RaymondCode/simple-demo/middleware"
	"github.com/RaymondCode/simple-demo/model"
	"github.com/gin-gonic/gin"
)

func initRouter(r *gin.Engine) {
	model.InitDB()
	// 用static作为静态变量存放目录
	r.Static("/static", "./static")

	apiRouter := r.Group("/douyin")

	// basic apis
	apiRouter.GET("/feed/", controller.FeedVideoListController)
	apiRouter.GET("/user/", middleware.JWTAuthTo(), controller.UserInfoController)
	apiRouter.POST("/user/register/", middleware.SHAPassword(), controller.RegisterController)
	apiRouter.POST("/user/login/", middleware.SHAPassword(), controller.LoginController)
	apiRouter.POST("/publish/action/", middleware.JWTAuthTo(), controller.PublishController)
	apiRouter.GET("/publish/list/", middleware.JWTAuthTo(), controller.PublishListController)

	// extra apis - I
	apiRouter.POST("/favorite/action/", middleware.JWTAuthTo(), controller.FavoriteActionController)
	apiRouter.GET("/favorite/list/", middleware.JWTAuthTo(), controller.FavoriteListController)
	apiRouter.POST("/comment/action/", middleware.JWTAuthTo(), controller.CommentActionController)
	apiRouter.GET("/comment/list/", middleware.JWTAuthTo(), controller.CommentListController)

	// extra apis - II
	//开发中
	apiRouter.POST("/relation/action/", controller.RelationAction)
	apiRouter.GET("/relation/follow/list/", controller.FollowList)
	apiRouter.GET("/relation/follower/list/", controller.FollowerList)
	apiRouter.GET("/relation/friend/list/", controller.FriendList)
	apiRouter.GET("/message/chat/", controller.MessageChat)
	apiRouter.POST("/message/action/", controller.MessageAction)
}
