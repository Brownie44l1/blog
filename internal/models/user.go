package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `db:"id" json:"id"`
	Username  string    `db:"Username" json:"username"`
	Email     string    `db:"email" json:"email"`
	Password  string    `db:"password" json:"-"`
	BlogCount int64     `db:"blog_count" json:"blog_count"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}
