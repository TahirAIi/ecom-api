package main

import (
	"database/sql"
	data "ecom-api/inernal/data/models"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"os"
)

type application struct {
	logger *log.Logger
	models data.Models
}

func main() {
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
	db, err := openDbConnection()

	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()

	app := application{
		logger: logger,
		models: data.NewModel(db),
	}

	app.server()
}

func openDbConnection() (*sql.DB, error) {
	db, err := sql.Open("mysql", "dev:tahirdev@tcp(localhost:3306)/api-ecom")

	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
