package repository

import (
	"github.com/braejan/practice-microblogging/src/domain/microblog/entities"
	"github.com/braejan/practice-microblogging/src/domain/microblog/repository/database"
	userEntities "github.com/braejan/practice-microblogging/src/domain/user/entities"
)

type MicroblogRepository struct {
	repository database.IMicroblogPostgres
}

type IMicroblogRepository interface {
	CreatePost(owner *userEntities.User, text string) (err error)
	GetPostByID(ID string, internal bool) (post *entities.MicroBlog, err error)
	GetAllPosts() (posts []*entities.MicroBlog, err error)
	GetAllPostsByUserID(userID string) (posts []*entities.MicroBlog, err error)
	LikePost(ID string, user *userEntities.User) (err error)
	DislikePost(ID string, user *userEntities.User) (err error)
}

func NewMicroblogRepository() (repository IMicroblogRepository, err error) {
	microblogPostgresRepository, err := database.NewMicroblogPostgres()
	if err != nil {
		return
	}
	repository = &MicroblogRepository{
		repository: microblogPostgresRepository,
	}
	return
}

func (microblogRepository *MicroblogRepository) CreatePost(owner *userEntities.User, text string) (err error) {
	return microblogRepository.repository.CreatePost(owner, text)
}
func (microblogRepository *MicroblogRepository) GetPostByID(ID string, internal bool) (post *entities.MicroBlog, err error) {
	return microblogRepository.repository.GetPostByID(ID, internal)
}
func (microblogRepository *MicroblogRepository) GetAllPosts() (posts []*entities.MicroBlog, err error) {
	return microblogRepository.repository.GetAllPosts()
}
func (microblogRepository *MicroblogRepository) GetAllPostsByUserID(userID string) (posts []*entities.MicroBlog, err error) {
	return microblogRepository.repository.GetAllPostsByUserID(userID)
}
func (microblogRepository *MicroblogRepository) LikePost(ID string, user *userEntities.User) (err error) {
	return microblogRepository.repository.LikePost(ID, user)
}
func (microblogRepository *MicroblogRepository) DislikePost(ID string, user *userEntities.User) (err error) {
	return microblogRepository.repository.DislikePost(ID, user)
}
