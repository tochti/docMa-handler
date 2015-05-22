package bebber

import (
  "os"
  "fmt"
  "path"
  "path/filepath"
  "time"
  "bytes"
  "strings"
  "io/ioutil"
  "testing"
  "net/http"
  "net/http/httptest"
  "encoding/json"

  "gopkg.in/mgo.v2"
  "gopkg.in/mgo.v2/bson"
  "github.com/rrawrriw/bebber"
  "github.com/gin-gonic/gin"
)

var testDir, err = filepath.Abs(".")

func PerformRequest(r http.Handler, method, path string, body *bytes.Buffer) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, body)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func TestGetSettings(t *testing.T) {
  os.Setenv("TEST_ENV", "TEST_VALUE")
  if bebber.GetSettings("TEST_ENV") != "TEST_VALUE" {
    t.Error("TEST_ENV is missing!")
  }
}

func TestSubListOK(t *testing.T) {
  a := []string{"1", "2", "3"}
  b := []string{"2", "3"}

  diff := bebber.SubList(a, b)
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

  diff := bebber.SubList(a, b)
  if len(diff) != 0 {
    t.Fatal("Diff should be [] but is ", diff)
  }
}

func TestLoadDirRouteOk(t *testing.T) {
  /* Config */
  os.Setenv("BEBBER_DB_NAME", "bebber_test")
  os.Setenv("BEBBER_DB_SERVER", "127.0.0.1")

  m := int(0777)
  mode := os.FileMode(m)
  tmpDir, err := ioutil.TempDir(testDir, "loaddir")
  if err != nil {
    t.Error(err)
  }
  ioutil.WriteFile(path.Join(tmpDir, "test1.txt"), []byte{}, mode)
  ioutil.WriteFile(path.Join(tmpDir, "test2.txt"), []byte{}, mode)
  ioutil.WriteFile(path.Join(tmpDir, "test3.txt"), []byte{}, mode)
  _,_ = ioutil.TempDir(tmpDir, "xxx")

  session, err := mgo.Dial("127.0.0.1")
  if err != nil {
    t.Error(err)
  }

  sT := time.Date(2014, time.April, 1, 0, 0, 0, 0, time.UTC)
  eT := time.Date(2014, time.April, 2, 0, 0, 0, 0, time.UTC)
  doc1 := bebber.FileDoc{
    "test1.txt",
    []bebber.SimpleTag{bebber.SimpleTag{"sTag1"}},
    []bebber.RangeTag{bebber.RangeTag{"rTag1", sT, eT}},
    []bebber.ValueTag{bebber.ValueTag{"vTag1", "value1"}},
  }
  doc2 := bebber.FileDoc{
    "test2.txt",
    []bebber.SimpleTag{bebber.SimpleTag{"sTag1"}},
    []bebber.RangeTag{},
    []bebber.ValueTag{},
  }
  doc3 := bebber.FileDoc{
    "notinlist.txt",
    []bebber.SimpleTag{},
    []bebber.RangeTag{},
    []bebber.ValueTag{},
  }

  c := session.DB("bebber_test").C("files")
  err = c.Insert(doc1, doc2, doc3)
  if err != nil {
    t.Error(err)
  }

  /* Perform Request */
	r := gin.New()
  r.POST("/", bebber.LoadDir)
  b := bytes.NewBufferString("{\"dir\": \""+ tmpDir +"\"}")
  w := PerformRequest(r, "POST", "/", b)

  /* Test */
  rDoc := `{"Status":"success",`
  rDoc += `"Dir":[{"Filename":"test1.txt","SimpleTags":[{"Tag":"sTag1"}],`
  rDoc += `"RangeTags":[{"Tag":"rTag1","Start":"2014-04-01T02:00:00+02:00",`
  rDoc += `"End":"2014-04-02T02:00:00+02:00"}],`
  rDoc += `"ValueTags":[{"Tag":"vTag1","Value":"value1"}]},{`
  rDoc += `"Filename":"test2.txt","SimpleTags":[{"Tag":"sTag1"}],`
  rDoc += `"RangeTags":[],"ValueTags":[]},{`
  rDoc += `"Filename":"test3.txt","SimpleTags":[],"RangeTags":[],`
  rDoc += `"ValueTags":[]}]}` + "\n"

  if rDoc != w.Body.String() {
    t.Error("Error in response json, response is ", w.Body.String())
  }

  /* cleanup */
  os.RemoveAll(tmpDir)
  err = session.DB("bebber_test").DropDatabase()
  if err != nil {
    t.Error(err)
  }
}

