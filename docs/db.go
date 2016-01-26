package docs

import "gopkg.in/gorp.v1"

func AddTables(db *gorp.DbMap) {
	tMap := db.AddTableWithName(Doc{}, DocsTable).
		SetKeys(true, "id")
	tMap.ColMap("name").SetUnique(true).SetNotNull(true)

	db.AddTableWithName(DocAccountData{}, DocAccountDataTable).
		SetKeys(false, "doc_id")

	db.AddTableWithName(DocNumber{}, DocNumbersTable).
		SetKeys(false, "doc_id", "number")
}
