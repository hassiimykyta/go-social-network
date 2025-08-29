package models

import "time"

type Post struct {
	Id          int64      `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	UserId      int64      `json:"user_id"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
}

type PostPublic struct {
	Id          int64     `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	UserId      int64     `json:"user_id"`
	CreatedAt   time.Time `json:"created_at"`
}

func (p Post) Public() PostPublic {
	return PostPublic{
		Id:          p.Id,
		Title:       p.Title,
		Description: p.Description,
		UserId:      p.UserId,
		CreatedAt:   p.CreatedAt,
	}
}