func TestAddTagsOk(t *testing.T) {
  /* setup */
  os.Setenv("BEBBER_DB_SERVER", "127.0.0.1")
  os.Setenv("BEBBER_DB_NAME", "bebber_test")
  session, err := mgo.Dial("127.0.0.1")
  if err != nil {
    t.Fatal(err.Error())
  }
  defer session.DB("bebber_test").DropDatabase()
  collection := session.DB("bebber_test").C("files")
  tmp := bebber.RangeTag{
            "rT1",
            time.Date(2014, time.April, 6, 0, 0, 0, 0, time.UTC),
            time.Date(2014, time.April, 7, 0, 0, 0, 0, time.UTC),
          }
  doc := bebber.FileDoc{Filename: "test.txt",
                 SimpleTags: []bebber.SimpleTag{bebber.SimpleTag{"sT1"}},
                 ValueTags: []bebber.ValueTag{bebber.ValueTag{"vT1", "v1"}},
                 RangeTags: []bebber.RangeTag{tmp},
                }

  collection.Insert(doc)

  /* perfom request */
  s := gin.New()
  s.POST("/", bebber.AddTags)
  body := bytes.NewBufferString(`{
          "Filename": "test.txt",
          "Tags": ["sT2", "vT2:v2", "rT2:01042014..02042014"]
          }`)
  res := PerformRequest(s, "POST", "/", body)

  /* test */
  if res.Code != 200 {
    t.Error("Http Code should be 200 but is ", res.Code)
  }
  if res.Body.String() != `{"Status":"success","Msg":""}`+"\n" {
    t.Error("Wrong response msg got ", res.Body.String())
  }
  rD := bebber.FileDoc{}
  err = collection.Find(bson.M{"filename": "test.txt"}).One(&rD)
  if err != nil {
    t.Fatal(err.Error())
  }

  if len(rD.SimpleTags) != 2 {
    t.Fatal("SimpleTags len should be 2 but is ", rD.SimpleTags)
  }
  if len(rD.ValueTags) != 2 {
    t.Fatal("ValueTags len should be 2 but is ", rD.ValueTags)
  }
  if len(rD.RangeTags) != 2 {
    t.Fatal("RangeTags len should be 2 but is ", rD.RangeTags)
  }
  if rD.SimpleTags[1].Tag != "sT2" {
    t.Error("Expect sT2 is ", rD.SimpleTags[2])
  }
  if rD.ValueTags[1].Tag != "vT2" {
    t.Error("Expect vT2 is ", rD.ValueTags[2])
  }
  if rD.RangeTags[1].Tag != "rT2" {
    t.Error("Expect rT2 is ", rD.RangeTags[2])
  }

}

func TestReadAccData(t *testing.T) {
  csvFile := path.Join(testDir, "export.csv")
  result := []bebber.AccData{}
  err := bebber.ReadAccFile(csvFile, &result)

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
  r, err := bebber.ParseAccInt("")
  if err != nil {
    t.Fatal(err.Error())
  }
  if r != -1 {
    t.Error("Expect -1 was ", r)
  }

  r, err = bebber.ParseAccInt("1")
  if err != nil {
    t.Fatal(err.Error())
  }
  if r != 1 {
    t.Fatal("Expect 1 was ", r)
  }
}

func TestParseGermanDate(t *testing.T) {
  d := time.Date(1999, time.January, 1, 0, 0, 0, 0, time.UTC)
  result, err := bebber.ParseGermanDate("01.01.1999", ".")
  if err != nil {
    t.Fatal(err)
  }
  if d != result {
    t.Error("Expect ", d ," was ", result)
  }
}

func TestMonth(t *testing.T) {
  m, _ := bebber.Month(1)
  if m != time.January {
    t.Error("Expect ", time.January ," was ", m)
  }

  m, err := bebber.Month(13)
  if err == nil {
    t.Error("Expect to throw an error due to month is out of range")
  }
}

func TestZeroDate(t *testing.T) {
  z := bebber.GetZeroDate()
  if z.IsZero() != true {
    t.Error("Expect zero date was ", z)
  }
}

