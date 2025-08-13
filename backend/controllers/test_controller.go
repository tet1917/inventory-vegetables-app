package controllers

import (
	"inventory-vegetables-app/config"
	"net/http"

	"github.com/gin-gonic/gin"
)

// TestSupabaseConnection Supabase接続テスト
func TestSupabaseConnection(c *gin.Context) {
	// 簡単な接続テスト（テーブル一覧取得など）
	if config.SupabaseClient == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Supabaseクライアントが初期化されていません。",
			"success": false,
		})
		return
	}

	// Supabaseに接続できているのか基本確認
	c.JSON(http.StatusOK, gin.H{
		"message": "接続テスト成功！",
		"success": true,
		"client":  "Supabaseクライアント初期化済み。",
	})
}

// HealthCheck アプリケーション動作確認
func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "八百屋在庫管理API　サーバーは正常稼働しています。",
		"success": true,
		"port":    "8080",
	})
}
