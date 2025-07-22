package userservice

// TODO add jwt

import (
	"fmt"
	"gamegolang/entity"
	"gamegolang/pkg/jwt"
	"gamegolang/pkg/phone_number"
	"hash/fnv"
	"strconv"
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
	FindUserDataByPhoneNumber(phoneNumber string) (*entity.User, error)
	GetProfileByID(userID uint) (*GetProfileResponse, error)
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

type GetProfileResponse struct {
	Name string `json:"name"`
}

func hash(s string) int {
	h := fnv.New32a()
	h.Write([]byte(s))
	return int(h.Sum32())
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
		Password:    strconv.Itoa(hash(req.Password)),
	}

	// create new user in storage
	createdUser, err := s.Repo.Register(user)
	if err != nil {
		return nil, fmt.Errorf("unexpected error: %w", err)
	}

	// return created user
	return createdUser, nil
}

func (s LoginService) Login(req LoginCredentials) (string, error) {

	userData, FindPhoneError := s.Repo.FindUserDataByPhoneNumber(req.PhoneNumber)
	if FindPhoneError != nil {
		return "", FindPhoneError
	}

	if userData.Password != strconv.Itoa(hash(req.Password)) {
		return "", fmt.Errorf("password or phone number does not match")
	}

	token, tError := jwt.BuildToken(userData.Name, userData.ID)
	if tError != nil {
		return "", fmt.Errorf("error creating token")
	}

	return token, nil
}

func (s LoginService) GetProfile(userID uint) (*GetProfileResponse, error) {

	if userID <= 0 {
		return nil, fmt.Errorf("user ID is required")
	}

	username, uError := s.Repo.GetProfileByID(userID)
	if uError != nil {
		return nil, uError
	}

	return username, nil
}
