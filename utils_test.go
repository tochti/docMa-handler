package bebber

import (
  "os"
  "fmt"
  "path"
  "time"
  "strings"
  "testing"
  "crypto/sha1"
  "encoding/json"

  "gopkg.in/mgo.v2"
  "gopkg.in/mgo.v2/bson"
)


func Test_GetSettings(t *testing.T) {
  os.Setenv("TEST_ENV", "TEST_VALUE")
  if GetSettings("TEST_ENV") != "TEST_VALUE" {
    t.Error("TEST_ENV is missing!")
  }
}

func Test_SubListOK(t *testing.T) {
  a := []string{"1", "2", "3"}
  b := []string{"2", "3"}

  diff := SubList(a, b)
  if len(diff) != 1 {
    t.Fatal("Diff should be [1] but is ", diff)
  }
  if diff[0] != "1" {
    t.Error("Diff should be [1] but is ", diff)
  }

}

func Test_SubListEmpty(t *testing.T) {
  a := []string{}
  b := []string{}

  diff := SubList(a, b)
  if len(diff) != 0 {
    t.Fatal("Diff should be [] but is ", diff)
  }
}

func Test_ReadAccData(t *testing.T) {
  csvFile := path.Join(testDir, "export.csv")
  result := []AccData{}
  err := ReadAccFile(csvFile, &result)

  if err != nil {
    t.Fatal(err.Error())
  }

  d1 := time.Date(2013, time.August, 29, 0, 0, 0, 0, time.UTC)
  d2 := time.Date(2013, time.September, 01, 0, 0, 0, 0, time.UTC)
  if (result[0].Belegdatum != d1) ||
     (result[0].Buchungsdatum != d2) ||
     (result[0].Belegnummernkreis != "B") ||
     (result[0].Belegnummer != "6") ||
     (result[0].Buchungstext != "Lastschrift Strato") ||
     (result[0].Buchungsbetrag != 7.99) ||
     (result[0].Sollkonto != 71003) ||
     (result[0].Habenkonto != 1210) ||
     (result[0].Steuerschlüssel != 0) ||
     (result[0].Kostenstelle1 != "") ||
     (result[0].Kostenstelle2 != "") ||
     (result[0].BuchungsbetragEuro != 7.99) ||
     (result[0].Waehrung != "EUR") {
    t.Error("Error in CSV result ", result[0])
  }

  if len(result) != 7 {
    t.Error("Len of result should be 7, was ", len(result))
  }

}

func Test_ParseAccInt(t *testing.T) {
  r, err := ParseAccInt("")
  if err != nil {
    t.Fatal(err.Error())
  }
  if r != -1 {
    t.Error("Expect -1 was ", r)
  }

  r, err = ParseAccInt("1")
  if err != nil {
    t.Fatal(err.Error())
  }
  if r != 1 {
    t.Fatal("Expect 1 was ", r)
  }
}

func Test_ParseGermanDate(t *testing.T) {
  d := time.Date(1999, time.January, 1, 0, 0, 0, 0, time.UTC)
  result, err := ParseGermanDate("01.01.1999", ".")
  if err != nil {
    t.Fatal(err)
  }
  if d != result {
    t.Error("Expect ", d ," was ", result)
  }
}

func Test_Month(t *testing.T) {
  m, _ := Month(1)
  if m != time.January {
    t.Error("Expect ", time.January ," was ", m)
  }

  m, err := Month(13)
  if err == nil {
    t.Error("Expect to throw an error due to month is out of range")
  }
}

func Test_ZeroDate(t *testing.T) {
  z := GetZeroDate()
  if z.IsZero() != true {
    t.Error("Expect zero date was ", z)
  }
}

