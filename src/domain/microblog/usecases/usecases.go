package usecases

import (
	"fmt"

	commonUsecases "github.com/braejan/practice-microblogging/src/domain/common/usecases"
	"github.com/braejan/practice-microblogging/src/domain/microblog/entities"
	"github.com/braejan/practice-microblogging/src/domain/microblog/repository"
	userEntities "github.com/braejan/practice-microblogging/src/domain/user/entities"
	userUsecases "github.com/braejan/practice-microblogging/src/domain/user/usecases"
)

type MicroblogUsecases struct {
	microblogRepository repository.IMicroblogRepository
	userUsecases        *userUsecases.UserUsecases
	commonUsecases      *commonUsecases.CommonUsecases
}

func NewMicroblogUsecases(userUsecases *userUsecases.UserUsecases, commonUsecases *commonUsecases.CommonUsecases) (usecases *MicroblogUsecases, err error) {
	repository, err := repository.NewMicroblogRepository()
	if err != nil {
		return
	}

	usecases = &MicroblogUsecases{
		microblogRepository: repository,
		userUsecases:        userUsecases,
		commonUsecases:      commonUsecases,
	}
	return
}

func (microblogUsecases *MicroblogUsecases) CreatePost(owner string, text string) (err error) {
	user, err := microblogUsecases.userUsecases.FindUserByID(owner)
	if err != nil {
		return
	}
	err = microblogUsecases.validateUser(user)
	if err != nil {
		return
	}
	err = microblogUsecases.commonUsecases.ValidatePostLength(text)
	if err != nil {
		return
	}
	return microblogUsecases.microblogRepository.CreatePost(user, text)
}

func (microblogUsecases *MicroblogUsecases) GetPostByID(ID string, userID string, internal bool) (post *entities.MicroBlog, err error) {
	user, err := microblogUsecases.userUsecases.FindUserByID(userID)
	if err != nil {
		return
	}
	err = microblogUsecases.validateUser(user)
	if err != nil {
		return
	}
	return microblogUsecases.microblogRepository.GetPostByID(ID, internal)
}

func (microblogUsecases *MicroblogUsecases) GetAllPosts() (posts []*entities.MicroBlog, err error) {
	return microblogUsecases.microblogRepository.GetAllPosts()
}

func (microblogUsecases *MicroblogUsecases) GetAllPostsByUserID(userID string) (posts []*entities.MicroBlog, err error) {
	user, err := microblogUsecases.userUsecases.FindUserByID(userID)
	if err != nil {
		return
	}
	err = microblogUsecases.validateUser(user)
	if err != nil {
		return
	}
	return microblogUsecases.microblogRepository.GetAllPostsByUserID(userID)
}

func (microblogUsecases *MicroblogUsecases) LikePost(ID string, userID string) (err error) {
	user, err := microblogUsecases.userUsecases.FindUserByID(userID)
	if err != nil {
		return
	}
	err = microblogUsecases.validateUser(user)
	if err != nil {
		return
	}
	post, err := microblogUsecases.GetPostByID(ID, userID, true)
	if err != nil {
		return
	}
	if post == nil {
		err = fmt.Errorf("post id %s not found", ID)
		return
	}
	return microblogUsecases.microblogRepository.LikePost(ID, user)
}

func (microblogUsecases *MicroblogUsecases) DislikePost(ID string, userID string) (err error) {
	user, err := microblogUsecases.userUsecases.FindUserByID(userID)
	if err != nil {
		return
	}
	_, err = microblogUsecases.GetPostByID(ID, userID, true)
	if err != nil {
		return
	}
	return microblogUsecases.microblogRepository.DislikePost(ID, user)
}

func (microblogUsecases *MicroblogUsecases) validateUser(user *userEntities.User) (err error) {

	if user == nil {
		err = fmt.Errorf("user %s not found", user.ID)
	}
	return
}
