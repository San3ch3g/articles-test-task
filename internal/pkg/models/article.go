package models

type Article struct {
	Id       uint32 `gorm:"primaryKey"`
	Title    string
	Content  string
	AuthorId uint32
	Author   Author `gorm:"foreignKey:AuthorId;references:Id"`
}
