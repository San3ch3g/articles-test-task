package models

type Author struct {
	Id       uint32 `gorm:"primaryKey;"`
	Username string `gorm:"unique"`
	Password string
}
