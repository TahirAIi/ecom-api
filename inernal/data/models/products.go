package data

import (
	"database/sql"
	"time"
)

type Product struct {
	Id          int32     `json:"id"`
	CategoryId  int32     `json:"category_id"`
	Title       string    `json:"title"`
	Description *string   `json:"description"`
	Price       int32       `json:"price"`
	MainPicture string    `json:"main_picture"`
	CreatedAt   time.Time `json:"-"`
	UpdatedAt   time.Time `json:"-"`
	DeletedAt   time.Time `json:"-"`
}

type ProductModel struct {
	Db *sql.DB
}

func (productModel ProductModel) Insert(product *Product) error {
	query := `INSERT INTO products (category_id, title, price,description, main_picture)
	VALUES(?,?,?,?,?)`

	result , err := productModel.Db.Exec(query, product.CategoryId, product.Title, product.Price, product.Description, product.MainPicture)

	if err != nil {
		return err
	}
	val, err := result.LastInsertId()
	if err != nil {
		return err
	}
	product.Id = int32(val)
	return nil
}

func (productModel ProductModel) Get(categoryId int32, productId int32, product *Product) error {
	query := `SELECT id, category_id, title, price, main_picture, description
	FROM products
	WHERE category_id = ? AND id = ?
	LIMIT 1`
	err := productModel.Db.QueryRow(query, categoryId, productId).Scan(&product.Id, &product.CategoryId, &product.Title, &product.Price, &product.MainPicture, &product.Description)
	if err != nil {
		return err
	}
	return nil
}

func (productModel ProductModel) GetAll(categoryId int32, limit int32, offset int32) ([]*Product, error) {
	query := `SELECT id, category_id, title, description, price, main_picture
		FROM products
		WHERE category_id = ?
		LIMIT ? OFFSET ?`

	var products []*Product
	rows, err := productModel.Db.Query(query, categoryId, limit, offset)
	if err != nil {
		return products, err
	}

	for rows.Next() {
		var product Product
		err := rows.Scan(
			&product.Id,
			&product.CategoryId,
			&product.Title,
			&product.Description,
			&product.Price,
			&product.MainPicture,
		)
		if err != nil {
			return nil, err
		}
		products = append(products, &product)
	}
	return products, nil
}

func (productModel ProductModel) Update(caegoryId int32, product *Product) error {
	query := `UPDATE products SET
		category_id = ?, title = ?, description = ?, main_picture = ?
		WHERE category_id = ? AND id = ?
		LIMIT 1 		
		`

	_, err := productModel.Db.Exec(query, product.CategoryId, product.Title, product.Description, product.MainPicture, caegoryId, product.Id)

	if err != nil {
		return err
	}

	return nil
}

func (productModel ProductModel) Delete(categoryId int32, productId int32) error {
	query := `UPDATE products 
			SET deleted_at = ?
			WHERE category_id = ? AND id = ?`

	_, err := productModel.Db.Exec(query, time.Now(), categoryId, productId)
	if err != nil {
		return err
	}
	return nil
}

func (productModel ProductModel) Count(categoryId int) (*int32, error) {
	query := `SELECT COUNT(id) FROM products
	where category_id = ? AND deleted_at IS NULL`

	var count *int32
	err := productModel.Db.QueryRow(query, categoryId).Scan(&count)
	if err != nil {
		return nil, err
	}
	return count, nil
}
