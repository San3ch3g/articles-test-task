package models

type Author struct {
	Id       uint32 `gorm:"primaryKey;"` //TODO ForeignKey для ID
	Username string `gorm:"unique"`
	Password string
}
