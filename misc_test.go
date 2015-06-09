package bebber

import (
  "os"
  "fmt"
  "path"
  "time"
  "testing"

  "gopkg.in/mgo.v2"
)


func TestGetSettings(t *testing.T) {
  os.Setenv("TEST_ENV", "TEST_VALUE")
  if GetSettings("TEST_ENV") != "TEST_VALUE" {
    t.Error("TEST_ENV is missing!")
  }
}

func TestSubListOK(t *testing.T) {
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

func TestSubListEmpty(t *testing.T) {
  a := []string{}
  b := []string{}

  diff := SubList(a, b)
  if len(diff) != 0 {
    t.Fatal("Diff should be [] but is ", diff)
  }
}

func TestReadAccData(t *testing.T) {
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

func TestParseAccInt(t *testing.T) {
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

func TestParseGermanDate(t *testing.T) {
  d := time.Date(1999, time.January, 1, 0, 0, 0, 0, time.UTC)
  result, err := ParseGermanDate("01.01.1999", ".")
  if err != nil {
    t.Fatal(err)
  }
  if d != result {
    t.Error("Expect ", d ," was ", result)
  }
}

func TestMonth(t *testing.T) {
  m, _ := Month(1)
  if m != time.January {
    t.Error("Expect ", time.January ," was ", m)
  }

  m, err := Month(13)
  if err == nil {
    t.Error("Expect to throw an error due to month is out of range")
  }
}

func TestZeroDate(t *testing.T) {
  z := GetZeroDate()
  if z.IsZero() != true {
    t.Error("Expect zero date was ", z)
  }
}

func TestSpotTagType(t *testing.T) {
  ty, err := SpotTagType("test")
  if ty != "SimpleTag" {
    t.Error("Should be SimpleTag is ", ty)
  }
  if err != nil {
    t.Error("Error should be empty is ", err.Error())
  }

  ty, err = SpotTagType("test:")
  if ty != "" {
    t.Error("Should be an wrong tag is ", ty)
  }
  if err == nil {
    t.Error("Error should be (Missing value) but is nil")
  }

  ty, err = SpotTagType("test:1234")
  if ty != "ValueTag" {
    t.Error("Should be a ValueTag is ", ty)
  }
  if err != nil {
    t.Error("Error should be empty is ", err.Error())
  }

  ty, err = SpotTagType("test:\"hallo hallo\"")
  if ty != "ValueTag" {
    t.Error("Should be a ValueTag is ", ty)
  }
  if err != nil {
    t.Error("Error should be empty is ", err.Error())
  }

  ty, err = SpotTagType("test:er:li")
  if ty != "ValueTag" {
    t.Error("Should be a ValueTag is ", ty)
  }
  if err != nil {
    t.Error("Error should be empty is ", err.Error())
  }

  ty, err = SpotTagType("test:01012014..02022014")
  if ty != "RangeTag" {
    t.Error("Should be a RangeTag is ", ty)
  }
  if err != nil {
    t.Error("Error should be empty is ", err.Error())
  }

  ty, err = SpotTagType("test:1102014..02022104")
  if ty != "RangeTag" {
    t.Error("Should be RangeTag is ", ty)
  }
  if err.Error() != "Error in range" {
    t.Error("Error msg should be (Error in range) is ", err.Error())
  }

  ty, err = SpotTagType("test:1102014..02022104")
  if ty != "RangeTag" {
    t.Error("Should be RangeTag is ", ty)
  }
  if err.Error() != "Error in range" {
    t.Error("Error msg should be (Error in range) is ", err.Error())
  }

  ty, err = SpotTagType("test:1102014..2022104")
  if ty != "RangeTag" {
    t.Error("Should be RangeTag is ", ty)
  }
  if err.Error() != "Error in range" {
    t.Error("Error msg should be (Error in range) is ", err.Error())
  }

  ty, err = SpotTagType("test:..02022015")
  if ty != "RangeTag" {
    t.Error("Should be RangeTag is ", ty)
  }
  if err != nil {
    t.Error("No error should occur ", err.Error())
  }

  ty, err = SpotTagType("test:02022015..")
  if ty != "RangeTag" {
    t.Error("Should be RangeTag is ", ty)
  }
  if err != nil {
    t.Error("No error should occur ", err.Error())
  }

}

func TestCreateUpdateDocSimpleTag(t *testing.T) {
  doc := FileDoc{Filename: "test.txt"}
  err := CreateUpdateDoc([]string{}, &doc)
  if err != nil {
    t.Error(err.Error())
  }
  if doc.Filename != "test.txt" {
    t.Error("#1 wrong filename (", doc.Filename, ")")
  }
  if len(doc.SimpleTags) != 0 {
    t.Error("expect [] is ", doc.SimpleTags)
  }

  doc = FileDoc{
          Filename: "test.txt",
          SimpleTags: []SimpleTag{SimpleTag{"sTag"}},
        }
  err = CreateUpdateDoc([]string{"sTag1", "sTag2"}, &doc)
  if err != nil {
    t.Error(err.Error())
  }
  if doc.Filename != "test.txt" {
    t.Error("#2 wrong filename (", doc.Filename, ")")
  }
  if len(doc.SimpleTags) != 3 {
    t.Fatal("Expect 3 SimpleTags got ", len(doc.SimpleTags))
  }
  if doc.SimpleTags[0].Tag != "sTag" || doc.SimpleTags[1].Tag != "sTag1" || doc.SimpleTags[2].Tag != "sTag2" {
    t.Error("wrong tags ", doc.SimpleTags)
  }
}

func TestCreateUpdateValueTag(t *testing.T) {
  tags := []string{"vTag1:1234", "vTag2:\"foo bar\"", "vTag3:va:lue"}
  doc := FileDoc{Filename: "test.txt"}
  err := CreateUpdateDoc(tags, &doc)
  if err != nil {
    t.Error(err.Error())
  }
  if doc.Filename != "test.txt" {
    t.Error("wrong filename (", doc.Filename, ")")
  }

  if doc.ValueTags[0].Tag != "vTag1" || doc.ValueTags[0].Value != "1234" {
    t.Error("wrong value tag #1 ", doc.ValueTags[0])
  }
  if doc.ValueTags[1].Tag != "vTag2" || doc.ValueTags[1].Value != "foo bar" {
    t.Error("wrong value tag #2 ", doc.ValueTags[1])
  }
  if doc.ValueTags[2].Tag != "vTag3" || doc.ValueTags[2].Value != "va:lue" {
    t.Error("wrong value tag #3 ", doc.ValueTags[2])
  }
}

func TestCreateUpdateRangeTag(t *testing.T) {
  tags := []string{"rT1:01042014..02042014", "rt2:..02042014", "rt3:01042014.."}
  doc := FileDoc{Filename: "test.txt"}
  err := CreateUpdateDoc(tags, &doc)
  if err != nil {
    t.Fatal(err.Error())
  }
  if doc.Filename != "test.txt" {
    t.Error("wrong filename (", doc.Filename, ")")
  }

  sD := time.Date(2014, time.April, 1, 0, 0, 0, 0, time.UTC)
  eD := time.Date(2014, time.April, 2, 0, 0, 0, 0, time.UTC)
  if doc.RangeTags[0].Tag != "rT1" || doc.RangeTags[0].Start != sD || doc.RangeTags[0].End != eD {
    t.Error("wrong range tag ", doc.RangeTags[0])
  }
}

func TestJoinAccFile(t *testing.T) {
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
  err = session.DB("bebber_test").DropDatabase()
  if err != nil {
    t.Fatal(err.Error())
  }

  c := session.DB("bebber_test").C("files")

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

func TestFileDocsMethods(t *testing.T) {
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

func TestEmptyAccData(t *testing.T) {
  ad := AccData{}
  if ad.Empty() == false {
    t.Fatal("Expect true was false")
  }

  ad = AccData{Belegnummer: "1"}
  if ad.Empty() == true {
    t.Fatal("Expect false was true")
  }
}

func TestLoadUserOk(t *testing.T) {
  session, err := mgo.Dial("127.0.0.1")
  if err != nil {
    t.Fatal(err.Error())
  }
  defer session.Close()

  userTmp := User{Username: "XXX", Password: "", Dirs: map[string]string{"i":"ih"}}
  userExpect := User{Username: "Haschel", Password: "", Dirs: map[string]string{"i":"ah"}}
  db := session.DB("bebber_test")
  col := db.C(UsersCollection)
  defer db.DropDatabase()

  err = col.Insert(userExpect, userTmp)
  if err != nil {
    t.Fatal(err.Error())
  }

  user := User{}
  err = user.Load("Haschel", col)

  if (userExpect.Username != user.Username) && (err == nil) {
    t.Fatal("Expect,", userExpect, "was,", user)
  }

}

func TestLoadUserFail(t *testing.T) {
  session, err := mgo.Dial("127.0.0.1")
  if err != nil {
    t.Fatal(err.Error())
  }
  defer session.Close()
  col := session.DB("bebber_test").C(UsersCollection)
  user := User{}
  err = user.Load("Haschel", col)

  if err.Error() != "Cannot find user Haschel" {
    t.Fatal("Expect 'Cannot found user Haschel' error was", err.Error())
  }

}
