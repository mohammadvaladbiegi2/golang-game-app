package mysql

import (
	"database/sql"
	"fmt"
	"gamegolang/entity"
)

func (d *MySQLDB) IsPhoneNumberUnique(phoneNumber string) (bool, error) {
	query := `select * from user where phone_number = ?`
	user := entity.User{}
	var create_at []uint8
	var update_at []uint8

	result := d.db.QueryRow(query, phoneNumber)
	err := result.Scan(&user.ID, &user.Name, &user.PhoneNumber, &create_at, &update_at)
	if err != nil {
		if err == sql.ErrNoRows {
			return true, nil
		}

		return false, err
	}
	fmt.Println(user, string(create_at), string(update_at))
	return false, nil
}

func (d *MySQLDB) Register(u entity.User) (entity.User, error) {
	res, err := d.db.Exec(`INSERT INTO user (name, phone_number) VALUES (?, ?)`, u.Name, u.PhoneNumber)
	if err != nil {
		return entity.User{}, fmt.Errorf("unxepected error => %w", err)
	}

	id, _ := res.LastInsertId()

	return entity.User{
		ID:          uint(id),
		Name:        u.Name,
		PhoneNumber: u.PhoneNumber,
	}, nil
}
