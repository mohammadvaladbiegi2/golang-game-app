package mysql

import (
	"database/sql"
	"errors"
	"fmt"
	"gamegolang/entity"
	"gamegolang/pkg/richerror"
	userservice "gamegolang/service/user_service"
)

func (d *MySQLDB) IsPhoneNumberUnique(phoneNumber string) (bool, error) {
	query := `select * from user where phone_number = ?`
	user := entity.User{}
	var create_at []uint8
	var update_at []uint8

	result := d.db.QueryRow(query, phoneNumber)
	err := result.Scan(&user.ID, &user.Name, &user.PhoneNumber, &create_at, &update_at, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return true, nil
		}

		return false, err
	}
	fmt.Println(user, string(create_at), string(update_at))
	return false, nil
}

func (d *MySQLDB) Register(u entity.User) (*entity.User, error) {
	res, err := d.db.Exec(`INSERT INTO user (name, phone_number, password) VALUES (?, ?,?)`, u.Name, u.PhoneNumber, u.Password)
	if err != nil {
		return nil, fmt.Errorf("unxepected error => %w", err)
	}

	id, _ := res.LastInsertId()

	return &entity.User{
		ID:          uint(id),
		Name:        u.Name,
		PhoneNumber: u.PhoneNumber,
		Password:    u.Password,
	}, nil
}

func (d *MySQLDB) FindUserDataByPhoneNumber(phoneNumber string) (*entity.User, error) {
	query := `select * from user where phone_number = ?`
	user := entity.User{}
	var create_at []uint8
	var update_at []uint8
	result := d.db.QueryRow(query, phoneNumber)
	err := result.Scan(&user.ID, &user.Name, &create_at, &update_at, &user.PhoneNumber, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}

		return nil, fmt.Errorf("server Error %v", err)
	}

	return &user, nil
}

func (d *MySQLDB) GetProfileByID(userID uint) (*userservice.GetProfileResponse, error) {
	query := `select name from user where id = ?`
	userName := userservice.GetProfileResponse{}
	result := d.db.QueryRow(query, userID)
	err := result.Scan(&userName.Name)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, richerror.NewError(richerror.RichError{
				WrappedError: err,
				StatusCode:   404,
				Message:      "user not found",
				MetaData:     nil,
			})
		}

		return nil, richerror.NewError(richerror.RichError{
			WrappedError: err,
			StatusCode:   500,
			Message:      "server error",
			MetaData: map[string]interface{}{
				"method": "GetProfileByID",
			},
		})
	}

	return &userName, nil
}
