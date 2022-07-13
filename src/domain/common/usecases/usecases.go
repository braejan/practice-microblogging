package usecases

import (
	"fmt"
)

type CommonUsecases struct{}

func NewCommonUsecases() (usecases *CommonUsecases) {
	usecases = &CommonUsecases{}
	return
}

const (
	maxLength = 200
)

func (usecases *CommonUsecases) ValidatePostLength(text string) (err error) {
	length := len(text)
	if length == 0 || length > maxLength {
		err = fmt.Errorf("max length allowed %d but got %d", maxLength, length)
	}
	return
}
