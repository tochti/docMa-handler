package docs

import (
	"reflect"
	"testing"

	"github.com/tochti/docMa-handler/common"
	"github.com/tochti/docMa-handler/labels"
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
