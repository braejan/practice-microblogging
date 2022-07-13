package usecases

import (
	"github.com/braejan/practice-microblogging/src/domain/user/entities"
	"github.com/braejan/practice-microblogging/src/domain/user/repository"
)

type UserUsecases struct {
	userRepository repository.IUserRepository
}

func NewUserUsecases() (usecases *UserUsecases, err error) {
	repository, err := repository.NewUserRepository()
	usecases = &UserUsecases{
		userRepository: repository,
	}
	return
}

func (usecases *UserUsecases) CreateUser(ID string, name string) (err error) {
	return usecases.userRepository.CreateUser(ID, name)
}

func (usecases *UserUsecases) FindUserByID(ID string) (user *entities.User, err error) {
	return usecases.userRepository.FindUserByID(ID)
}
