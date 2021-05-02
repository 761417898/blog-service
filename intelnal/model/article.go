package model

import (
	"blog-service/pkg/app"
	"github.com/jinzhu/gorm"
)

type Article struct {
	*Model
	Title         string `json:"title"`
	Desc          string `json:"desc"`
	Content       string `json:"content"`
	CoverImageURL string `json:"cover_image_url"`
	State         uint8  `json:"state"`
}

type ArticleSwagger struct {
	List  []*Article
	Pager *app.Pager
}

func (article Article) TableName() string {
	return "article"
}

func (article Article) Create(db *gorm.DB) error {
	return db.Create(&article).Error
}

func (article Article) Update(db *gorm.DB) error {
	return db.Model(&Article{}).Where("id = ? AND is_del = ?", article.ID, 0).Update(article).Error
}

func (article Article) Delete(db *gorm.DB) error {
	return db.Where("id = ? AND is_del = ?", article.ID, 0).Delete(&article).Error
}

func (article Article) Get(db *gorm.DB) (*Article, error) {
	var retArticle Article
	err := db.Where("id = ? AND is_del = ?", article.ID, 0).First(&retArticle).Error
	if err != nil {
		return nil, err
	}
	return &retArticle, nil
}

func (article Article) Count(db *gorm.DB) (int, error) {
	var count int
	if article.Title != "" {
		db = db.Where("title = ?", article.Title)
	}
	db = db.Where("state = ?", article.State)
	if err := db.Model(&article).Where("is_del = ?", 0).Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

func (article Article) List(db *gorm.DB, pageOffset, pageSize int) ([]*Article, error) {
	var articles []*Article
	var err error
	if pageOffset >= 0 && pageSize > 0 {
		db = db.Offset(pageOffset).Limit(pageSize)
	}
	if article.Title != "" {
		db = db.Where("title = ?", article.Title)
	}
	db = db.Where("state = ?", article.State)
	if err = db.Where("is_del = ?", 0).Find(&articles).Error; err != nil {
		return nil, err
	}

	return articles, nil
}
