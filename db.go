package bebber

import (
	"log"

	_ "github.com/go-sql-driver/mysql"

	"github.com/tochti/docMa-handler/labels"
	"github.com/tochti/gin-gum/gumspecs"
	"gopkg.in/gorp.v1"
)

func InitMySQL() *gorp.DbMap {

	mysql := gumspecs.ReadMySQL()

	sqlDB, err := mysql.DB()
	if err != nil {
		log.Fatal(err)
	}

	db := &gorp.DbMap{
		Db: sqlDB,
		Dialect: gorp.MySQLDialect{
			"InnonDB",
			"UTF8",
		},
	}

	labels.AddTables(db)

	return db
}
