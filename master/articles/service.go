package articles

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"kedaiprogrammer/helpers"
	"net/http"
	"strings"

	"github.com/aws/aws-sdk-go/aws/session"
)

type Services interface {
	GetAll(tag, search string, limit int, offset int, OrderColumn string, orderDirection string) ([]map[string]interface{}, int, int, error)
	Save(s3 *session.Session, request CreateRequest, uuid string) (Article, error)
	GetOne(articleID string) (map[string]interface{}, error)
}

type services struct {
	repository Repository
}

func NewServices(repository Repository) *services {
	return &services{repository}
}
func (s *services) GetAll(tag, search string, limit int, offset int, OrderColumn string, orderDirection string) ([]map[string]interface{}, int, int, error) {
	return s.repository.GetAllWithCounts(tag, search, limit, offset, OrderColumn, orderDirection)
}
func (s *services) GetOne(articleID string) (map[string]interface{}, error) {
	return s.repository.GetOne(articleID)
}

func (s *services) Save(s3 *session.Session, request CreateRequest, uuid string) (Article, error) {
	article := Article{}
	article.ArticleID = helpers.GenerateUUID()
	article.AuthorID = uuid
	article.Body = request.Body
	article.CategoryID = request.CategoryID
	article.Description = request.Description
	decodedImage, _ := base64.StdEncoding.DecodeString(request.MainImage)
	file := bytes.NewReader(decodedImage)
	fileType := http.DetectContentType(decodedImage)
	if fileType != "image/jpeg" && fileType != "image/jpg" && fileType != "image/png" {
		return article, fmt.Errorf("format image not permitted: %v", fileType)
	}
	image := fmt.Sprintf("%s.%s", request.Title, strings.Split(fileType, "/")[1])
	article.MainImage = strings.ReplaceAll(strings.ToLower(image), " ", "-")
	article.Slug = strings.ReplaceAll(strings.ToLower(request.Title), " ", "-")
	article.Status = 1
	article.Title = request.Title

	err := helpers.SaveObjectToS3(s3, "kedaiedukasi/articles/", strings.ReplaceAll(strings.ToLower(image), " ", "-"), file)
	if err != nil {
		return article, err
	}
	newArticle, err := s.repository.Save(article)
	if err != nil {
		return newArticle, err
	}
	return newArticle, nil
}
