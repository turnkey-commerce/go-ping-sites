package database

import "testing"

// TestCreateDb tests the creation of the database.
func TestCreateDb(t *testing.T) {
	db, err := InitializeDB()
	if err != nil {
		t.Fatal("Failed to create database:", err)
	}

	errPing := db.Ping()
	if errPing != nil {
		t.Fatal("Failed to ping database:", errPing)
	}
}
