package model

import "github.com/jinzhu/gorm"

type ArticleTag struct {
	*Model
	ArticleId uint32 `json:"article_id"`
	TagId     uint32 `json:"tag_id"`
}

func (a ArticleTag) TableName() string {
	return "blog_article_tag"
}

func (a ArticleTag) Create(db *gorm.DB) error {
	return db.Create(&a).Error
}

func (a ArticleTag) Get(db *gorm.DB) (ArticleTag, error) {
	var articleTag ArticleTag
	db = db.Where("article_id = ?", a.ArticleId)
	err := db.First(&articleTag).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return articleTag, err
	}

	return articleTag, nil
}
