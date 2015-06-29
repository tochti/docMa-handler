package bebber

import (
  "path"
  "time"
  "strings"
  "testing"
)

func Test_ReadAccProcessFile_OK(t *testing.T) {
  csvFile := path.Join(testDir, "export.csv")
  result, err := ReadAccProcessFile(csvFile)

  if err != nil {
    t.Fatal(err.Error())
  }

  d1 := time.Date(2013, time.August, 29, 0, 0, 0, 0, time.UTC)
  d2 := time.Date(2013, time.September, 01, 0, 0, 0, 0, time.UTC)
  if (result[0].DocDate != d1) ||
     (result[0].DateOfEntry != d2) ||
     (result[0].DocNumberRange != "B") ||
     (result[0].DocNumber != "6") ||
     (result[0].PostingText != "Lastschrift Strato") ||
     (result[0].AmountPosted != 7.99) ||
     (result[0].DebitAcc != 71003) ||
     (result[0].CreditAcc != 1210) ||
     (result[0].TaxCode != 0) ||
     (result[0].CostUnit1 != "") ||
     (result[0].CostUnit2 != "") ||
     (result[0].AmountPostedEuro != 7.99) ||
     (result[0].Currency != "EUR") {
    t.Error("Error in CSV result ", result[0])
  }

  if len(result) != 7 {
    t.Error("Len of result should be 7, was ", len(result))
  }

}

func Test_ParseAccInt_OK(t *testing.T) {
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

/*
func Test_JoinAccFile_OK(t *testing.T) {
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
*/

func Test_FindAccProcessByDocNumber_OK(t *testing.T) {
  globals := MakeTestGlobals(t)
  session := globals.MongoDB.Session.Copy()
  defer session.Close()

  db := session.DB(TestDBName)
  defer db.DropDatabase()

  ImportAccProcess(db, path.Join(testDir, "export.csv"))

  docNumbers := []string{"B6", "13"}
  result, err := FindAccProcessByDocNumbers(db, docNumbers)

  if err != nil {
    t.Fatal(err.Error())
  }

  if len(result) != 2 {
    t.Fatal("Expect 2 results was", result)
  }

  if result[0].DocNumberRange != "B" {
    t.Fatal("Expect B was", result[0].DocNumberRange)
  }
  if result[0].DocNumber != "6" {
    t.Fatal("Expect 6 was", result[0].DocNumber)
  }

  if result[1].DocNumberRange != "" {
    t.Fatal("Expect empty string was", result[0].DocNumberRange)
  }
  if result[1].DocNumber != "13" {
    t.Fatal("Expect 13 was", result[0].DocNumber)
  }


}

func Test_FindAccProcessByAccNumber_OK(t *testing.T) {
  globals := MakeTestGlobals(t)
  session := globals.MongoDB.Session.Copy()
  defer session.Close()

  db := session.DB(TestDBName)
  defer db.DropDatabase()

  ImportAccProcess(db, path.Join(testDir, "export.csv"))

  accNumber := 71002
  fromDate := time.Date(2013,8,21,0,0,0,0,time.UTC)
  toDate := time.Date(2013,9,1,0,0,0,0,time.UTC)
  result, err := FindAccProcessByAccNumber(db, accNumber, fromDate, toDate)

  if err != nil {
    t.Fatal(err.Error())
  }

  if len(result) != 2 {
    t.Fatal("Expect 2 results was", result)
  }

  if result[0].DocNumber != "9" {
    t.Fatal("Expect 9 was", result[0].DocNumber)
  }
  if result[1].DocNumber != "13" {
    t.Fatal("Expect 13 was", result[1].DocNumber)
  }


}

func Test_FindDocByDocNumber_OK(t *testing.T) {

}

func Test_FindDocByAccNumber_OK(t *testing.T) {

}

func Test_SplitDocNumber_OK(t *testing.T) {
  rang, number, err := SplitDocNumber("BB123")
  if err != nil {
    t.Fatal("Expect nil was", err)
  }
  if rang != "BB" {
    t.Fatal("Expect BB was", rang)
  }
  if number != "123" {
    t.Fatal("Expect 123 was", number)
  }

  rang, number, err = SplitDocNumber("987")
  if err != nil {
    t.Fatal("Expect nil was", err)
  }
  if rang != "" {
    t.Fatal("Expect empty string was", rang)
  }
  if number != "987" {
    t.Fatal("Expect 987 was", number)
  }
}

func Test_SplitDocNumber_Fail(t *testing.T) {
  _, _, err := SplitDocNumber("BB")

  expectErrMsg := "Invalid docnumber!"
  if err != nil {
    if strings.Contains(err.Error(), expectErrMsg) == false {
      t.Fatal("Expect", expectErrMsg, "was nil")
    }
  } else {
    t.Fatal("Expect", expectErrMsg, "was nil")
  }
}
