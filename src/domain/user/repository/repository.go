package repository

import (
	"github.com/braejan/practice-microblogging/src/domain/user/entities"
	"github.com/braejan/practice-microblogging/src/domain/user/repository/database"
)

type UserRepository struct {
	repository database.IUserPostgres
}
type IUserRepository interface {
	CreateUser(ID string, name string) (err error)
	FindUserByID(ID string) (user *entities.User, err error)
}

func NewUserRepository() (repository IUserRepository, err error) {
	userPostgresRepository, err := database.NewUserPostgres()
	if err != nil {
		return
	}
	repository = &UserRepository{
		repository: userPostgresRepository,
	}
	return
}

func (userRepository *UserRepository) CreateUser(ID string, name string) (err error) {
	return userRepository.repository.CreateUser(ID, name)
}

func (userRepository *UserRepository) FindUserByID(ID string) (user *entities.User, err error) {
	return userRepository.repository.FindUserByID(ID)
}
