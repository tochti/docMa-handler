package docs

import (
	"reflect"
	"testing"

	"github.com/tochti/docMa-handler/common"
	"github.com/tochti/docMa-handler/labels"
	"github.com/tochti/gin-gum/gumtest"
)

func Test_FindLabelsOfDoc(t *testing.T) {
	db := common.InitTestDB(t, AddTables, labels.AddTables)

	label := labels.Label{
		ID:   1,
		Name: "label",
	}

	doc := Doc{
		ID:   1,
		Name: "karl.pdf",
	}

	docsLabels := DocsLabels{
		DocID:   doc.ID,
		LabelID: label.ID,
	}

	if err := db.Insert(&label); err != nil {
		t.Fatal(err)
	}

	if err := db.Insert(&doc); err != nil {
		t.Fatal(err)
	}

	if err := db.Insert(&docsLabels); err != nil {
		t.Fatal(err)
	}

	r, err := FindLabelsOfDoc(db, doc.ID)
	if err != nil {
		t.Fatal(err)
	}

	expect := []labels.Label{label}
	ok := reflect.DeepEqual(expect, r)
	if !ok {
		t.Fatalf("Expect %v was %v", expect, r)
	}

}

func Test_ReadDocNumbers(t *testing.T) {
	db := initDB(t)

	docNumber := DocNumber{
		DocID:  1,
		Number: "1",
	}

	if err := db.Insert(&docNumber); err != nil {
		t.Fatal(err)
	}

	docNumbers, err := ReadDocNumbers(db, 1)
	if err != nil {
		t.Fatal(err)
	}

	if len(docNumbers) != 1 {
		t.Fatalf("Expect %v was %v", 1, len(docNumbers))
	}

	if docNumbers[0].Number != docNumber.Number {
		t.Fatalf("Expect %v was %v", docNumber.Number, docNumbers[0].Number)
	}

}

func Test_ReadAccountData(t *testing.T) {
	db := initDB(t)

	accountData := DocAccountData{
		DocID:         1,
		AccountNumber: 12,
		PeriodFrom:    gumtest.SimpleNow(),
		PeriodTo:      gumtest.SimpleNow(),
	}

	if err := db.Insert(&accountData); err != nil {
		t.Fatal(err)
	}

	r, err := ReadAccountData(db, 1)
	if err != nil {
		t.Fatal(err)
	}

	if r.DocID != accountData.DocID ||
		r.AccountNumber != accountData.AccountNumber ||
		!r.PeriodFrom.Equal(accountData.PeriodFrom) ||
		!r.PeriodTo.Equal(accountData.PeriodTo) {
		t.Fatalf("Expect %v was %v", accountData, r)
	}

}
