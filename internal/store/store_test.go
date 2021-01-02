package store

import (
	"github.com/labstack/gommon/log"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"testing"
)

const datasourceUrl = "./accountinformation.sqlite"

func init() {
	_ = os.Setenv("DRIVER_NAME", "sqlite3")
	_ = os.Setenv("DATASOURCE_URL", datasourceUrl)
}

func Test_dbConn(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"test_db_conn_success"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := LoadDBConnection()
			if db == nil {
				t.Errorf("dbConn() = %v", db)
			}

			_ = db.Close()
		})
	}

	err := os.Remove(datasourceUrl)
	if err != nil {
		log.Warn(err)
	}
}
