package model

import (
	"github.com/go-programming-tour-book/blog-service/pkg/app"
	"github.com/jinzhu/gorm"
)

type Tag struct {
	*Model
	Name  string `json:"name"`
	State uint8  `json:"state"`
}

type TagSwagger struct {
	List []*Tag
	Pager *app.Pager
}

func (t Tag) TableName() string {
	return "blog_tag"
}

func (t Tag) Create(db *gorm.DB) error {
	return db.Create(&t).Error // todo 这里为啥需要指针
}

func (t Tag) Update(db *gorm.DB, values interface{}) error {
	//return db.Model(&Tag{}).Where("id = ? AND is_del = ?", t.ID, 0).Update(t).Error // todo 这里为啥需要新的Tag指针
	// update(t)会自动填充id = {}条件， update(values)不会，需要自己手动where
	return db.Model(&Tag{}).Where("id = ? AND is_del = ?", t.ID, 0).Update(values).Error
}

func (t Tag) Count(db *gorm.DB) (int, error) {
	var count int
	if t.Name != "" {
		db = db.Where("name = ?", t.Name)
	}
	db = db.Where("state = ?", t.State)
	err := db.Model(&t).Where("is_del = ?", 0).Count(&count).Error
	if err != nil {
		return 0, err
	}

	return count, nil
}

	func (t Tag) List(db *gorm.DB, pageOffset, pageSize int) ([]*Tag, error) {
	var tags []*Tag
	var err error
	if pageOffset >= 0 && pageSize > 0 {
		db = db.Offset(pageOffset).Limit(pageSize)
	}
	if t.Name != "" {
		db = db.Where("name = ?", t.Name)
	}
	db = db.Where("state = ?", t.State)
	err = db.Model(&t).Where("is_del = ?", 0).Find(&tags).Error
	if err != nil {
		return nil, err
	}

	return tags, nil
}

func (t Tag) Delete(db *gorm.DB) error {
	return db.Where("id = ? AND is_del = ?", t.ID, 0).Delete(&t).Error // todo 这里为啥需要指针
}

func (t Tag) Get(db *gorm.DB) (Tag, error) {
	var tag Tag
	db = db.Where("id = ?", t.ID)
	err := db.First(&tag).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return tag, err
	}

	return tag, nil
}