package models

import "time"

type User struct {
	Id          string    `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	DisplayName string    `json:"display_name"`
	Email       string    `json:"email"`
}
