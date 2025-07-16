package userservice

import (
	"fmt"
	"gamegolang/entity"
	"gamegolang/pkg/phone_number"
)

type RegisterResponse struct {
	User entity.User
}

type RegisterRepository interface {
	IsPhoneNumberUnique(phoneNumber string) (bool, error)
	Register(u entity.User) (*entity.User, error)
}

type LoginCredentials struct {
	PhoneNumber string `json:"phone_number"`
	Password    string `json:"password"`
}

type LoginRepository interface {
	FindUserDataByPhoneNumber(phoneNumber string) (*LoginCredentials, error)
}

type RegisterService struct {
	Repo RegisterRepository
}

type LoginService struct {
	Repo LoginRepository
}

type RegisterRequest struct {
	Name        string `json:"name"`
	PhoneNumber string `json:"phone_number"`
	Password    string `json:"password"`
}

func (s RegisterService) Register(req RegisterRequest) (*entity.User, error) {
	// TODO - we should verify phone number by verification code

	// validate phone number
	if !phone_number.IsValidPhoneNumber(req.PhoneNumber) {

		return nil, fmt.Errorf("phone number is not valid")
	}

	// check uniqueness of phone number
	if isUnique, err := s.Repo.IsPhoneNumberUnique(req.PhoneNumber); err != nil || !isUnique {
		if err != nil {
			return nil, fmt.Errorf("unexpected error: %w", err)
		}

		if !isUnique {
			return nil, fmt.Errorf("phone number is not unique")
		}
	}

	// validate name
	if len(req.Name) < 3 {
		return nil, fmt.Errorf("name length should be greater than 3")
	}
	if len(req.Password) < 6 {
		return nil, fmt.Errorf("password length should be greater than 8")
	}

	user := entity.User{
		ID:          0,
		PhoneNumber: req.PhoneNumber,
		Name:        req.Name,
		Password:    req.Password,
	}

	// create new user in storage
	createdUser, err := s.Repo.Register(user)
	if err != nil {
		return nil, fmt.Errorf("unexpected error: %w", err)
	}

	// return created user
	return createdUser, nil
}

func (s LoginService) Login(req LoginCredentials) (*bool, error) {

	userData, FindPhoneError := s.Repo.FindUserDataByPhoneNumber(req.PhoneNumber)
	status := new(bool)
	*status = false
	if FindPhoneError != nil {
		return status, FindPhoneError
	}

	if userData.Password != req.Password {
		*status = false
		return status, fmt.Errorf("password or phone number does not match")
	}

	*status = true
	return status, nil
}
