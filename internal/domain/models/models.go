package models

import "time"

type Task struct {
	Id          int64     `json:"id"  db:"id"`
	Title       string    `json:"title" db:"title"`
	Description string    `json:"description" db:"description"`
	Status      string    `json:"status" db:"status"`
	Created_at  time.Time `json:"created_at" db:"created_at"`
	Update_at   time.Time `json:"updated_at" db:"updated_at"`
}
