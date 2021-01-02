package configmanager

import (
	"github.com/kaanaktas/openbanking-accountinformation/api"
	"github.com/kaanaktas/openbanking-accountinformation/internal/store"
	"log"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	if os.Getenv("DRIVER_NAME") == "" {
		_ = os.Setenv("DRIVER_NAME", "sqlite3")

	}
	if os.Getenv("DATASOURCE_URL") == "" {
		_ = os.Setenv("DATASOURCE_URL", "../../testdata/accountinformation.sqlite")
	}

	log.Print("CONFIG-MANAGER START")
	dbx := store.LoadDBConnection()

	api.RunSql(dbx, "../../testdata/insert_data.down.sql")
	api.RunSql(dbx, "../../testdata/insert_data.up.sql")

	exitCode := m.Run()
	log.Print("CONFIG-MANAGER END")
	os.Exit(exitCode)
}
