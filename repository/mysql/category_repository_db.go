package mysql

import (
	"fmt"
	"gamegolang/entity"
	categoryservice "gamegolang/service/category_service"
)

func (d *MySQLDB) Create(req categoryservice.CreateRequest) (*entity.Category, error) {

	query := `INSERT INTO category (title, description) VALUES (?, ?)`

	result, err := d.db.Exec(query, req.Title, req.Description)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	id, LastInsertIdID := result.LastInsertId()
	if LastInsertIdID != nil {
		fmt.Println("LastInsertIdID:", LastInsertIdID)
		return nil, LastInsertIdID
	}

	return &entity.Category{
		uint(id),
		req.Title,
		req.Description,
	}, nil

}
