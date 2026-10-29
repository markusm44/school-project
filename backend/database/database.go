package driver

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/mattes/migrate/source/file"
	_ "github.com/mattn/go-sqlite3"
)

var DBase *sql.DB

func ConnDB() {
	db, err := sql.Open("sqlite3", "./database/forum.db?_foreign_key=on")
	if err != nil {
		log.Fatal(err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatal("Cannot ping database!")
	}
	err = migrateDB(db, "./database/migrations/")
	if err != nil {
		log.Fatalf("Error applying database migrations: %d", err)
	}
	DBase = db
	log.Println("DBase conn established")
}

func migrateDB(db *sql.DB, migrationsDir string) error {
	driver, err := sqlite3.WithInstance(db, &sqlite3.Config{}) 
	if err != nil {
		return fmt.Errorf("failed to initialize sqlite3 driver: %d", err)
	}

	m, err := migrate.NewWithDatabaseInstance( 
		"file://"+migrationsDir,
		"sqlite3", driver)
	if err != nil {
		return fmt.Errorf("failed to initialize migrate instance: %d", err)
	}

	err = m.Up() 

	if err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to apply database migrations: %d", err)
	}

	return nil
}

func CloseDB() {
	DBase.Close()
	log.Println("Database connection closed")
}
