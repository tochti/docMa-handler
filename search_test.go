package bebber

import (
	"strings"
	"testing"
	"time"
)

func Test_SearchDocs_OK(t *testing.T) {
	g := MakeTestGlobals(t)
	session := g.MongoDB.Session.Copy()
	defer session.Close()
	db := session.DB(g.Config["MONGODB_DBNAME"])
	defer db.DropDatabase()

	expectDateOfScan := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
	doc1 := Doc{
		Name:   "Tasty",
		Labels: []Label{Label("l1"), Label("l2")},
		AccountData: DocAccountData{
			DocNumbers: []string{"1", "2"},
		},
		Infos: DocInfos{
			DateOfScan: expectDateOfScan,
		},
	}

	tmpDate := time.Date(1999, 1, 1, 0, 0, 0, 0, time.UTC)
	doc2 := Doc{
		Labels: []Label{Label("l1"), Label("l2")},
		AccountData: DocAccountData{
			DocNumbers: []string{"1", "2"},
		},
		Infos: DocInfos{
			DateOfScan: tmpDate,
		},
	}

	docsColl := db.C(DocsColl)
	err = docsColl.Insert(doc1, doc2)
	if err != nil {
		t.Fatal(err.Error())
	}

	searchStr := `{
    "Labels": "l1,l2",
    "DocNumbers": "1,2",
    "DateOfScan": {
      "From": "2000-01-01T00:00:00Z",
      "To": "2000-01-01T00:00:00Z"
    }
  }`
	result, err := SearchDocs(searchStr, db)

	if err != nil {
		t.Fatal(err.Error())
	}

	if len(result) != 1 {
		t.Fatal("Expect len 1 was", len(result))
	}

	if doc1.Name != result[0].Name {
		t.Fatal("Expect", doc1.Name, "was", result[0].Name)
	}
}

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
