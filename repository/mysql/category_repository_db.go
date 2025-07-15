package mysql

import (
	"fmt"
	"gamegolang/entity"
	"gamegolang/service/category_service"
)

func (d *MySQLDB) CreateCategory(req category_service.CreateCategoryRequestStruct) (entity.Category, error) {

	query := `INSERT INTO category (title, description) VALUES (?, ?)`

	result, err := d.db.Exec(query, req.Title, req.Description)
	if err != nil {
		fmt.Println(err)
		return entity.Category{}, err
	}
	id, LastInsertIdID := result.LastInsertId()
	if LastInsertIdID != nil {
		fmt.Println("LastInsertIdID:", LastInsertIdID)
		return entity.Category{}, LastInsertIdID
	}

	return entity.Category{
		uint(id),
		req.Title,
		req.Description,
	}, nil

}
