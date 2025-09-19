package models

import "time"

type Base struct {
	ID        uint64    `json:"id" gorm:"primaryKey,autoIncrement,index,unique,uniqueIndex"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (b Base) IsValid() bool {
	return b.ID > 0 && !b.CreatedAt.IsZero()
}
