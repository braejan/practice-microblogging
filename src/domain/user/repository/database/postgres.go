package database

import (
	"context"
	"database/sql"
	"fmt"

	commonEntities "github.com/braejan/practice-microblogging/src/domain/common/entities"
	"github.com/braejan/practice-microblogging/src/domain/common/interfaces"
	userEntities "github.com/braejan/practice-microblogging/src/domain/user/entities"

	_ "github.com/lib/pq"
)

type UserPostgres struct {
	*commonEntities.Postgres
}

func NewUserPostgres() (userPostgres IUserPostgres, err error) {
	userPostgres = &UserPostgres{
		Postgres: commonEntities.NewPostgres(),
	}
	return
}

type IUserPostgres interface {
	CreateUser(userID string, name string) (err error)
	FindUserByID(userID string) (user *userEntities.User, err error)
	interfaces.IPostgres
}

const (
	CreateUserDML = `INSERT INTO "user"(id, name) VALUES ($1, $2)`
)

func (userPostgres *UserPostgres) CreateUser(userID string, name string) (err error) {
	if err != nil {
		return err
	}
	tx, err := userPostgres.BeginTx(context.Background())
	if err != nil {
		return
	}
	defer func() {
		err = userPostgres.DB.Close()
	}()
	_, err = tx.ExecContext(context.Background(), CreateUserDML, userID, name)
	if err != nil {
		_ = userPostgres.RollbackTx(tx)
		return
	}
	err = userPostgres.CommitTx(tx)
	return
}

const (
	FindUserByIDDML = `SELECT id, name FROM "user" WHERE id = $1`
)

func (userPostgres *UserPostgres) FindUserByID(userID string) (user *userEntities.User, err error) {
	userPostgres.DB, err = userPostgres.Open()
	defer func() {
		_ = userPostgres.DB.Close()
	}()
	rows, err := userPostgres.DB.QueryContext(context.Background(), FindUserByIDDML, userID)
	if err != nil {
		return
	}
	defer rows.Close()
	if rows.Next() {
		user = &userEntities.User{}
		err = rows.Scan(&user.ID, &user.Name)
		if err != nil {
			user = nil
			return
		}
	}
	return
}

func (userPostgres *UserPostgres) Open() (DB *sql.DB, err error) {
	sqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		commonEntities.Host,
		commonEntities.Port,
		commonEntities.DBUser,
		commonEntities.DBSecret,
		commonEntities.DBName,
	)
	DB, err = sql.Open("postgres", sqlInfo)
	return
}

func (userPostgres *UserPostgres) BeginTx(ctx context.Context) (tx *sql.Tx, err error) {
	userPostgres.DB, err = userPostgres.Open()
	if err != nil {
		return
	}
	tx, err = userPostgres.DB.BeginTx(ctx, nil)
	return
}
func (userPostgres *UserPostgres) CommitTx(tx *sql.Tx) (err error) {
	err = tx.Commit()
	return
}
func (userPostgres *UserPostgres) RollbackTx(tx *sql.Tx) (err error) {
	err = tx.Rollback()
	return
}
