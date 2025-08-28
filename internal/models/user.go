package models

import "time"

type User struct {
	Id           int64      `json:"id"`
	Username     string     `json:"username"`
	Email        string     `json:"email"`
	PasswordHash string     `json:"-"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	DeletedAt    *time.Time `json:"deleted_at,omitempty"`
}

type UserPublic struct {
	Id        int64     `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

func (u User) Public() UserPublic {
	return UserPublic{
		Id:        u.Id,
		Username:  u.Username,
		Email:     u.Email,
		CreatedAt: u.CreatedAt,
	}
}