func Test_JoinAccFile(t *testing.T) {
  /* setup */
  // Invoices
  invo1 := AccData{
    Belegdatum: time.Date(2014,time.March,1, 0,0,0,0,time.UTC),
    Buchungsdatum: time.Date(2014, time.March,2, 0,0,0,0,time.UTC),
    Belegnummernkreis: "1",
    Belegnummer: "1",
    Sollkonto: 0,
    Habenkonto: 0,
  }
  invo2 := AccData{
    Belegdatum: time.Date(2014,time.March,1, 0,0,0,0,time.UTC),
    Buchungsdatum: time.Date(2014, time.March,2, 0,0,0,0,time.UTC),
    Belegnummernkreis: "1",
    Belegnummer: "2",
    Sollkonto: 0,
    Habenkonto: 0,
  }
  stat1 := AccData{
    Belegdatum: time.Date(2014,time.March,1, 0,0,0,0,time.UTC),
    Buchungsdatum: time.Date(2014, time.March,2, 0,0,0,0,time.UTC),
    Belegnummernkreis: "S",
    Belegnummer: "3",
    Sollkonto: 10001,
    Habenkonto: 0,
  }
  stat2 := AccData{
    Belegdatum: time.Date(2014,time.April,1, 0,0,0,0,time.UTC),
    Buchungsdatum: time.Date(2014, time.April,6, 0,0,0,0,time.UTC),
    Belegnummernkreis: "S",
    Belegnummer: "4",
    Sollkonto: 0,
    Habenkonto: 20001,
  }
  stat3 := AccData{
    Belegdatum: time.Date(2014,time.April,6, 0,0,0,0,time.UTC),
    Buchungsdatum: time.Date(2014, time.April,6, 0,0,0,0,time.UTC),
    Belegnummernkreis: "S",
    Belegnummer: "5",
    Sollkonto: 0,
    Habenkonto: 20001,
  }

  // Tmp statement to check if validCSV works !bad!
  stat4 := AccData{
    Belegdatum: time.Date(2013,time.April,6, 0,0,0,0,time.UTC),
    Buchungsdatum: time.Date(2013, time.April,6, 0,0,0,0,time.UTC),
    Belegnummernkreis: "S",
    Belegnummer: "99999",
    Sollkonto: 0,
    Habenkonto: 0,
  }

  acd := []AccData{invo1, invo2, stat1, stat2, stat3, stat4}
  // Fill database 
  session, err := mgo.Dial("127.0.0.1")
  if err != nil {
    t.Fatal(err.Error())
  }
  defer session.Close()
  err = session.DB(TestDBName).DropDatabase()
  if err != nil {
    t.Fatal(err.Error())
  }

  c := session.DB(TestDBName).C("files")

  f1 := FileDoc{
    Filename: "i1.pdf",
    ValueTags: []ValueTag{ValueTag{"Belegnummer", "11"}},
  }
  f2 := FileDoc{
    Filename: "i2.pdf",
    ValueTags: []ValueTag{ValueTag{"Belegnummer", "12"}},
  }
  f3 := FileDoc{
    Filename: "inone.pdf",
    ValueTags: []ValueTag{ValueTag{"Belegnummer", "13"}},
  }
  sD := time.Date(2014,time.February,14, 0,0,0,0,time.UTC)
  eD := time.Date(2014,time.March,1, 0,0,0,0,time.UTC)
  rT1 := RangeTag{"Belegzeitraum", sD, eD}
  f4 := FileDoc{
    Filename: "s1.pdf",
    ValueTags: []ValueTag{ValueTag{"Kontonummer", "10001"}},
    RangeTags: []RangeTag{rT1},
  }
  sD = time.Date(2014,time.April,1, 0,0,0,0,time.UTC)
  eD = time.Date(2014,time.April,18, 0,0,0,0,time.UTC)
  rT2 := RangeTag{"Belegzeitraum", sD, eD}
  f5 := FileDoc{
    Filename: "s2.pdf",
    ValueTags: []ValueTag{ValueTag{"Kontonummer", "20001"}},
    RangeTags: []RangeTag{rT2},
  }
  // Zeitraum wrong, Kontonummer right. 
  sD = time.Date(2014,time.April,20, 0,0,0,0,time.UTC)
  eD = time.Date(2014,time.April,24, 0,0,0,0,time.UTC)
  rT3 := RangeTag{"Belegzeitraum", sD, eD}
  f6 := FileDoc{
    Filename: "snone.pdf",
    ValueTags: []ValueTag{ValueTag{"Kontonummer", "10001"}},
    RangeTags: []RangeTag{rT3},
  }
  // Zeitraum right, Kontonummer wrong.
  sD = time.Date(2014,time.April,1, 0,0,0,0,time.UTC)
  eD = time.Date(2014,time.April,18, 0,0,0,0,time.UTC)
  rT4 := RangeTag{"Belegzeitraum", sD, eD}
  f7 := FileDoc{
    Filename: "snone2.pdf",
    ValueTags: []ValueTag{ValueTag{"Kontonummer", "1"}},
    RangeTags: []RangeTag{rT4},
  }

  err = c.Insert(f1, f2, f3, f4, f5, f6, f7)
  if err != nil {
    t.Fatal(err.Error())
  }

  eresult := []AccFile{
    AccFile{invo1, f1},
    AccFile{invo2, f2},
    AccFile{stat1, f4},
    AccFile{stat2, f5},
    AccFile{stat3, f5},
  }

  result, err := JoinAccFile(acd, c, false)

  if err != nil {
    t.Fatal(err.Error())
  }

  if len(result) != len(eresult) {
    t.Fatal("Expect len ", len(eresult), " was ", len(result))
    fmt.Println("Expect len ", len(eresult), " was ", len(result))
  }

  for i := range eresult {
    if eresult[i].FileDoc.Filename != result[i].FileDoc.Filename {
      t.Error("Expect ", eresult[i].FileDoc.Filename, " was ",
              result[i].FileDoc.Filename)
    }
  }

}

