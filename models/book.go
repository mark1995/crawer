package models

import "time"

type Book struct {
	Id        int    `gorm:"primary_key;AUTO_INCREMENT;column:id"`
	Name      string `gorm:"column:name;not null" binding:"required"`
	Author    string `gorm:"column:author;not null" binding:"required"`
	Image     string
	Status    int
	From      string
	Url       string    `gorm:"column:url"`
	CreatedAt time.Time `gorm:"column:createAt"`
	UpdatedAt time.Time `gorm:"column:updateAt"`
}

func BookAdd(book *Book) error {
	return DB.Self.Create(book).Error
}

func BookUpdate(book *Book) error {
	return DB.Self.Save(book).Error
}

func BookDeleteByID(id int) error {
	book := &Book{Id: id}
	return BookDelete(book)
}

func BookDelete(book *Book) error {
	return DB.Self.Delete(book).Error
}

func GetBookByName(name string) (*Book, error) {
	book := &Book{}
	b := DB.Self.Where("name = ?", name).First(&book)
	return book, b.Error
}

func GetBookByID(ID int) (*Book, error) {
	book := &Book{}
	b := DB.Self.Where("id = ?", ID).First(&book)
	return book, b.Error
}