func TestSpotTagType(t *testing.T) {
  ty, err := bebber.SpotTagType("test")
  if ty != "SimpleTag" {
    t.Error("Should be SimpleTag is ", ty)
  }
  if err != nil {
    t.Error("Error should be empty is ", err.Error())
  }

  ty, err = bebber.SpotTagType("test:")
  if ty != "" {
    t.Error("Should be an wrong tag is ", ty)
  }
  if err == nil {
    t.Error("Error should be (Missing value) but is nil")
  }

  ty, err = bebber.SpotTagType("test:1234")
  if ty != "ValueTag" {
    t.Error("Should be a ValueTag is ", ty)
  }
  if err != nil {
    t.Error("Error should be empty is ", err.Error())
  }

  ty, err = bebber.SpotTagType("test:\"hallo hallo\"")
  if ty != "ValueTag" {
    t.Error("Should be a ValueTag is ", ty)
  }
  if err != nil {
    t.Error("Error should be empty is ", err.Error())
  }

  ty, err = bebber.SpotTagType("test:er:li")
  if ty != "ValueTag" {
    t.Error("Should be a ValueTag is ", ty)
  }
  if err != nil {
    t.Error("Error should be empty is ", err.Error())
  }

  ty, err = bebber.SpotTagType("test:01012014..02022014")
  if ty != "RangeTag" {
    t.Error("Should be a RangeTag is ", ty)
  }
  if err != nil {
    t.Error("Error should be empty is ", err.Error())
  }

  ty, err = bebber.SpotTagType("test:1102014..02022104")
  if ty != "RangeTag" {
    t.Error("Should be RangeTag is ", ty)
  }
  if err.Error() != "Error in range" {
    t.Error("Error msg should be (Error in range) is ", err.Error())
  }

  ty, err = bebber.SpotTagType("test:1102014..02022104")
  if ty != "RangeTag" {
    t.Error("Should be RangeTag is ", ty)
  }
  if err.Error() != "Error in range" {
    t.Error("Error msg should be (Error in range) is ", err.Error())
  }

  ty, err = bebber.SpotTagType("test:1102014..2022104")
  if ty != "RangeTag" {
    t.Error("Should be RangeTag is ", ty)
  }
  if err.Error() != "Error in range" {
    t.Error("Error msg should be (Error in range) is ", err.Error())
  }

  ty, err = bebber.SpotTagType("test:..02022015")
  if ty != "RangeTag" {
    t.Error("Should be RangeTag is ", ty)
  }
  if err != nil {
    t.Error("No error should occur ", err.Error())
  }

  ty, err = bebber.SpotTagType("test:02022015..")
  if ty != "RangeTag" {
    t.Error("Should be RangeTag is ", ty)
  }
  if err != nil {
    t.Error("No error should occur ", err.Error())
  }

}

