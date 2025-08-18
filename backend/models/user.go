package models

import (
	"time"

	"github.com/go-playground/validator/v10"
)

type User struct {
	ID           string    `json:"id" db:"id"`
	Name         string    `json:"name" db:"name" validate:"required,min=1,max=100"`
	Email        string    `json:"email" db:"email" validate:"required,email"`
	PasswordHash string    `json:"-" db: "password_hash"` //json出力はしない。
	Role         string    `json:"role" validate:"required,oneof=admin member"`
	StoreID      *string   `json:"store_id" db:"store_id"` //ポインタ型
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

// Store 店舗モデル
type Store struct {
	ID        string    `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	Location  string    `json:"location" db:"location"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// UserRegisterRequest ユーザー登録リクエスト
type UserRegisterRequest struct {
	Name     string  `json:"name" validate:"required,min=1,max=100"`
	Email    string  `json:"email" validate:"required,email"`
	Password string  `json:"password" validate:"required,min=6"`
	Role     string  `json:"role" validate:"required,oneof=admin member"`
	StoreID  *string `json:"store_id"` // 管理者の場合はnull可
}

// UserResponse ユーザー情報レスポンス（パスワードなし）
type UserResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	StoreID   *string   `json:"store_id"`
	StoreName *string   `json:"store_name,omitempty"` // JOIN結果用
	CreatedAt time.Time `json:"created_at"`
}

// ユーザーログインリクエスト
type UserLoginRequest struct {
	LoginID  string `json:"login_id" validate:"required"` // ユーザー名 or メールアドレス
	Password string `json:"password" validate:"required"` // パスワード
}

// ユーザーログインレスポンス
type UserLoginResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	StoreID   *string   `json:"store_id"`
	StoreName *string   `json:"store_name,omitempty"`
	Token     string    `json:"token"` // 後でJWT実装予定
	CreatedAt time.Time `json:"created_at"`
}

// Validate バリデーション実行
func (req *UserRegisterRequest) Validate() error {
	validator := validator.New()
	return validator.Struct(req)
}

func (req *UserLoginRequest) ValidateLogin() error {
	validate := validator.New()
	return validate.Struct(req)
}
