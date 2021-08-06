package dao

import "github.com/go-programming-tour-book/blog-service/internal/model"

func (d *Dao) CreateArticleTag(articleId, tagId uint32, createdBy string) error {
	articleTag := model.ArticleTag{
		Model:     &model.Model{CreatedBy: createdBy},
		ArticleId: articleId,
		TagId:     tagId,
	}

	return articleTag.Create(d.engine)
}

func (d *Dao) GetArticleTagByAID(articleId uint32) (model.ArticleTag, error) {
	articleTag := model.ArticleTag{ArticleId: articleId}

	return articleTag.Get(d.engine)
}
