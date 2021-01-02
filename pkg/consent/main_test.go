package consent

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

	_ = os.Setenv("CLIENT_CA_CERT_PEM", "../../certs/ob_issuer.cer,../../certs/danske_sandbox.cer,../../certs/ozone_sandbox.cer")
	_ = os.Setenv("CLIENT_CERT_PEM", "../../certs/ob_transport.pem")
	_ = os.Setenv("CLIENT_KEY_PEM", "../../certs/ob_transport.key")
	_ = os.Setenv("OB_SIGN_KEY", "../../certs/ob_signing.key")

	log.Print("CONSENT START")
	dbx := store.LoadDBConnection()

	api.RunSql(dbx, "../../testdata/insert_data.down.sql")
	api.RunSql(dbx, "../../testdata/insert_data.up.sql")

	exitCode := m.Run()
	log.Print("CONSENT END")
	os.Exit(exitCode)
}
