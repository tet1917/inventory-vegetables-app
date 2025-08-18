package controllers

import (
	"encoding/json"
	"inventory-vegetables-app/config"
	"inventory-vegetables-app/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// ユーザー登録
func RegisterUser(c *gin.Context) {
	var req models.UserRegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "リクエスト形式が正しくありません。",
			"details": err.Error(),
			"success": false,
		})
		return
	}

	// バリデーション
	if err := req.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "入力値が正しくありません。",
			"details": err.Error(),
			"success": false,
		})
		return
	}

	// メール重複チェック
	result, _, err := config.SupabaseClient.From("users").
		Select("id", "", false).
		Eq("email", req.Email).
		Execute()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "データベースエラー",
			"details": err.Error(),
			"success": false,
		})
		return
	}

	// JSONデータをパース
	var existingUsers []map[string]interface{}
	if err := json.Unmarshal(result, &existingUsers); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "データベースレスポンスの解析エラー",
			"details": err.Error(),
			"success": false,
		})
		return
	}

	// 既に存在するユーザーがいる場合
	if len(existingUsers) > 0 {
		c.JSON(http.StatusConflict, gin.H{
			"error":   "このメールアドレスはすでに登録されています。",
			"success": false,
		})
		return
	}

	// store_idが指定されている場合、店舗存在チェック
	if req.StoreID != nil && *req.StoreID != "" {
		storeResult, _, err := config.SupabaseClient.From("stores").
			Select("id", "", false).
			Eq("id", *req.StoreID).
			Execute()

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "店舗存在チェックでエラーが発生しました。",
				"details": err.Error(),
				"success": false,
			})
			return
		}

		var stores []map[string]interface{}
		if err := json.Unmarshal(storeResult, &stores); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "店舗データの解析エラー",
				"details": err.Error(),
				"success": false,
			})
			return
		}

		if len(stores) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "指定された店舗が存在しません。",
				"success": false,
			})
			return
		}
	}

	// パスワードのハッシュ化
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "パスワード処理でエラーが発生しました。",
			"details": err.Error(),
			"success": false,
		})
		return
	}

	// ユーザー情報をDBに挿入
	userData := map[string]interface{}{
		"name":          req.Name,
		"email":         req.Email,
		"password_hash": string(hashedPassword),
		"role":          req.Role,
		"store_id":      req.StoreID,
	}

	insertResult, _, err := config.SupabaseClient.From("users").
		Insert(userData, false, "", "", "").
		Execute()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "ユーザー登録に失敗しました",
			"details": err.Error(),
			"success": false,
		})
		return
	}

	// 挿入結果をパース
	var insertedUser []map[string]interface{}
	if err := json.Unmarshal(insertResult, &insertedUser); err != nil {
		// パースエラーでも登録は成功しているので、エラーではなく警告として扱う
		c.JSON(http.StatusCreated, gin.H{
			"message": "✅ ユーザー登録に成功しました。",
			"success": true,
			"note":    "レスポンスデータの解析でエラーが発生しましたが、登録は完了しました。",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "✅ ユーザー登録に成功しました。",
		"success": true,
		"data":    insertedUser,
	})
}

// ユーザーログイン
func LoginUser(c *gin.Context) {
	var req models.UserLoginRequest
	// JSONリクエストのパース
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "リクエスト形式が正しくありません。",
			"details": err.Error(),
			"success": false,
		})
		return
	}
	// バリデーション
	if err := req.ValidateLogin(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "ログインIDとパスワードは必須です。",
			"details": err.Error(),
			"success": false,
		})
		return
	}

	// ユーザー検索（メールアドレス）
	result, _, err := config.SupabaseClient.From("users").Select("id,name,email,password_hash,role,store_id,created_at", "", false).Eq("email", req.LoginID).Execute()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "データベースエラー",
			"details": err.Error(),
			"success": false,
		})
		return
	}
	var users []map[string]interface{}
	if err := json.Unmarshal(result, &users); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "データベースレスポンスの解析エラー",
			"details": err.Error(),
			"success": false,
		})
		return
	}
	// ユーザーが見つからない場合、ユーザー名で検索
	if len(users) == 0 {
		result, _, err = config.SupabaseClient.From("users").Select("id,name,email,password_hash,role,store_id,created_at", "", false).Eq("name", req.LoginID).Execute()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "データベースエラー",
				"details": err.Error(),
				"success": false,
			})
			return
		}
		if err := json.Unmarshal(result, &users); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "データベースレスポンスの解析エラー",
				"details": err.Error(),
				"success": false,
			})
			return
		}
	}
	// 最終的にユーザーが見つからない場合。
	if len(users) == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "ログインIDまたはパスワードが正しくありません。",
			"success": false,
		})
		return
	}

	user := users[0]
	// パスワード確認
	storedPassword, ok := user["password_hash"].(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "パスワードのデータ取得でエラーが発生しました。",
			"success": false,
		})
		return
	}
	// パスワード照合
	err = bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(req.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "ログインIDまたはパスワードが正しくありません。",
			"success": false,
		})
		return
	}
	// 店舗情報を取得（store_idが存在する場合）
	var storeName *string
	if storeID := user["store_id"]; storeID != nil {
		storeResult, _, err := config.SupabaseClient.From("stores").
			Select("name", "", false).
			Eq("id", storeID.(string)).
			Execute()

		if err == nil {
			var stores []map[string]interface{}
			if err := json.Unmarshal(storeResult, &stores); err == nil && len(stores) > 0 {
				if name, ok := stores[0]["name"].(string); ok {
					storeName = &name
				}
			}
		}
	}
	// ログイン成功レスポンス
	response := models.UserLoginResponse{
		ID:        user["id"].(string),
		Name:      user["name"].(string),
		Email:     user["email"].(string),
		Role:      user["role"].(string),
		StoreName: storeName,
		Token:     "temp-token-" + user["id"].(string), // 仮のトークン（後でJWT実装予定）
	}

	// store_idの処理（nilの可能性がある。）
	if storeID := user["store_id"]; storeID != nil {
		storeIDStr := storeID.(string)
		response.StoreID = &storeIDStr
	}
	// created_atの処理
	if createdAt, ok := user["created_at"].(string); ok {
		if parsedTime, err := time.Parse(time.RFC3339, createdAt); err == nil {
			response.CreatedAt = parsedTime
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "ログイン成功",
		"success": true,
		"user":    response,
	})
}

// GetStores 店舗一覧取得（登録時の選択用）
func GetStores(c *gin.Context) {
	result, _, err := config.SupabaseClient.From("stores").
		Select("id, name, location", "", false).
		Execute()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "店舗情報の取得に失敗しました。",
			"details": err.Error(),
			"success": false,
		})
		return
	}

	// JSONデータをパース
	var stores []map[string]interface{}
	if err := json.Unmarshal(result, &stores); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "店舗データの解析エラー",
			"details": err.Error(),
			"success": false,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "✅ 店舗一覧取得成功",
		"success": true,
		"data":    stores,
	})
}
