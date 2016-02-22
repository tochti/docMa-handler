package docs

import (
	"testing"
	"time"

	"github.com/tochti/docMa-handler/common"
	"github.com/tochti/docMa-handler/labels"
	"github.com/tochti/gin-gum/gumtest"
)

func Test_SearchDocs_ByLabels(t *testing.T) {
	db := common.InitTestDB(t, AddTables, labels.AddTables)

	d := gumtest.SimpleNow()
	doc1 := Doc{
		ID:            1,
		Name:          "FindMe.pdf",
		Barcode:       "",
		DateOfScan:    d,
		DateOfReceipt: d,
	}

	label1 := labels.Label{
		ID:   1,
		Name: "l1",
	}

	DocLabel1 := DocsLabels{
		DocID:   doc1.ID,
		LabelID: label1.ID,
	}

	if err := db.Insert(&doc1,
		&label1,
		&DocLabel1); err != nil {
		t.Fatal(err)
	}

	// Add docs that we shouldn't find
	d2 := gumtest.SimpleNow()
	err := db.Insert(
		&Doc{
			Name:          "DontFindMe1.pdf",
			Barcode:       "",
			DateOfScan:    d2,
			DateOfReceipt: d2,
		},
		&Doc{
			Name:          "DontFindMe2.pdf",
			Barcode:       "",
			DateOfScan:    d,
			DateOfReceipt: d,
		},
		&labels.Label{
			ID:   3,
			Name: "catchMe",
		},
		&DocsLabels{
			DocID:   3,
			LabelID: 3,
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	searchForm := SearchForm{
		Labels: "l1,l2",
	}

	r, err := SearchDocs(db, searchForm)

	if err != nil {
		t.Fatal(err)
	}

	if len(r) != 1 {
		t.Fatal("Expect len 1 was len", len(r))
	}

	if doc1.Name != r[0].Name ||
		doc1.Barcode != r[0].Barcode ||
		!doc1.DateOfScan.Equal(r[0].DateOfScan) ||
		!doc1.DateOfReceipt.Equal(r[0].DateOfReceipt) {
		t.Fatalf("Expect %v was %v", doc1, r[0])
	}
}

func Test_SearchDocs_FromDateOfScan(t *testing.T) {
	db := common.InitTestDB(t, AddTables)

	d := gumtest.SimpleNow()
	doc1 := Doc{
		ID:            1,
		Name:          "FindMe.pdf",
		Barcode:       "",
		DateOfScan:    d,
		DateOfReceipt: d,
	}

	if err := db.Insert(&doc1); err != nil {
		t.Fatal(err)
	}

	d2 := gumtest.SimpleNow().Add(-48 * time.Hour)
	err := db.Insert(
		&Doc{
			Name:          "DontFindMe1.pdf",
			Barcode:       "",
			DateOfScan:    d2,
			DateOfReceipt: d2,
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	searchForm := SearchForm{
		DateOfScan: Interval{
			From: d,
		},
	}

	r, err := SearchDocs(db, searchForm)

	if err != nil {
		t.Fatal(err)
	}

	if len(r) != 1 {
		t.Fatal("Expect len 1 was len", len(r))
	}

	if doc1.Name != r[0].Name ||
		doc1.Barcode != r[0].Barcode ||
		!doc1.DateOfScan.Equal(r[0].DateOfScan) ||
		!doc1.DateOfReceipt.Equal(r[0].DateOfReceipt) {
		t.Fatalf("Expect %v was %v", doc1, r[0])
	}
}

func Test_SearchDocs_ToDateOfScan(t *testing.T) {
	db := common.InitTestDB(t, AddTables)

	d := gumtest.SimpleNow()
	doc1 := Doc{
		ID:            1,
		Name:          "FindMe.pdf",
		Barcode:       "",
		DateOfScan:    d,
		DateOfReceipt: d,
	}

	if err := db.Insert(&doc1); err != nil {
		t.Fatal(err)
	}

	d2 := gumtest.SimpleNow().Add(+48 * time.Hour)
	err := db.Insert(
		&Doc{
			Name:          "DontFindMe1.pdf",
			Barcode:       "",
			DateOfScan:    d2,
			DateOfReceipt: d2,
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	searchForm := SearchForm{
		DateOfScan: Interval{
			To: d.Add(+24 * time.Hour),
		},
	}

	r, err := SearchDocs(db, searchForm)

	if err != nil {
		t.Fatal(err)
	}

	if len(r) != 1 {
		t.Fatal("Expect len 1 was len", len(r))
	}

	if doc1.Name != r[0].Name ||
		doc1.Barcode != r[0].Barcode ||
		!doc1.DateOfScan.Equal(r[0].DateOfScan) ||
		!doc1.DateOfReceipt.Equal(r[0].DateOfReceipt) {
		t.Fatalf("Expect %v was %v", doc1, r[0])
	}
}

func Test_SearchDocs_BetweenDates(t *testing.T) {
	db := common.InitTestDB(t, AddTables)

	d := gumtest.SimpleNow()
	doc1 := Doc{
		ID:            1,
		Name:          "FindMe.pdf",
		Barcode:       "",
		DateOfScan:    d,
		DateOfReceipt: d,
	}

	if err := db.Insert(&doc1); err != nil {
		t.Fatal(err)
	}

	d2 := gumtest.SimpleNow().Add(-48 * time.Hour)
	err := db.Insert(
		&Doc{
			Name:          "DontFindMe1.pdf",
			Barcode:       "",
			DateOfScan:    d2,
			DateOfReceipt: d2,
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	searchForm := SearchForm{
		DateOfScan: Interval{
			From: d.Add(-24 * time.Hour),
			To:   d.Add(+24 * time.Hour),
		},
	}

	r, err := SearchDocs(db, searchForm)

	if err != nil {
		t.Fatal(err)
	}

	if len(r) != 1 {
		t.Fatal("Expect len 1 was len", len(r))
	}

	if doc1.Name != r[0].Name ||
		doc1.Barcode != r[0].Barcode ||
		!doc1.DateOfScan.Equal(r[0].DateOfScan) ||
		!doc1.DateOfReceipt.Equal(r[0].DateOfReceipt) {
		t.Fatalf("Expect %v was %v", doc1, r[0])
	}
}

func Test_SearchDocs_ByDocNumbers(t *testing.T) {
	db := common.InitTestDB(t, AddTables, labels.AddTables)

	d := gumtest.SimpleNow()
	doc1 := Doc{
		ID:            1,
		Name:          "FindMe.pdf",
		Barcode:       "",
		DateOfScan:    d,
		DateOfReceipt: d,
	}

	docNumber1 := DocNumber{
		DocID:  doc1.ID,
		Number: "1",
	}

	if err := db.Insert(&doc1, &docNumber1); err != nil {
		t.Fatal(err)
	}

	// Add docs that we shouldn't find
	d2 := gumtest.SimpleNow()
	err := db.Insert(
		&Doc{
			Name:          "DontFindMe1.pdf",
			Barcode:       "",
			DateOfScan:    d2,
			DateOfReceipt: d2,
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	searchForm := SearchForm{
		Labels:     "",
		DocNumbers: "1,2",
		DateOfScan: Interval{
			From: d,
			To:   d,
		},
	}

	r, err := SearchDocs(db, searchForm)

	if err != nil {
		t.Fatal(err)
	}

	if len(r) != 1 {
		t.Fatal("Expect len 1 was len", len(r))
	}

	if doc1.Name != r[0].Name ||
		doc1.Barcode != r[0].Barcode ||
		!doc1.DateOfScan.Equal(r[0].DateOfScan) ||
		!doc1.DateOfReceipt.Equal(r[0].DateOfReceipt) {
		t.Fatalf("Expect %v was %v", doc1, r[0])
	}
}
func Test_SearchDocs(t *testing.T) {
	db := common.InitTestDB(t, AddTables, labels.AddTables)

	d := gumtest.SimpleNow()
	doc1 := Doc{
		ID:            1,
		Name:          "FindMe.pdf",
		Barcode:       "",
		DateOfScan:    d,
		DateOfReceipt: d,
	}

	docNumber1 := DocNumber{
		DocID:  doc1.ID,
		Number: "1",
	}

	docNumber2 := DocNumber{
		DocID:  doc1.ID,
		Number: "2",
	}

	label1 := labels.Label{
		ID:   1,
		Name: "l1",
	}

	label2 := labels.Label{
		ID:   2,
		Name: "l2",
	}

	DocLabel1 := DocsLabels{
		DocID:   doc1.ID,
		LabelID: label1.ID,
	}

	DocLabel2 := DocsLabels{
		DocID:   doc1.ID,
		LabelID: label2.ID,
	}

	if err := db.Insert(&doc1,
		&docNumber1,
		&docNumber2,
		&label1,
		&label2,
		&DocLabel1,
		&DocLabel2); err != nil {
		t.Fatal(err)
	}

	// Add docs that we shouldn't find
	d2 := gumtest.SimpleNow()
	err := db.Insert(
		&Doc{
			Name:          "DontFindMe1.pdf",
			Barcode:       "",
			DateOfScan:    d2,
			DateOfReceipt: d2,
		},
		&Doc{
			Name:          "DontFindMe2.pdf",
			Barcode:       "",
			DateOfScan:    d,
			DateOfReceipt: d,
		},
		&Doc{
			ID:            4,
			Name:          "DontFindMe3.pdf",
			Barcode:       "",
			DateOfScan:    d,
			DateOfReceipt: d,
		},
		&DocsLabels{
			DocID:   4,
			LabelID: label1.ID,
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	searchForm := SearchForm{
		Labels:     "l1,l2",
		DocNumbers: "1,2",
		DateOfScan: Interval{
			From: d,
			To:   d,
		},
	}

	r, err := SearchDocs(db, searchForm)

	if err != nil {
		t.Fatal(err)
	}

	if len(r) != 1 {
		t.Fatal("Expect len 1 was len", len(r))
	}

	if doc1.Name != r[0].Name ||
		doc1.Barcode != r[0].Barcode ||
		!doc1.DateOfScan.Equal(r[0].DateOfScan) ||
		!doc1.DateOfReceipt.Equal(r[0].DateOfReceipt) {
		t.Fatalf("Expect %v was %v", doc1, r[0])
	}
}

func Test_ParseQueryString_None(t *testing.T) {
	l := parseQueryString("")

	expect := []string{}
	if !compareStrings(expect, l) {
		t.Fatalf("Expect %v was %v", expect, l)
	}
}

func Test_ParseQueryString_OneLabel(t *testing.T) {
	l := parseQueryString("l1")

	expect := []string{"l1"}
	if !compareStrings(expect, l) {
		t.Fatalf("Expect %v was %v", expect, l)
	}
}

func Test_ParseQueryString_OneInterestLabel(t *testing.T) {
	l := parseQueryString("l1, ")

	expect := []string{"l1"}
	if !compareStrings(expect, l) {
		t.Fatalf("Expect %v was %v", expect, l)
	}
}

func Test_ParseQueryString_Spaces(t *testing.T) {
	l := parseQueryString(" li rum , la rum, \t   loeffel stihl ")

	expect := []string{"li rum", "la rum", "loeffel stihl"}
	if !compareStrings(expect, l) {
		t.Fatalf("Expect %v was %v", expect, l)
	}
}

func Test_ParseQueryString(t *testing.T) {
	l := parseQueryString("l1,l2")

	expect := []string{"l1", "l2"}
	if !compareStrings(expect, l) {
		t.Fatalf("Expect %v was %v", expect, l)
	}
}

func compareStrings(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	for i, e := range a {
		if e != b[i] {
			return false
		}
	}

	return true

}

/*
func Test_SearchDocs_WrongQueryFail(t *testing.T) {
	g := MakeTestGlobals(t)
	session := g.MongoDB.Session.Copy()
	defer session.Close()
	db := session.DB(g.Config["MONGODB_DBNAME"])
	defer db.DropDatabase()

	searchStr := `{
    "Labels": "l1,l2"
    "DocNumbers": "1,2",
    "DateOfScan": {
      "From": "2000-01-01T00:00:00Z",
      "To": "2000-01-01T00:00:00Z"
    }
  }`
	result, err := SearchDocs(searchStr, db)

	if len(result) != 0 {
		t.Fatal("Expect result of zero was", result)
	}

	errMsg := "invalid character '\"' after object key:value pair"
	if strings.Contains(err.Error(), errMsg) == false {
		t.Fatal("Expect", errMsg, "was", err.Error())
	}
}

func Test_SearchDocs_EmptyQueryFail(t *testing.T) {
	g := MakeTestGlobals(t)
	session := g.MongoDB.Session.Copy()
	defer session.Close()
	db := session.DB(g.Config["MONGODB_DBNAME"])
	defer db.DropDatabase()

	searchStr := ``
	result, err := SearchDocs(searchStr, db)

	if len(result) != 0 {
		t.Fatal("Expect result of zero was", result)
	}

	errMsg := "unexpected end of JSON input"
	if strings.Contains(err.Error(), errMsg) == false {
		t.Fatal("Expect", errMsg, "was", err.Error())
	}
}

func Test_SearchDocs_EmptyJsonQueryFail(t *testing.T) {
	g := MakeTestGlobals(t)
	session := g.MongoDB.Session.Copy()
	defer session.Close()
	db := session.DB(g.Config["MONGODB_DBNAME"])
	defer db.DropDatabase()

	searchStr := `{}`
	result, err := SearchDocs(searchStr, db)

	if len(result) != 0 {
		t.Fatal("Expect result of zero was", result)
	}

	errMsg := "Empty search query"
	if strings.Contains(err.Error(), errMsg) == false {
		t.Fatal("Expect", errMsg, "was", err.Error())
	}
}

func Test_SearchDocs_WrongValueTypeFail(t *testing.T) {
	g := MakeTestGlobals(t)
	session := g.MongoDB.Session.Copy()
	defer session.Close()
	db := session.DB(g.Config["MONGODB_DBNAME"])
	defer db.DropDatabase()

	searchStr := `{
    "Labels":["Neu"],
    "DocNumbers": ["15"],
    "DateOfScan": []}`
	result, err := SearchDocs(searchStr, db)

	if len(result) != 0 {
		t.Fatal("Expect result of zero was", result)
	}

	errMsg := "Empty search query"
	if strings.Contains(err.Error(), errMsg) == false {
		t.Fatal("Expect", errMsg, "was", err.Error())
	}
}

func Test_SearchDocs_WrongValueType2Fail(t *testing.T) {
	g := MakeTestGlobals(t)
	session := g.MongoDB.Session.Copy()
	defer session.Close()
	db := session.DB(g.Config["MONGODB_DBNAME"])
	defer db.DropDatabase()

	searchStr := `{
    "Labels":["Neu"],
    "DocNumbers": ["15"],
    "DateOfScan": {
      "From": [],
      "To": []
    }
  }`
	result, err := SearchDocs(searchStr, db)

	if len(result) != 0 {
		t.Fatal("Expect result of zero was", result)
	}

	errMsg := "Empty search query"
	if strings.Contains(err.Error(), errMsg) == false {
		t.Fatal("Expect", errMsg, "was", err.Error())
	}
}
*/
