package userservice

import (
	"fmt"
	"gamegolang/entity"
	"gamegolang/pkg/phone_number"
)

type Repository interface {
	IsPhoneNumberUnique(phoneNumber string) (bool, error)
	Register(u entity.User) (entity.User, error)
}

type Service struct {
	Repo Repository
}

type RegisterRequest struct {
	Name        string
	PhoneNumber string
}

type RegisterResponse struct {
	User entity.User
}

func (s Service) Register(req RegisterRequest) (RegisterResponse, error) {
	// TODO - we should verify phone number by verification code

	// validate phone number
	if !phone_number.IsValidPhoneNumber(req.PhoneNumber) {
		
		return RegisterResponse{}, fmt.Errorf("phone number is not valid")
	}

	// check uniqueness of phone number
	if isUnique, err := s.Repo.IsPhoneNumberUnique(req.PhoneNumber); err != nil || !isUnique {
		if err != nil {
			return RegisterResponse{}, fmt.Errorf("unexpected error: %w", err)
		}

		if !isUnique {
			return RegisterResponse{}, fmt.Errorf("phone number is not unique")
		}
	}

	// validate name
	if len(req.Name) < 3 {
		return RegisterResponse{}, fmt.Errorf("name length should be greater than 3")
	}

	user := entity.User{
		ID:          0,
		PhoneNumber: req.PhoneNumber,
		Name:        req.Name,
	}

	// create new user in storage
	createdUser, err := s.Repo.Register(user)
	if err != nil {
		return RegisterResponse{}, fmt.Errorf("unexpected error: %w", err)
	}

	// return created user
	return RegisterResponse{
		User: createdUser,
	}, nil
}
