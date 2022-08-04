package main

import (
	"context"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/braejan/practice-microblogging/src/domain/user/entities"
	"github.com/braejan/practice-microblogging/src/domain/user/usecases"
)

func HandleRequest(ctx context.Context, request entities.User) (err error) {
	userUsesCases, err := usecases.NewUserUsecases()
	if err != nil {
		return
	}
	err = userUsesCases.CreateUser(request.ID, request.Name)
	return err
}

func main() {
	lambda.Start(HandleRequest)
}

//GOOS=linux go build -o create_user src/app/lambda/user/create_user.go
