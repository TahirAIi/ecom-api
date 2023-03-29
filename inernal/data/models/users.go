package data

import (
	"database/sql"
	"time"
)

type User struct {
	ID int32 `json:"id"`
	FullName string `json:"fullname"`
	Email string `json:"email"`
	Password string `json:"-"`
	IsAdmin bool `json:"is_admin"`
	CreatedAt time.Time `json:"-"`
	UpatedAt time.Time	`json:"-"`
	DeletedAt time.Time	`json:"-"`
}

type UserModel struct {
	Db *sql.DB
}

func (userModel UserModel) Insert(user *User) error {
	query := `INSERT INTO users (full_name, email, password, is_admin)
			values(?,?,?,?)`
	
	result ,err := userModel.Db.Exec(query, user.FullName, user.Email, user.Password, user.IsAdmin)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	user.ID = int32(id) 
	return nil
}

func (userModel UserModel) GetByEmail(email string) (*User, error) {
	query := `SELECT id, full_name, email, password, is_admin FROM users
	WHERE email = ? AND deleted_at IS NULL LIMIT 1`
	var user User
	err := userModel.Db.QueryRow(query, email).Scan(&user.ID, &user.FullName, &user.Email, &user.Password, &user.IsAdmin)
	if err != nil {
		return nil, err
	}

	return &user, nil
}