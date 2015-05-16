package bebber

import (
  "os"
  _"fmt"
  "path"
  "path/filepath"
  "time"
  "bytes"
  _"strings"
  "io/ioutil"
  "testing"
  "net/http"
  "net/http/httptest"
  _"encoding/json"

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

  c := session.DB("bebber_test").C("files")
  err = c.Insert(doc1, doc2)
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
    t.Error("Result Error - Status: ", w.Code)
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
