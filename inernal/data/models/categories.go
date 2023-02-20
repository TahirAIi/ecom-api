package data

import (
	"database/sql"
	"time"
)

type Category struct {
	Id          int64     `json:"-"`
	Uuid        string    `json:"uuid"`
	ParentId    *int      `json:"parentId"`
	Title       string    `json:"title"`
	Description *string   `json:"description"`
	CreatedAt   time.Time `json:"-"`
	UpdatedAt   time.Time `json:"-"`
	DeletedAt   time.Time `json:"-"`
}

type CategoryModel struct {
	Db *sql.DB
}

func (categoryModel CategoryModel) Insert(category *Category) error {
	query := `INSERT INTO categories (uuid, title, description, parent_id, created_at, updated_at)
                VALUES(?, ?, ?, ?, ?, ?)`
	args := []interface{}{category.Uuid, category.Title, category.Description, category.ParentId, category.CreatedAt, category.UpdatedAt}
	_, err := categoryModel.Db.Exec(query, args...)
	if err != nil {
		return err
	}
	return nil
}

func (categoryModel CategoryModel) Update(category *Category) error {
	query := `UPDATE categories
    SET title = ?, description = ?
    where uuid = ? LIMIT 1`

	_, err := categoryModel.Db.Exec(query, category.Title, category.Description, category.Uuid)
	if err != nil {
		return err
	}
	return nil
}

func (categoryModel CategoryModel) Delete(uuid string) error {
	query := `DELETE FROM categories 
    WHERE uuid = ?LIMIT 1`
	_, err := categoryModel.Db.Exec(query, uuid)

	if err != nil {
		return err
	}
	return nil
}

func (categoryModel CategoryModel) GetAll(limit int, offset int) ([]*Category, error) {
	query := `SELECT count(id), uuid, parent_id, title, description
    FROM categories
     where deleted_at IS NULL 
      GROUP BY id LIMIT ? OFFSET ?`
	rows, err := categoryModel.Db.Query(query, limit, offset)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	categories := []*Category{}
	totalCategories := 0
	for rows.Next() {
		var category Category

		err := rows.Scan(
			&totalCategories,
			&category.Uuid,
			&category.ParentId,
			&category.Title,
			&category.Description,
		)

		if err != nil {
			return nil, err
		}
		categories = append(categories, &category)
	}
	return categories, nil
}

func (categoryModel CategoryModel) Get(uuid string) (*Category, error) {
	query := `SELECT uuid, title, parent_id, description FROM categories 
            where uuid = ?`

	var category Category

	err := categoryModel.Db.QueryRow(query, uuid).Scan(&category.Uuid, &category.Title, &category.ParentId, &category.Description)
	if err != nil {
		return nil, err
	}

	return &category, nil
}

func (categoryModel CategoryModel) GetTotalCount() (int, error) {
	query := `SELECT COUNT(id) FROM categories 
            where deleted_at IS NULL`

	var totalCount int

	err := categoryModel.Db.QueryRow(query).Scan(&totalCount)
	if err != nil {
		return 0, err
	}
	return totalCount, nil
}
