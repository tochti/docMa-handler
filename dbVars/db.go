package dbVars

import "gopkg.in/gorp.v1"

func AddTables(db *gorp.DbMap) {
	db.AddTableWithName(DBVar{}, DBVarsTable).SetKeys(false, "name")
}
