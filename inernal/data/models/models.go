package data

import "database/sql"

type Models struct {
	Category CategoryModel
}

func NewModel(db *sql.DB) Models {
	return Models{
		Category: CategoryModel{Db: db},
	}
}
