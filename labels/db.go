package labels

import "gopkg.in/gorp.v1"

var (
	LabelsTable = "labels"
)

func AddTables(db *gorp.DbMap) {
	db.AddTableWithName(Label{}, LabelsTable).
		SetKeys(true, "id").
		ColMap("name").
		SetUnique(true)
}

func CreateTables(db *gorp.DbMap) error {
	AddTables(db)
	err := db.CreateTablesIfNotExists()
	if err != nil {
		return err
	}

	return nil
}
