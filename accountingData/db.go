package accountingData

import "gopkg.in/gorp.v1"

func AddTables(db *gorp.DbMap) {
	db.AddTableWithName(AccountingData{}, AccountingDataTable).SetKeys(true, "id")
}
