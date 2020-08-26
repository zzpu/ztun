package model

import "time"

type Tunnel struct {
	ID        int64     `gorm:"column:id" json:"id"`
	Name      string    `gorm:"column:name" json:"name"`
	Port      string    `gorm:"port:hash" json:"port"`
	Remote    string    `gorm:"column:remote" json:"remote"`
	Status    int       `gorm:"column:status" json:"status"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
}
