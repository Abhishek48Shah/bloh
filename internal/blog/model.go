package blog

import "time"

type Blog struct {
	Id       int       `json:"id"`
	Title    string    `json:"title"`
	Content  string    `json:"content"`
	CreateAt time.Time `json:"create_at"`
	UpdateAt time.Time `json:"update_at"`
}
type CreateBlogInput struct {
	Title   string `json:"title" validate:"required,min=20,max=200"`
	Content string `json:"content" validate:"required"`
}