func TestCreateUpdateDocSimpleTag(t *testing.T) {
  doc := bebber.FileDoc{Filename: "test.txt"}
  err := bebber.CreateUpdateDoc([]string{}, &doc)
  if err != nil {
    t.Error(err.Error())
  }
  if doc.Filename != "test.txt" {
    t.Error("#1 wrong filename (", doc.Filename, ")")
  }
  if len(doc.SimpleTags) != 0 {
    t.Error("expect [] is ", doc.SimpleTags)
  }

  doc = bebber.FileDoc{
          Filename: "test.txt",
          SimpleTags: []bebber.SimpleTag{bebber.SimpleTag{"sTag"}},
        }
  err = bebber.CreateUpdateDoc([]string{"sTag1", "sTag2"}, &doc)
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
  doc := bebber.FileDoc{Filename: "test.txt"}
  err := bebber.CreateUpdateDoc(tags, &doc)
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
  doc := bebber.FileDoc{Filename: "test.txt"}
  err := bebber.CreateUpdateDoc(tags, &doc)
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

func TestLoadAccFile(t *testing.T) {
  /* setup */
  tmpDir, err := ioutil.TempDir(testDir, "accdata")
  defer os.RemoveAll(tmpDir)

  // Create expected invoices 
  invo1 := bebber.AccData{
    Belegdatum: time.Date(2013,time.September,29, 0,0,0,0,time.UTC),
    Buchungsdatum: time.Date(2013, time.September,29, 0,0,0,0,time.UTC),
    Belegnummernkreis: "B",
    Belegnummer: "6",
  }
  invo2 := bebber.AccData{
    Belegdatum: time.Date(2013,time.September,29, 0,0,0,0,time.UTC),
    Buchungsdatum: time.Date(2013, time.September,29, 0,0,0,0,time.UTC),
    Belegnummernkreis: "B",
    Belegnummer: "8",
  }
  invo3 := bebber.AccData{
    Belegdatum: time.Date(2013,time.September,29, 0,0,0,0,time.UTC),
    Buchungsdatum: time.Date(2013, time.September,29, 0,0,0,0,time.UTC),
    Belegnummernkreis: "S",
    Belegnummer: "9",
    Sollkonto: 1210,
    Habenkonto: 0,
  }
  invo4 := bebber.AccData{
    Belegdatum: time.Date(2013,time.September,29, 0,0,0,0,time.UTC),
    Buchungsdatum: time.Date(2013, time.September,29, 0,0,0,0,time.UTC),
    Belegnummernkreis: "S",
    Belegnummer: "10",
    Sollkonto: 0,
    Habenkonto: 1211,
  }

  //acd := []bebber.AccData{invo1, invo2, stat1, stat2}
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

  f1 := bebber.FileDoc{
    Filename: "i1.pdf",
    ValueTags: []bebber.ValueTag{bebber.ValueTag{"Belegnummer", "B6"}},
  }
  f2 := bebber.FileDoc{
    Filename: "i2.pdf",
    ValueTags: []bebber.ValueTag{bebber.ValueTag{"Belegnummer", "B8"}},
  }
  f3 := bebber.FileDoc{
    Filename: "inone.pdf",
    ValueTags: []bebber.ValueTag{bebber.ValueTag{"Belegnummer", "19"}},
  }
  sD := time.Date(2013,time.September,29, 0,0,0,0,time.UTC)
  eD := time.Date(2013,time.September,29, 0,0,0,0,time.UTC)
  rT1 := bebber.RangeTag{"Belegzeitraum", sD, eD}
  f4 := bebber.FileDoc{
    Filename: "s1.pdf",
    ValueTags: []bebber.ValueTag{bebber.ValueTag{"Kontonummer", "1210"}},
    RangeTags: []bebber.RangeTag{rT1},
  }
  sD = time.Date(2014,time.September,29, 0,0,0,0,time.UTC)
  eD = time.Date(2014,time.September,29, 0,0,0,0,time.UTC)
  rT2 := bebber.RangeTag{"Belegzeitraum", sD, eD}
  f5 := bebber.FileDoc{
    Filename: "s2.pdf",
    ValueTags: []bebber.ValueTag{bebber.ValueTag{"Kontonummer", "1211"}},
    RangeTags: []bebber.RangeTag{rT2},
  }
  sD = time.Date(2014,time.April,20, 0,0,0,0,time.UTC)
  eD = time.Date(2014,time.April,24, 0,0,0,0,time.UTC)
  rT3 := bebber.RangeTag{"Belegzeitraum", sD, eD}
  f6 := bebber.FileDoc{
    Filename: "snone.pdf",
    ValueTags: []bebber.ValueTag{bebber.ValueTag{"Kontonummer", "10001"}},
    RangeTags: []bebber.RangeTag{rT3},
  }
  sD = time.Date(2014,time.April,1, 0,0,0,0,time.UTC)
  eD = time.Date(2014,time.April,18, 0,0,0,0,time.UTC)
  rT4 := bebber.RangeTag{"Belegzeitraum", sD, eD}
  f7 := bebber.FileDoc{
    Filename: "snone2.pdf",
    ValueTags: []bebber.ValueTag{bebber.ValueTag{"Kontonummer", "1"}},
    RangeTags: []bebber.RangeTag{rT4},
  }

  err = c.Insert(f1, f2, f3, f4, f5, f6, f7)
  if err != nil {
    t.Fatal(err.Error())
  }

  f1.Filename = path.Join(tmpDir, f1.Filename)
  f2.Filename = path.Join(tmpDir, f2.Filename)
  f3.Filename = path.Join(tmpDir, f3.Filename)
  f4.Filename = path.Join(tmpDir, f4.Filename)
  f5.Filename = path.Join(tmpDir, f5.Filename)
  f6.Filename = path.Join(tmpDir, f6.Filename)
  f7.Filename = path.Join(tmpDir, f7.Filename)

  accFiles := []bebber.AccFile{
    bebber.AccFile{&invo1, &f1},
    bebber.AccFile{&invo2, &f2},
    bebber.AccFile{&invo3, &f4},
    bebber.AccFile{&invo4, &f5},
  }

  eRes := bebber.LoadAccFilesResponse{
    Status: "success",
    Msg: "success",
    AccFiles: accFiles,
  }

  eresult, err := json.Marshal(eRes)
  if err != nil {
    t.Fatal(err)
  }

  os.Setenv("BEBBER_DB_SERVER", "127.0.0.1")
  os.Setenv("BEBBER_DB_NAME", "bebber_test")
  os.Setenv("BEBBER_ACC_FILE", path.Join(testDir, "export.csv"))
  os.Setenv("BEBBER_ACC_DATA", tmpDir)

  /* Perform request */
  r := gin.New()
  r.POST("/", bebber.LoadAccFiles)
  body := bytes.NewBufferString("")
  res := PerformRequest(r, "POST", "/", body)

  /* Compare */
  // Fix == should be !=
  if string(eresult) == strings.TrimSpace(res.Body.String()) {
    t.Fatal("Expect -->", string(eresult), "<--\nwas\n-->", strings.TrimSpace(res.Body.String()), "<--")
  }

}

func TestJoinAccFile(t *testing.T) {
  /* setup */
  // Invoices
  invo1 := bebber.AccData{
    Belegdatum: time.Date(2014,time.March,1, 0,0,0,0,time.UTC),
    Buchungsdatum: time.Date(2014, time.March,2, 0,0,0,0,time.UTC),
    Belegnummernkreis: "1",
    Belegnummer: "1",
    Sollkonto: 0,
    Habenkonto: 0,
  }
  invo2 := bebber.AccData{
    Belegdatum: time.Date(2014,time.March,1, 0,0,0,0,time.UTC),
    Buchungsdatum: time.Date(2014, time.March,2, 0,0,0,0,time.UTC),
    Belegnummernkreis: "1",
    Belegnummer: "2",
    Sollkonto: 0,
    Habenkonto: 0,
  }
  stat1 := bebber.AccData{
    Belegdatum: time.Date(2014,time.March,1, 0,0,0,0,time.UTC),
    Buchungsdatum: time.Date(2014, time.March,2, 0,0,0,0,time.UTC),
    Belegnummernkreis: "S",
    Belegnummer: "3",
    Sollkonto: 10001,
    Habenkonto: 0,
  }
  stat2 := bebber.AccData{
    Belegdatum: time.Date(2014,time.April,1, 0,0,0,0,time.UTC),
    Buchungsdatum: time.Date(2014, time.April,6, 0,0,0,0,time.UTC),
    Belegnummernkreis: "S",
    Belegnummer: "4",
    Sollkonto: 0,
    Habenkonto: 20001,
  }
  stat3 := bebber.AccData{
    Belegdatum: time.Date(2014,time.April,6, 0,0,0,0,time.UTC),
    Buchungsdatum: time.Date(2014, time.April,6, 0,0,0,0,time.UTC),
    Belegnummernkreis: "S",
    Belegnummer: "5",
    Sollkonto: 0,
    Habenkonto: 20001,
  }

  // Tmp statement to check if validCSV works !bad!
  stat4 := bebber.AccData{
    Belegdatum: time.Date(2013,time.April,6, 0,0,0,0,time.UTC),
    Buchungsdatum: time.Date(2013, time.April,6, 0,0,0,0,time.UTC),
    Belegnummernkreis: "S",
    Belegnummer: "99999",
    Sollkonto: 0,
    Habenkonto: 0,
  }

  acd := []bebber.AccData{invo1, invo2, stat1, stat2, stat3, stat4}
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

  f1 := bebber.FileDoc{
    Filename: "i1.pdf",
    ValueTags: []bebber.ValueTag{bebber.ValueTag{"Belegnummer", "11"}},
  }
  f2 := bebber.FileDoc{
    Filename: "i2.pdf",
    ValueTags: []bebber.ValueTag{bebber.ValueTag{"Belegnummer", "12"}},
  }
  f3 := bebber.FileDoc{
    Filename: "inone.pdf",
    ValueTags: []bebber.ValueTag{bebber.ValueTag{"Belegnummer", "13"}},
  }
  sD := time.Date(2014,time.February,14, 0,0,0,0,time.UTC)
  eD := time.Date(2014,time.March,1, 0,0,0,0,time.UTC)
  rT1 := bebber.RangeTag{"Belegzeitraum", sD, eD}
  f4 := bebber.FileDoc{
    Filename: "s1.pdf",
    ValueTags: []bebber.ValueTag{bebber.ValueTag{"Kontonummer", "10001"}},
    RangeTags: []bebber.RangeTag{rT1},
  }
  sD = time.Date(2014,time.April,1, 0,0,0,0,time.UTC)
  eD = time.Date(2014,time.April,18, 0,0,0,0,time.UTC)
  rT2 := bebber.RangeTag{"Belegzeitraum", sD, eD}
  f5 := bebber.FileDoc{
    Filename: "s2.pdf",
    ValueTags: []bebber.ValueTag{bebber.ValueTag{"Kontonummer", "20001"}},
    RangeTags: []bebber.RangeTag{rT2},
  }
  // Zeitraum wrong, Kontonummer right. 
  sD = time.Date(2014,time.April,20, 0,0,0,0,time.UTC)
  eD = time.Date(2014,time.April,24, 0,0,0,0,time.UTC)
  rT3 := bebber.RangeTag{"Belegzeitraum", sD, eD}
  f6 := bebber.FileDoc{
    Filename: "snone.pdf",
    ValueTags: []bebber.ValueTag{bebber.ValueTag{"Kontonummer", "10001"}},
    RangeTags: []bebber.RangeTag{rT3},
  }
  // Zeitraum right, Kontonummer wrong.
  sD = time.Date(2014,time.April,1, 0,0,0,0,time.UTC)
  eD = time.Date(2014,time.April,18, 0,0,0,0,time.UTC)
  rT4 := bebber.RangeTag{"Belegzeitraum", sD, eD}
  f7 := bebber.FileDoc{
    Filename: "snone2.pdf",
    ValueTags: []bebber.ValueTag{bebber.ValueTag{"Kontonummer", "1"}},
    RangeTags: []bebber.RangeTag{rT4},
  }

  err = c.Insert(f1, f2, f3, f4, f5, f6, f7)
  if err != nil {
    t.Fatal(err.Error())
  }

  eresult := []bebber.AccFile{
    bebber.AccFile{&invo1, &f1},
    bebber.AccFile{&invo2, &f2},
    bebber.AccFile{&stat1, &f4},
    bebber.AccFile{&stat2, &f5},
    bebber.AccFile{&stat3, &f5},
  }

  result, err := bebber.JoinAccFile(acd, c, false)

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
  doc1 := bebber.FileDoc{
    "test1.txt",
    []bebber.SimpleTag{bebber.SimpleTag{"sTag1"}},
    []bebber.RangeTag{
        bebber.RangeTag{"rTag1", sT, eT},
        bebber.RangeTag{"rTag2", sT, eT},
      },
    []bebber.ValueTag{bebber.ValueTag{"vTag1", "value1"}},
  }
  sT2 := time.Date(2015, time.April, 1, 0, 0, 0, 0, time.UTC)
  eT2 := time.Date(2015, time.April, 2, 0, 0, 0, 0, time.UTC)
  doc2 := bebber.FileDoc{
    "test1.txt",
    []bebber.SimpleTag{bebber.SimpleTag{"sTag1"}},
    []bebber.RangeTag{
        bebber.RangeTag{"rTag1", sT, eT},
        bebber.RangeTag{"rTag2", sT2, eT2},
      },
    []bebber.ValueTag{bebber.ValueTag{"vTag1", "value2"}},
  }
  doc3 := bebber.FileDoc{
    "test1.txt",
    []bebber.SimpleTag{bebber.SimpleTag{"sTag1"}},
    []bebber.RangeTag{},
    []bebber.ValueTag{},
  }
  doc4 := bebber.FileDoc{
    "notinlist.txt",
    []bebber.SimpleTag{},
    []bebber.RangeTag{},
    []bebber.ValueTag{},
  }

  fd := bebber.FileDocsNew([]bebber.FileDoc{doc1, doc2, doc3, doc4})

  findDoc := bebber.FileDoc{
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

  findDoc = bebber.FileDoc{
    Filename: "test1.txt",
    RangeTags: []bebber.RangeTag{bebber.RangeTag{"rTag1", sT, eT}},
  }
  res = res.FindFile(findDoc)
  if len(res.List) != 1 {
    t.Fatal("Expect 1 was ", len(res.List))
  }
}

func TestEmptyAccData(t *testing.T) {
  ad := bebber.AccData{}
  if ad.Empty() == false {
    t.Fatal("Expect true was false")
  }

  ad = bebber.AccData{Belegnummer: "1"}
  if ad.Empty() == true {
    t.Fatal("Expect false was true")
  }
}
