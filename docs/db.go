package docs

import "gopkg.in/gorp.v1"

func AddTables(db *gorp.DbMap) {
	tMap := db.AddTableWithName(Doc{}, DocsTable).
		SetKeys(true, "id")
	tMap.ColMap("name").SetUnique(true).SetNotNull(true)

	db.AddTableWithName(DocAccountData{}, DocAccountDataTable).
		SetKeys(true, "id")

	db.AddTableWithName(DocNumber{}, DocNumbersTable).SetKeys(true, "id")
}
