package routes

import (
	"inventory-vegetables-app/controllers"

	"github.com/gin-gonic/gin"
)

func SetupAuthRoutes(router *gin.Engine) {
	auth := router.Group("/auth")
	{
		// ユーザー登録
		auth.POST("/register", controllers.RegisterUser)

		// ユーザーログイン
		auth.POST("/login", controllers.LoginUser)

		// 店舗一覧取得
		auth.GET("/stores", controllers.GetStores)
	}
}
