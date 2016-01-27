package docs

import (
	"github.com/tochti/docMa-handler/labels"
	"gopkg.in/gorp.v1"
)

func AddTables(db *gorp.DbMap) {
	tMap := db.AddTableWithName(Doc{}, DocsTable).
		SetKeys(true, "id")
	tMap.ColMap("name").SetUnique(true).SetNotNull(true)

	db.AddTableWithName(DocAccountData{}, DocAccountDataTable).
		SetKeys(false, "doc_id")

	db.AddTableWithName(DocNumber{}, DocNumbersTable).
		SetKeys(false, "doc_id", "number")

	db.AddTableWithName(DocsLabels{}, DocsLabelsTable).
		SetKeys(false, "doc_id", "label_id")
}

func FindLabelsOfDoc(db *gorp.DbMap, docID int64) ([]labels.Label, error) {

	q := Q(`
	SELECT labels.id,labels.name 
	FROM %v as labels, %v as labels_docs
	WHERE labels_docs.doc_id=?
	AND labels.id=labels_docs.label_id`,
		labels.LabelsTable, DocsLabelsTable)

	labelList := []labels.Label{}
	if _, err := db.Select(&labelList, q, docID); err != nil {
		return []labels.Label{}, err
	}

	return labelList, nil

}
