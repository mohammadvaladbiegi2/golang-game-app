package userservice

import (
	"fmt"
	"gamegolang/entity"
	"gamegolang/pkg/phone_number"
)

type RegisterResponse struct {
	User entity.User
}

type Repository interface {
	IsPhoneNumberUnique(phoneNumber string) (bool, error)
	Register(u entity.User) (*entity.User, error)
	Login(u LoginRequest) (*bool, error)
}

type Service struct {
	Repo Repository
}

type RegisterRequest struct {
	Name        string `json:"name"`
	PhoneNumber string `json:"phone_number"`
	Password    string `json:"password"`
}

type LoginRequest struct {
	PhoneNumber string `json:"phone_number"`
	Password    string `json:"password"`
}

func (s Service) Register(req RegisterRequest) (*entity.User, error) {
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

func (s Service) Login(req LoginRequest) (*bool, error) {
	//	TODO check phone number exist and get User
	//	TODO check password
	//	TODO return status
	status := new(bool)
	*status = true
	return status, nil
}
