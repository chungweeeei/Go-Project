package database

import (
	"fmt"

	"example.com/restful-server/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// declare global database variable
var DB *gorm.DB

func InitDB() {

	createDatabaseIfNotExists("go_db")

	dsn := "host=localhost user=root password=root dbname=go_db port=5432 sslmode=disable TimeZone=Asia/Taipei"
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("Could not connect to database")
	}

	// Do migrate tables
	err = DB.AutoMigrate(&models.User{}, &models.Task{})
	if err != nil {
		panic("Could not auto migrate User model")
	}

	fmt.Println("Database connection established and User model migrated successfully")
}

func createDatabaseIfNotExists(dbName string) {
	// Connect to the default postgres database to check if the target database exists
	var err error
	dsn := "host=localhost user=root password=root dbname=postgres port=5432 sslmode=disable TimeZone=Asia/Taipei"
	tmpDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("Could not connect to default postgres database")
	}

	// Get the underlying sql.DB to close the connection later
	sqlDB, err := tmpDB.DB()
	if err != nil {
		panic("Could not get underlying database from gorm")
	}

	defer sqlDB.Close()

	// Check if the database already exists (use postgres default database to check)
	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname = $1)"
	err = tmpDB.Raw(query, dbName).Scan(&exists).Error

	if err != nil {
		panic("Error checking if database exists.")
	}

	// If the database exists, return early
	if exists {
		fmt.Printf("Database '%s' already exists\n", dbName)
		return
	}

	// Create database if it doesn't exist
	fmt.Printf("Database '%s' does not exist, creating it...\n", dbName)
	createQuery := fmt.Sprintf("CREATE DATABASE %s", dbName)
	err = tmpDB.Exec(createQuery).Error

	if err != nil {
		panic("Error creating database.")
	}

	fmt.Printf("Database '%s' created successfully\n", dbName)
}
