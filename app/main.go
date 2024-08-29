package main

import (
	"database/sql"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	_ "github.com/mattn/go-sqlite3"

	"github.com/bxcodec/go-clean-arch/bmi"
	"github.com/bxcodec/go-clean-arch/internal/repository/sqlitez"
	"github.com/bxcodec/go-clean-arch/internal/rest"
	"github.com/bxcodec/go-clean-arch/internal/rest/middleware"
	"github.com/joho/godotenv"
)

const (
	defaultTimeout = 30
	defaultAddress = ":9090"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	dbConn, err := sql.Open(`sqlite3`, "./internal/repository/sqlitez/bmi.db")
	if err != nil {
		log.Fatal("failed to open connection to database", err)
	}
	createTableSQL := `
    CREATE TABLE IF NOT EXISTS bmi_records (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        user_name TEXT NOT NULL,
        weight REAL NOT NULL,
        height REAL NOT NULL,
        bmi REAL NOT NULL,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP
    );
    `
	_, err = dbConn.Exec(createTableSQL)
	if err != nil {
		log.Fatalf("Failed to create table: %v\n", err)
	}

	err = dbConn.Ping()
	if err != nil {
		log.Fatal("failed to ping database ", err)
	}

	defer func() {
		err := dbConn.Close()
		if err != nil {
			log.Fatal("got error when closing the DB connection", err)
		}
	}()

	e := echo.New()
	e.Use(middleware.CORS)

	timeout, err := strconv.Atoi(os.Getenv("CONTEXT_TIMEOUT"))
	if err != nil {
		log.Println("failed to parse timeout, using default timeout")
		timeout = defaultTimeout
	}

	e.Use(middleware.SetRequestContextWithTimeout(time.Duration(timeout) * time.Second))

	// BMI Repo
	bmiRepo := sqlitez.NewBMIRepository(dbConn)

	// BMI Service
	bmiSvc := bmi.NewService(bmiRepo)
	rest.NewBMIHandler(e, bmiSvc)

	// Start Server
	address := os.Getenv("SERVER_ADDRESS")
	if address == "" {
		address = defaultAddress
	}
	log.Fatal(e.Start(address)) //nolint
}
