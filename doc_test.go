package bebber

import (
  "time"
  "strings"
  "testing"
  "encoding/json"

  "gopkg.in/mgo.v2/bson"
)

func Test_FindDoc_OK(t *testing.T) {
  globals := MakeTestGlobals(t)
  session := globals.MongoDB.Session.Copy()
  defer session.Close()

  db := session.DB(TestDBName)
  defer db.DropDatabase()

  docID := bson.NewObjectId()
  expectDoc := Doc{ID: docID, Name: "karl.pdf", Note: "Note", Labels: []Label{}}
  err := db.C(DocsColl).Insert(expectDoc)
  if err != nil {
    t.Fatal(err.Error())
  }

  doc := Doc{Name: "karl.pdf"}
  err = doc.Find(db)
  if err != nil {
    t.Fatal(err.Error())
  }

  expectDocJSON, err := json.Marshal(expectDoc)
  if err != nil {
    t.Fatal(err.Error())
  }

  docJSON, err := json.Marshal(doc)
  if err != nil {
    t.Fatal(err.Error())
  }

  if string(expectDocJSON) != string(docJSON) {
    t.Fatal("Exepect", string(expectDocJSON), "was", string(docJSON))
  }

}

func Test_FindDoc_Fail(t *testing.T) {
  globals := MakeTestGlobals(t)
  session := globals.MongoDB.Session.Copy()
  defer session.Close()

  db := session.DB(TestDBName)
  defer db.DropDatabase()

  doc := Doc{Name: "karl.pdf"}
  err = doc.Find(db)
  if strings.Contains(err.Error(), "Cannot find") == false {
    t.Fatal("Expect to fail with Cannot find document was", err)
  }
}

func Test_ChangeDoc_OK(t *testing.T) {
  globals := MakeTestGlobals(t)
  session := globals.MongoDB.Session.Copy()
  defer session.Close()

  db := session.DB(TestDBName)
  defer db.DropDatabase()

  doc := Doc{Name: "changeme", Infos: DocInfos{}, Note: "note"}

  docsColl := db.C(DocsColl)
  err := docsColl.Insert(doc)
  if err != nil {
    t.Fatal(err.Error())
  }

  accountData := DocAccountData{DocNumber: "123"}
  labels := []Label{"label1"}
  changeDoc := Doc{Name: "nicer", Barcode: "barcode",
                   AccountData: accountData, Labels: labels}
  err = doc.Change(changeDoc, db)
  if err != nil {
    t.Fatal(err.Error())
  }

  failDoc := Doc{Name: "changeme"}
  err = failDoc.Find(db)
  if err == nil {
    t.Fatal("Expect Cannot find document error was nil")
  }
  if strings.Contains(err.Error(), "Cannot find document") == false{
    t.Fatal("Expect Cannot find document error was", err.Error())
  }

  docUpdated := Doc{Name: "nicer"}
  err = docUpdated.Find(db)
  if err != nil {
    t.Fatal(err.Error())
  }

  if doc.Name != "nicer" {
    t.Fatal("Expect nicer was", doc.Name)
  }

  if doc.Barcode != "barcode" {
    t.Fatal("Expect barcode was", doc.Barcode)
  }

  if docUpdated.Barcode != "barcode" {
    t.Fatal("Expect barcode was", docUpdated.Barcode)
  }

  if docUpdated.Note != "note" {
    t.Fatal("Expect note was", docUpdated.Note)
  }

  if docUpdated.AccountData.DocNumber != "123" {
    t.Fatal("Expect 123 was", docUpdated.AccountData.DocNumber)
  }

  if docUpdated.Labels[0] != "label1" {
    t.Fatal("Expect label1 was", docUpdated.Labels[0])
  }
}

func Test_ChangeDoc_InfoFail(t *testing.T) {
  globals := MakeTestGlobals(t)
  session := globals.MongoDB.Session.Copy()
  defer session.Close()

  db := session.DB(TestDBName)
  defer db.DropDatabase()

  doc := Doc{Name: "changeme", Infos: DocInfos{}}

  docsColl := db.C(DocsColl)
  err := docsColl.Insert(doc)
  if err != nil {
    t.Fatal(err.Error())
  }

  docInfos := DocInfos{
              DateOfScan: time.Now(),
            }
  changeDoc := Doc{Name: "changeme", Infos: docInfos}
  err = doc.Change(changeDoc, db)
  if err == nil {
    t.Fatal("Expect Not allowed to change infos error was nil")
  }

  if strings.Contains(err.Error(), "Not allowed to change infos") == false {
    t.Fatal(err.Error())
  }

}

