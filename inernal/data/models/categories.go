package data

import (
	"database/sql"
	"time"
)

type Category struct {
	Id          int32     `json:"-"`
	ParentId    *int32      `json:"parentId"`
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
	query := `INSERT INTO categories (title, description, parent_id)
                VALUES(?, ?, ?)`
	args := []interface{}{category.Title, category.Description, category.ParentId}
	_, err := categoryModel.Db.Exec(query, args...)
	if err != nil {
		return err
	}
	return nil
}

func (categoryModel CategoryModel) Update(category *Category) error {
	query := `UPDATE categories
    SET title = ?, description = ?
    where id = ? LIMIT 1`

	_, err := categoryModel.Db.Exec(query, category.Title, category.Description, category.Id)
	if err != nil {
		return err
	}
	return nil
}

func (categoryModel CategoryModel) Delete(id int32) error {
	query := `DELETE FROM categories 
    WHERE id = ?
	LIMIT 1`
	_, err := categoryModel.Db.Exec(query, id)

	if err != nil {
		return err
	}
	return nil
}

func (categoryModel CategoryModel) GetAll(limit int32, offset int32) ([]*Category, error) {
	query := `SELECT id, parent_id, title, description
    FROM categories
    where deleted_at IS NULL 
    LIMIT ? OFFSET ?`
	rows, err := categoryModel.Db.Query(query, limit, offset)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	categories := []*Category{}
	for rows.Next() {
		var category Category

		err := rows.Scan(
			&category.Id,
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

func (categoryModel CategoryModel) Get(id int32) (*Category, error) {
	query := `SELECT id, title, parent_id, description FROM categories 
            where id = ?`

	var category Category

	err := categoryModel.Db.QueryRow(query, id).Scan(&category.Id, &category.Title, &category.ParentId, &category.Description)
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

func (categoryModel CategoryModel) GetCategoryId(id int) (*Category, error) {
	query := `SELECT id FROM categories 
            where id = ?`

	var category Category

	err := categoryModel.Db.QueryRow(query, id).Scan(&category.Id)
	if err != nil {
		return nil, err
	}

	return &category, nil
}
