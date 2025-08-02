package mysql

import (
	"database/sql"
	"gamegolang/entity"
	"gamegolang/pkg/richerror"
	userservice "gamegolang/service/user_service"
	"time"
)

func (d *MySQLDB) IsPhoneNumberUnique(phoneNumber string) (bool, richerror.RichError) {
	query := `select * from user where phone_number = ?`
	user := entity.User{}
	var create_at time.Time
	var update_at time.Time

	result := d.db.QueryRow(query, phoneNumber)
	err := result.Scan(&user.ID, &user.Name, &user.PhoneNumber, &create_at, &update_at, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return true, richerror.RichError{}
		}

		return false, richerror.NewError(
			err,
			500,
			"server error",
			map[string]interface{}{
				"method": "IsPhoneNumberUnique",
			},
		)
	}
	return false, richerror.RichError{}
}

func (d *MySQLDB) Register(u entity.User) (*entity.User, richerror.RichError) {
	res, err := d.db.Exec(`INSERT INTO user (name, phone_number, password) VALUES (?, ?,?)`, u.Name, u.PhoneNumber, u.Password)
	if err != nil {
		return nil, richerror.NewError(
			err,
			500,
			"server error",
			map[string]interface{}{
				"method": "Register",
			},
		)
	}

	id, _ := res.LastInsertId()

	return &entity.User{
		ID:          uint(id),
		Name:        u.Name,
		PhoneNumber: u.PhoneNumber,
		Password:    u.Password,
	}, richerror.RichError{}
}

func (d *MySQLDB) FindUserDataByPhoneNumber(phoneNumber string) (*entity.User, richerror.RichError) {
	query := `select * from user where phone_number = ?`
	user := entity.User{}
	var create_at time.Time
	var update_at time.Time
	result := d.db.QueryRow(query, phoneNumber)
	err := result.Scan(&user.ID, &user.Name, &create_at, &update_at, &user.PhoneNumber, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {

			return nil, richerror.NewError(
				err,
				404,
				"user not found",
				nil,
			)
		}

		return nil, richerror.NewError(
			err,
			500,
			"server error",
			map[string]interface{}{
				"method": "GetProfileByID",
			},
		)
	}

	return &user, richerror.RichError{}
}

func (d *MySQLDB) GetProfileByID(userID uint) (*userservice.GetProfileResponse, richerror.RichError) {
	query := `select name from user where id = ?`
	userName := userservice.GetProfileResponse{}
	result := d.db.QueryRow(query, userID)
	err := result.Scan(&userName.Name)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, richerror.NewError(
				err,
				404,
				"user not found",
				nil,
			)
		}

		return nil, richerror.NewError(
			err,
			500,
			"server error",
			map[string]interface{}{
				"method": "GetProfileByID",
			},
		)
	}

	return &userName, richerror.RichError{}
}
