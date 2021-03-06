package docs

import (
	"github.com/tochti/docMa-handler/labels"
	"gopkg.in/gorp.v1"
)

func AddTables(db *gorp.DbMap) {
	tMap := db.AddTableWithName(Doc{}, DocsTable).
		SetKeys(true, "id")
	tMap.ColMap("name").SetUnique(true).SetNotNull(true)
	tMap.ColMap("barcode").SetUnique(true)

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

func ReadDocNumbers(db *gorp.DbMap, docID int64) ([]DocNumber, error) {
	docNumbers := []DocNumber{}
	_, err := db.Select(
		&docNumbers,
		Q("SELECT * FROM %v WHERE doc_id=?", DocNumbersTable),
		docID,
	)
	if err != nil {
		return []DocNumber{}, err
	}

	return docNumbers, nil
}

func ReadAccountData(db *gorp.DbMap, docID int64) (DocAccountData, error) {
	accountData := DocAccountData{}
	err := db.SelectOne(
		&accountData,
		Q("SELECT * FROM %v WHERE doc_id=?", DocAccountDataTable),
		docID,
	)
	if err != nil {
		return DocAccountData{}, err
	}
	return accountData, nil
}

func FindDocsWithLabel(db *gorp.DbMap, labelID int64) ([]Doc, error) {
	d := []Doc{}

	q := Q(`
	SELECT 
		docs.id,
		docs.name,
		docs.barcode,
		docs.date_of_scan,
		docs.date_of_receipt,
		docs.note
	FROM %v as docs, %v as docs_labels
	WHERE docs_labels.label_id=?
	AND docs.id=docs_labels.doc_id`, DocsTable, DocsLabelsTable)
	_, err := db.Select(&d, q, labelID)
	if err != nil {
		return []Doc{}, err
	}

	return d, nil
}

// Remove all doc labels for one doc
func RemoveDocLabelConnection(db *gorp.DbMap, docID int64) error {
	q := Q("DELETE FROM %v WHERE doc_id=?", DocsLabelsTable)

	_, err := db.Exec(q, docID)
	if err != nil {
		return err
	}

	return nil
}

// Remove doc account data for one doc
func RemoveAccountData(db *gorp.DbMap, docID int64) error {
	q := Q("DELETE FROM %v WHERE doc_id=?", DocAccountDataTable)

	_, err := db.Exec(q, docID)
	if err != nil {
		return err
	}

	return nil
}

// Remove all doc numbers for one doc
func RemoveDocNumbers(db *gorp.DbMap, docID int64) error {
	q := Q("DELETE FROM %v WHERE doc_id=?", DocNumbersTable)

	_, err := db.Exec(q, docID)
	if err != nil {
		return err
	}

	return nil
}
