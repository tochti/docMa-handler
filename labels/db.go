package labels

import "gopkg.in/gorp.v1"

var (
	LabelsTable = "labels"
)

func AddTables(db *gorp.DbMap) *gorp.TableMap {
	return db.AddTableWithName(Label{}, LabelsTable).SetKeys(true, "id")
}

func CreateTables(db *gorp.DbMap) error {
	AddTables(db)
	err := db.CreateTablesIfNotExists()
	if err != nil {
		return err
	}

	return nil
}
