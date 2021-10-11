package domain

import "time"

type Product struct {
	ID        int64     `db:"id" json:"id" gorm:"primaryKey;autoIncrement:false" `
	Name      string    `db:"name" json:"name" `
	ImageURL  string    `db:"image_url" json:"image_url" `
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}
