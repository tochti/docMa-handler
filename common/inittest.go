package common

import (
	"os"
	"testing"

	"gopkg.in/gorp.v1"
)

var (
	TestDBName = "testing"
)

func setenvTest() {
	os.Clearenv()

	os.Setenv("MYSQL_USER", "tochti")
	os.Setenv("MYSQL_PASSWORD", "123")
	os.Setenv("MYSQL_HOST", "127.0.0.1")
	os.Setenv("MYSQL_PORT", "3306")
	os.Setenv("MYSQL_DB_NAME", TestDBName)
}

func InitTestDB(t *testing.T, AddTablesFunc func(*gorp.DbMap)) *gorp.DbMap {
	setenvTest()

	db := InitMySQL()
	AddTablesFunc(db)

	err := db.DropTablesIfExists()
	if err != nil {
		t.Fatal(err)
	}
	err = db.CreateTablesIfNotExists()
	if err != nil {
		t.Fatal(err)
	}

	return db
}
