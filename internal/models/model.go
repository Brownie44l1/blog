package models

import "time"

type User struct {
	ID       int64  `db:"id" json:"id"`
	Username string `db:"username" json:"username"`
	Password string `db:"password" json:"-"`
	BlogCount int    `db:"blog_count" json:"blog_count"`
}

type Blog struct {
	ID        int64     `db:"id" json:"id"`
	UserId    int64     `db:"user_id" json:"user_id"`
	Title     string    `db:"title" json:"title"`
	Content   string    `db:"content" json:"content"`
	//ImageURL  *string   `db:"image_url" json:"image_url,omitempty"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}
