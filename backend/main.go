package main

import (
	"inventory-vegetables-app/config"
	"inventory-vegetables-app/routes"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

	log.Println("アプリケーション開始。")
	// .envファイルの読み込み
	log.Println(".envファイルの読み込み中。。。")
	if err := godotenv.Load(); err != nil {
		log.Println(".envファイルが見つかりません。")
	} else {
		log.Print(".envファイル読み込み完了。")
	}

	// 環境変数の確認
	log.Printf("SUPABASE_URL: %s", os.Getenv("SUPABASE_URL"))
	log.Println("PORT:%s", os.Getenv("PORT"))

	// supabaseの初期化。
	log.Println("Supabase初期化中。。。")
	config.InitSupabase()

	// Ginルータの初期化
	log.Println("Ginルータの初期化中。。。")
	router := gin.Default()

	// CORSミドルウェア（後々必要になってくる）
	log.Println("CORS設定中。。。")
	router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
		c.Header("Access-Control-Allow-Headers", "COntent-Type,Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// テスト用ルートの設定
	log.Println("ルート設定中。。。")
	routes.SetupTestRoutes(router)

	// ポート設定
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("🚀 サーバーをポート %s で起動中...", port)
	log.Printf("📍 テストURL: http://localhost:%s/health", port)
	log.Printf("📍 Supabaseテスト: http://localhost:%s/test/supabase", port)
	log.Println("🔄 サーバー起動中... (Ctrl+C で停止)")

	// サーバーの起動
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("サーバー起動に失敗。%v", err)
	}
}
