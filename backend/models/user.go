package models

type User struct {
	ID           string `json:"id" db:"id"`
	Name         string `json:"name" db:"name" validate:"required,min=1,max=100"`
	Email        string `json:"email" db:"email" validate:"required,email"`
	PasswordHash string `json:"-" db: "password_hash"` //json出力はしない。
}
