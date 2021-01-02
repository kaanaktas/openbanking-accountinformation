package main

import (
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
)

func init() {
	_ = godotenv.Load()
}

func main() {
	doMigrate()
}

func doMigrate() {
	m, err := migrate.New(
		os.Getenv("MIGRATE_SCRIPT_URL"),
		os.Getenv("MIGRATE_DATABASE_URL"))
	defer m.Close()

	if err != nil {
		log.Fatalf("could not migrate %v", err)
	}

	log.Println(m.Version())
	version, _ := strconv.Atoi(os.Getenv("MIGRATE_VERSION"))
	err = m.Migrate(uint(version))
	if err != nil {
		log.Fatalf("err in migration, %v", err)
	}

	log.Println("migration has been completed")
	log.Println(m.Version())
}
