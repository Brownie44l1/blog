package models

type Image struct {
	ID     int64  `db:"id" json:"id"`
	BlogID int64  `db:"blog_id" json:"blog_id"`
	URL    string `db:"url" json:"url"`
}
