package routes

import (
	"inventory-vegetables-app/controllers"

	"github.com/gin-gonic/gin"
)

// SetupTestRoutes テスト用のルートを設定
func SetupTestRoutes(router *gin.Engine) {
	// ヘルスチェック
	router.GET("/health", controllers.HealthCheck)

	// Supabase接続テスト
	router.GET("/test/supabase", controllers.TestSupabaseConnection)

}