func Test_FileDocsMethods(t *testing.T) {
  sT := time.Date(2014, time.April, 1, 0, 0, 0, 0, time.UTC)
  eT := time.Date(2014, time.April, 2, 0, 0, 0, 0, time.UTC)
  doc1 := FileDoc{
    "test1.txt",
    []SimpleTag{SimpleTag{"sTag1"}},
    []RangeTag{
        RangeTag{"rTag1", sT, eT},
        RangeTag{"rTag2", sT, eT},
      },
    []ValueTag{ValueTag{"vTag1", "value1"}},
  }
  sT2 := time.Date(2015, time.April, 1, 0, 0, 0, 0, time.UTC)
  eT2 := time.Date(2015, time.April, 2, 0, 0, 0, 0, time.UTC)
  doc2 := FileDoc{
    "test1.txt",
    []SimpleTag{SimpleTag{"sTag1"}},
    []RangeTag{
        RangeTag{"rTag1", sT, eT},
        RangeTag{"rTag2", sT2, eT2},
      },
    []ValueTag{ValueTag{"vTag1", "value2"}},
  }
  doc3 := FileDoc{
    "test1.txt",
    []SimpleTag{SimpleTag{"sTag1"}},
    []RangeTag{},
    []ValueTag{},
  }
  doc4 := FileDoc{
    "notinlist.txt",
    []SimpleTag{},
    []RangeTag{},
    []ValueTag{},
  }

  fd := FileDocsNew([]FileDoc{doc1, doc2, doc3, doc4})

  findDoc := FileDoc{
    Filename: "test1.txt",
  }
  res := fd.FindFile(findDoc)
  if len(res.List) != 3 {
    t.Fatal("Expect 2 was ", len(res.List))
  }
  if (res.List[0].Filename != "test1.txt") ||
     (res.List[1].Filename != "test1.txt") ||
     (res.List[2].Filename != "test1.txt") {
    t.Fatal("Expect two times test1.txt was ",
             res.List[0].Filename, res.List[1].Filename)
  }

  findDoc = FileDoc{
    Filename: "test1.txt",
    RangeTags: []RangeTag{RangeTag{"rTag1", sT, eT}},
  }
  res = res.FindFile(findDoc)
  if len(res.List) != 1 {
    t.Fatal("Expect 1 was ", len(res.List))
  }
}

func Test_EmptyAccData(t *testing.T) {
  ad := AccData{}
  if ad.Empty() == false {
    t.Fatal("Expect true was false")
  }

  ad = AccData{Belegnummer: "1"}
  if ad.Empty() == true {
    t.Fatal("Expect false was true")
  }
}

func Test_ReadUser_OK(t *testing.T) {
  session, err := mgo.Dial("127.0.0.1")
  if err != nil {
    t.Fatal(err.Error())
  }
  defer session.Close()

  userTmp := User{Username: "XXX", Password: ""}
  userExpect := User{Username: "Haschel", Password: ""}
  db := session.DB(TestDBName)
  col := db.C(UsersColl)
  defer db.DropDatabase()

  err = col.Insert(userExpect, userTmp)
  if err != nil {
    t.Fatal(err.Error())
  }

  user := User{}
  err = user.Read("Haschel", db)

  if (userExpect.Username != user.Username) && (err == nil) {
    t.Fatal("Expect,", userExpect, "was,", user)
  }

}

func Test_ReadUser_Fail(t *testing.T) {
  session, err := mgo.Dial("127.0.0.1")
  if err != nil {
    t.Fatal(err.Error())
  }
  defer session.Close()

  db := session.DB(TestDBName)
  defer db.DropDatabase()

  user := User{}
  err = user.Read("Haschel", db)

  if err.Error() != "Cannot find user Haschel" {
    t.Fatal("Expect 'Cannot found user Haschel' error was", err.Error())
  }

}

func Test_SaveUser_OK(t *testing.T) {
  session, err := mgo.Dial("127.0.0.1")
  if err != nil {
    t.Fatal(err.Error())
  }
  defer session.Close()

  db := session.DB(TestDBName)
  defer db.DropDatabase()

  sha1Pass := fmt.Sprintf("%x", sha1.Sum([]byte("tt")))
  userExpect := User{Username: "test", Password: "tt"}
  err = userExpect.Save(db)

  user := User{}
  usersColl := db.C(UsersColl)
  err = usersColl.Find(bson.M{"username": "test"}).One(&user)
  if err != nil {
    t.Fatal(err.Error())
  }

  if (userExpect.Username != user.Username) && (err == nil) {
    t.Fatal("Expect", userExpect, "was", user)
  }

  if (sha1Pass != user.Password) {
    t.Fatal("Expect", sha1Pass, "was", user.Password)
  }

}

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

  accountData := DocAccountData{PostingText: "post-it"}
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

  if docUpdated.AccountData.PostingText != "post-it" {
    t.Fatal("Expect post-it was", docUpdated.AccountData.PostingText)
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

