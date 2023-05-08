package data

import "database/sql"

type Models struct {
	Category CategoryModel
	Product  ProductModel
	User     UserModel
}

func NewModel(db *sql.DB) Models {
	return Models{
		Category: CategoryModel{Db: db},
		Product:  ProductModel{Db: db},
		User:     UserModel{Db: db},
	}
}