func Test_RemoveDoc_OK(t *testing.T) {
  globals := MakeTestGlobals(t)
  session := globals.MongoDB.Session.Copy()
  defer session.Close()

  db := session.DB(TestDBName)
  defer db.DropDatabase()

  doc := Doc{Name: "Pimpel", Infos: DocInfos{}}

  docsColl := db.C(DocsColl)
  err := docsColl.Insert(doc)
  if err != nil {
    t.Fatal(err.Error())
  }

  err = doc.Remove(db)

  if err != nil {
    t.Fatal(err.Error())
  }

  err = doc.Find(db)
  if strings.Contains(err.Error(), "Cannot find document") == false {
    t.Fatal(err.Error())
  }

}

func Test_AppendLabels_OK(t *testing.T) {
  globals := MakeTestGlobals(t)
  session := globals.MongoDB.Session.Copy()
  defer session.Close()

  db := session.DB(TestDBName)
  defer db.DropDatabase()

  docTmp := Doc{
              Name: "Hoocker",
              Labels: []Label{"label1"},
            }
  docsColl := db.C(DocsColl)
  err := docsColl.Insert(docTmp)
  if err != nil {
    t.Fatal(err.Error())
  }

  labels := []Label{Label("l2"), Label("l3")}
  docTmp.AppendLabels(labels, db)

  doc := Doc{Name: "Hoocker"}
  err = doc.Find(db)
  if err != nil {
    t.Fatal(err)
  }

  if len(doc.Labels) != 3 {
    t.Fatal("Expect 3 labels was", doc.Labels)
  }
  if len(docTmp.Labels) != 3 {
    t.Fatal("Expect 3 labels was", doc.Labels)
  }
}

func Test_RemoveLabels_OK(t *testing.T) {
  globals := MakeTestGlobals(t)
  session := globals.MongoDB.Session.Copy()
  defer session.Close()

  db := session.DB(TestDBName)
  defer db.DropDatabase()

  docTmp := Doc{
              Name: "Hoocker",
              Labels: []Label{"l1", "l2", "l3"},
            }
  docsColl := db.C(DocsColl)
  err := docsColl.Insert(docTmp)
  if err != nil {
    t.Fatal(err.Error())
  }

  labels := []Label{"l2", "l3"}
  docTmp.RemoveLabels(labels, db)

  doc := Doc{Name: "Hoocker"}
  err = doc.Find(db)
  if err != nil {
    t.Fatal(err)
  }

  if len(doc.Labels) != 1 {
    t.Fatal("Expect 3 labels was", doc.Labels)
  }
  if len(docTmp.Labels) != 1 {
    t.Fatal("Expect 3 labels was", doc.Labels)
  }

  if docTmp.Labels[0] != "l1" {
    t.Fatal("Expect l1 was", docTmp.Labels[0])
  }
}

func Test_DocAccountDataIsValid_OK(t *testing.T) {
  docAccountData := DocAccountData{
    DocNumber: "docnumber",
  }

  ok, err := docAccountData.IsValid()
  if ok != true {
    t.Fatal(err.Error())
  }

  docAccountData = DocAccountData{
    DocPeriod: DocPeriod{
      From: time.Now(),
      To: time.Now(),
    },
    AccNumber: 1993321,
  }

  ok, err = docAccountData.IsValid()
  if ok != true {
    t.Fatal(err.Error())
  }
}

func Test_DocAccountDataIsValid_Fail(t *testing.T) {
  docAccountData := DocAccountData{
    DocNumber: "docnumber",
    AccNumber: 19999,
  }

  ok, err := docAccountData.IsValid()
  errMsg := "Accountant data mismatch!"
  if ok != false {
    t.Fatal("Expect", errMsg, "was", err)
  }

  docAccountData = DocAccountData{
    DocNumber: "docnumber",
    DocPeriod: DocPeriod{
      From: time.Now(),
      To: time.Now(),
    },
  }

  ok, err = docAccountData.IsValid()
  errMsg = "Accountant data mismatch!"
  if ok != false {
    t.Fatal("Expect", errMsg, "was", err)
  }

  docAccountData = DocAccountData{
  }

  ok, err = docAccountData.IsValid()
  errMsg = "Missing Accountant data!"
  if ok != false {
    t.Fatal("Exepct", errMsg, "was", err)
  }

  docAccountData = DocAccountData{
    AccNumber: 123456,
  }

  ok, err = docAccountData.IsValid()
  errMsg = "Missing document period data!"
  if ok != false {
    t.Fatal("Exepct", errMsg, "was", err)
  }

  docAccountData = DocAccountData{
    DocPeriod: DocPeriod{
      From: time.Now(),
      To: time.Now(),
    },
  }

  ok, err = docAccountData.IsValid()
  errMsg = "Missing account number!"
  if ok != false {
    t.Fatal("Exepct", errMsg, "was", err)
  }
}

