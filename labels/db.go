package labels

import (
	"database/sql"

	"gopkg.in/gorp.v1"
)

var (
	LabelsTable = "labels"
)

func InitGorp(sqlDB *sql.DB) *gorp.DbMap {
	db := &gorp.DbMap{
		Db: sqlDB,
		Dialect: gorp.MySQLDialect{
			"InnonDB",
			"UTF8",
		},
	}

	db.AddTableWithName(Label{}, LabelsTable).SetKeys(true, "id")

	return db
}
