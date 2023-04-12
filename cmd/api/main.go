package main

// Package classification of ecom-api.
//     Host: localhost
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
// swagger:meta

//go:generate swagger generate spec

import (
	"database/sql"
	data "ecom-api/inernal/data/models"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

type application struct {
	logger *log.Logger
	models data.Models
	config config
}

type config struct {
	multipartFormSize int
	fileBaseUrl       string
}

func main() {
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Llongfile)

	err := godotenv.Load()
	if err != nil {
		logger.Println(os.Getwd())
		logger.Fatal(err)
	}

	db, err := openDbConnection()

	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()
	config := config{
		multipartFormSize: 2048,
	}

	app := application{
		logger: logger,
		models: data.NewModel(db),
		config: config,
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
