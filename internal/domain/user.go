package domain

type User struct {
	ID           int64  `json:"id" db:"id" gorm:"primaryKey;autoIncrement:false"`
	Name         string `json:"name" db:"name"`
	Email        string `json:"email" db:"email"`
	Password     string `json:"-" db:"password"`
	PasswordSalt string `json:"-" db:"password_salt"`
}
