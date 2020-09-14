package persistence

import (
	"Repository-Pattern/domain/model"
	"Repository-Pattern/domain/repositories"
	"errors"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"os"
	"strings"
)

type PostRepo struct {
	db *gorm.DB
}

func NewPostRepository(db *gorm.DB) *PostRepo {
	return &PostRepo{db}
}

//PostRepo implements the repository.PostRepository interface
var _ repositories.PostRepository = &PostRepo{}

func (r *PostRepo) SavePost(post *model.Post) (*model.Post, map[string]string) {
	dbErr := map[string]string{}
	//The images are uploaded to digital ocean spaces. So we need to prepend the url. This might not be your use case, if you are not uploading image to Digital Ocean.
	post.PostImage = os.Getenv("DO_SPACES_URL") + post.PostImage

	err := r.db.Debug().Create(&post).Error
	if err != nil {
		//since our title is unique
		if strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "Duplicate") {
			dbErr["unique_title"] = "Post title already taken"
			return nil, dbErr
		}
		//any other db error
		dbErr["db_error"] = "database error"
		return nil, dbErr
	}
	return post, nil
}

func (r *PostRepo) GetPost(id uint64) (*model.Post, error) {
	var Post model.Post
	err := r.db.Debug().Where("id = ?", id).Take(&Post).Error
	if err != nil {
		return nil, errors.New("database error, please try again")
	}
	if gorm.IsRecordNotFoundError(err) {
		return nil, errors.New("post not found")
	}
	return &Post, nil
}

func (r *PostRepo) GetPostByIdUser(userUuid uuid.UUID) ([]model.Post, error) {
	var Post []model.Post
	err := r.db.Debug().Where("user_uuid = ?", userUuid).Find(&Post).Error
	if err != nil {
		return nil, errors.New("database error, please try again")
	}
	if gorm.IsRecordNotFoundError(err) {
		return nil, errors.New("post not found")
	}
	return Post, nil
}

func (r *PostRepo) GetAllPost() ([]model.Post, error) {
	var Posts []model.Post
	err := r.db.Debug().Limit(100).Order("created_at desc").Find(&Posts).Error
	if err != nil {
		return nil, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return nil, errors.New("post not found")
	}
	return Posts, nil
}

func (r *PostRepo) UpdatePost(post *model.Post) (*model.Post, map[string]string) {
	dbErr := map[string]string{}
	err := r.db.Debug().Save(&post).Error
	if err != nil {
		//since our title is unique
		if strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "Duplicate") {
			dbErr["unique_title"] = "title already taken"
			return nil, dbErr
		}
		//any other db error
		dbErr["db_error"] = "database error"
		return nil, dbErr
	}
	return post, nil
}

func (r *PostRepo) DeletePost(id uint64) error {
	var Post model.Post
	err := r.db.Debug().Where("id = ?", id).Delete(&Post).Error
	if err != nil {
		return errors.New("database error, please try again")
	}
	return nil
}
