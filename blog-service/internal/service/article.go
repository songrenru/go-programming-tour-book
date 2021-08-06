package service

import (
	"github.com/go-programming-tour-book/blog-service/internal/dao"
	"github.com/go-programming-tour-book/blog-service/internal/model"
	"github.com/go-programming-tour-book/blog-service/pkg/app"
)

type ArticleRequest struct {
	ID uint32 `form:"id" binding:"required,gte=1"`
	State uint8 `form:"state,default=1" binding:"oneof=0 1"`
}

type ArticleListRequest struct {
	TagId uint32 `form:"tag_id" binding:"gte=1"`
	State uint8 `form:"state,default=1" binding:"oneof=0 1"`
}

type CreateArticleRequest struct {
	TagId uint32 `form:"tag_id" binding:"required,gte=1"`
	Title string `form:"title" binding:"required,min=3,max=100"`
	Desc string `form:"desc" binding:"required,min=3,max=255"`
	CoverImageUrl string `form:"cover_image_url" binding:"required,url"`
	Content string `form:"content" binding:"required,min=3,max=4294967295"`
	CreateBy string `form:"create_by" binding:"required,min=3,max=100"`
	State uint8 `form:"state,default=1" binding:"oneof=0 1"`
}

type UpdateArticleRequest struct {
	ID uint32 `form:"id" binding:"required,gte=1"`
	Title string `form:"title" binding:"min=3,max=100"`
	Desc string `form:"desc" binding:"min=3,max=255"`
	CoverImageUrl string `form:"cover_image_url" binding:"url"`
	Content string `form:"content" binding:"min=3,max=4294967295"`
	State uint8 `form:"state" binding:"oneof=0 1"`
	ModifiedBy string `form:"modified_by" binding:"required,min=3,max=100"`
}

type DeleteArticleRequest struct {
	ID uint32 `form:"id" binding:"required,gte=1"`
}

type Article struct {
	ID            uint32     `json:"id"`
	Title         string     `json:"title"`
	Desc          string     `json:"desc"`
	Content       string     `json:"content"`
	CoverImageUrl string     `json:"cover_image_url"`
	State         uint8      `json:"state"`
	Tag           *model.Tag `json:"tag"`
}

func (s *Service) CreateArticle(param *CreateArticleRequest) error {
	article := dao.Article{
		Title:         param.Title,
		Desc:          param.Desc,
		Content:       param.Content,
		CoverImageUrl: param.CoverImageUrl,
		CreatedBy:     param.CreateBy,
		State:         param.State,
	}
	s.dao.Begin()
	articleModel, err := s.dao.CreateArticle(&article)
	if err != nil {
		s.dao.Rollback()
		return err
	}

	err = s.dao.CreateArticleTag(articleModel.ID, param.TagId, param.CreateBy)
	if err != nil {
		s.dao.Rollback()
		return err
	}

	s.dao.Commit()
	return nil
}

func (s *Service) GetArticleList(param *ArticleListRequest, pager *app.Pager) ([]*Article, int, error) {
	articleCount, err := s.dao.CountArticleByTagId(param.TagId, param.State)
	if err != nil {
		return nil, 0, err
	}

	articles, err := s.dao.GetArticleListByTagId(param.TagId, param.State, pager.Page, pager.PageSize)
	if err != nil {
		return nil, 0, nil
	}

	var articleList []*Article
	for _, article := range articles {
		articleList = append(articleList, &Article{
			ID:            article.ArticleID,
			Title:         article.ArticleTitle,
			Desc:          article.ArticleDesc,
			Content:       article.Content,
			CoverImageUrl: article.CoverImageUrl,
			Tag:           &model.Tag{Model: &model.Model{ID: article.TagID}, Name: article.TagName},
		})
	}

	return articleList, articleCount, nil
}

func (s *Service) GetArticle(param *ArticleRequest) (*Article, error) {
	article, err := s.dao.GetArticle(param.ID, param.State)
	if err != nil {
		return nil, err
	}

	articleTag, err := s.dao.GetArticleTagByAID(param.ID)
	if err != nil {
		return nil, err
	}

	tag, err := s.dao.GetTag(articleTag.TagId)
	if err != nil {
		return nil, err
	}

	return &Article{
		ID: article.ID,
		Title:         article.Title,
		Desc:          article.Desc,
		Content:       article.Content,
		CoverImageUrl: article.CoverImageUrl,
		Tag:           &model.Tag{Model: &model.Model{ID: tag.ID}, Name: tag.Name},

	}, nil
}

func (s *Service) UpdateArticle(param *UpdateArticleRequest) error {
	return s.dao.UpdateArticle(&dao.Article{
		ID:            param.ID,
		Title:         param.Title,
		Desc:          param.Desc,
		Content:       param.Content,
		CoverImageUrl: param.CoverImageUrl,
		ModifiedBy:    param.ModifiedBy,
		State:         param.State,
	})
}

func (s *Service) DeleteArticle(articleId uint32) error {
	return s.dao.DeleteArticle(articleId)
}
