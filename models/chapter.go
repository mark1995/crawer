package models

import "time"

type Chapter struct {
	Id        int       `gorm:"primary_key;AUTO_INCREMENT;column:id"`
	ChapterId int       `gorm:"column:name;not null" binding:"required"`
	Title     string    `gorm:"column:name;not null" binding:"required"`
	Content   string    `gorm:"column:content" binding:"required"`
	Sort      int       `gorm:"column:sort"`
	Pre       int       `gorm:"column:pre"`
	Next      int       `gorm:"column:next"`
	CreatedAt time.Time `gorm:"column:createAt"`
	UpdatedAt time.Time `gorm:"column:updateAt"`
}

func ChapterAdd(Chapter *Chapter) error {
	return DB.Self.Create(Chapter).Error
}

func ChapterUpdate(Chapter *Chapter) error {
	return DB.Self.Save(Chapter).Error
}

func ChapterDeleteByID(id int) error {
	Chapter := &Chapter{Id: id}
	return ChapterDelete(Chapter)
}

func ChapterDelete(Chapter *Chapter) error {
	return DB.Self.Delete(Chapter).Error
}

func GetChapterByName(name string) (*Chapter, error) {
	Chapter := &Chapter{}
	b := DB.Self.Where("name = ?", name).First(&Chapter)
	return Chapter, b.Error
}

func GetChapterByID(ID int) (*Chapter, error) {
	Chapter := &Chapter{}
	b := DB.Self.Where("id = ?", ID).First(&Chapter)
	return Chapter, b.Error
}
