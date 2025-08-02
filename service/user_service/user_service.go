package userservice

// TODO add jwt

import (
	"gamegolang/entity"
	"gamegolang/pkg/jwt"
	"gamegolang/pkg/phone_number"
	"gamegolang/pkg/richerror"
	"hash/fnv"
	"strconv"
)

type RegisterResponse struct {
	User entity.User
}

type RegisterRepository interface {
	IsPhoneNumberUnique(phoneNumber string) (bool, richerror.RichError)
	Register(u entity.User) (*entity.User, richerror.RichError)
}

type LoginCredentials struct {
	PhoneNumber string `json:"phone_number"`
	Password    string `json:"password"`
}

type LoginRepository interface {
	FindUserDataByPhoneNumber(phoneNumber string) (*entity.User, richerror.RichError)
	GetProfileByID(userID uint) (*GetProfileResponse, richerror.RichError)
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

func (s RegisterService) Register(req RegisterRequest) (*entity.User, richerror.RichError) {
	// TODO - we should verify phone number by verification code

	// validate phone number
	if !phone_number.IsValidPhoneNumber(req.PhoneNumber) {

		return nil, richerror.NewError(400, "phone number is not valid")
	}

	// check uniqueness of phone number
	isUnique, Unique := s.Repo.IsPhoneNumberUnique(req.PhoneNumber)
	if Unique.HaveError() || !isUnique {
		if Unique.HaveError() {
			return nil, Unique
		}

		if !isUnique {
			return nil, richerror.NewError(400, "phone number is not unique")
		}
	}

	// validate name
	if len(req.Name) < 3 {
		return nil, richerror.NewError(400, "name length should be greater than 3")
	}
	if len(req.Password) < 6 {
		return nil, richerror.NewError(400, "password length should be greater than 8")
	}

	user := entity.User{
		ID:          0,
		PhoneNumber: req.PhoneNumber,
		Name:        req.Name,
		Password:    strconv.Itoa(hash(req.Password)),
	}

	// create new user in storage
	createdUser, err := s.Repo.Register(user)
	if err.HaveError() {
		return nil, err
	}

	// return created user
	return createdUser, richerror.RichError{}
}

func (s LoginService) Login(req LoginCredentials) (string, richerror.RichError) {

	// TODO add rich error to project
	userData, FindPhoneError := s.Repo.FindUserDataByPhoneNumber(req.PhoneNumber)
	if FindPhoneError.HaveError() {
		return "", FindPhoneError
	}

	if userData.Password != strconv.Itoa(hash(req.Password)) {
		return "", richerror.NewError(403, "password or phone number does not match")
	}

	token, tError := jwt.BuildToken(userData.Name, userData.ID)
	if tError.HaveError() {
		return "", tError
	}

	return token, richerror.RichError{}
}

func (s LoginService) GetProfile(userID uint) (*GetProfileResponse, richerror.RichError) {

	if userID <= 0 {
		return nil, richerror.NewError(
			nil,
			400,
			"user ID is requir",
			nil,
		)
	}

	username, uError := s.Repo.GetProfileByID(userID)
	if uError.HaveError() {
		return nil, uError
	}

	return username, richerror.RichError{}
}
