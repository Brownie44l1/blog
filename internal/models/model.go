package models

import "time"

type User struct {
	ID       int64  `db:"id" json:"id"`
	Username string `db:"username" json:"username"`
	Password string `db:"password" json:"-"`
}

type Blog struct {
	ID        int64     `db:"id" json:"id"`
	UserId    int64     `db:"user_id" json:"user_id"`
	Title     string    `db:"title" json:"title"`
	Content   string    `db:"content" json:"content"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

/* {"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJpc3MiOiJibG9nLWFwaSIsInN1YiI6InVzZXJfYXV0aGVudGljYXRpb24iLCJleHAiOjE3NjQ2MDk5OTQsIm5iZiI6MTc2NDUyMzU5NCwiaWF0IjoxNzY0NTIzNTk0fQ.RKNgRSvtSBMkDzvzAXr7JN_LWyNB8hMQO9hbpzVm7sw","username":"john_doe"} */
