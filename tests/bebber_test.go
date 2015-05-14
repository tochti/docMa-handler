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
  _"encoding/json"

  "gopkg.in/mgo.v2"
  "github.com/rrawrriw/bebber"
  "github.com/gin-gonic/gin"
)

var testDir, err = filepath.Abs(".")

func createTestContext() (c *gin.Context, w *httptest.ResponseRecorder, r *gin.Engine) {
	w = httptest.NewRecorder()
	r = gin.New()
  c = &gin.Context{Engine: r}
	return
}

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
    t.Log("Diff should be [1] but is ", diff)
    t.FailNow()
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
    t.Log("Diff should be [] but is ", diff)
    t.FailNow()
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
  rDoc := `{"Status": "success", 
        "Dir": [
          {
            "Filename": "test1.txt",
            "SimpleTags": [{tag: "sTag1"}],
            "RangeTags": [{tag: "rTag1", 
                           start: "2014-04-01 00:00:00 +0000 UTC", 
                           end: "2014-04-02 00:00:00 +0000 UTC"}],
            "ValueTags": [{tag: "vTag1", value: "value1"}]
          },
          {
            "Filename": "test2.txt",
            "SimpleTags": [{"tag": "sTag1"}],
            "RangeTags": [],
            "ValueTags": []
          },
          {
            "Filename": "test3.txt",
            "SimpleTags": [],
            "RangeTags": [],
            "ValueTags": []
          }
            
        ]}`


  fmt.Println(w.Body.String())
  fmt.Println(strings.TrimSpace(rDoc))
  if strings.TrimSpace(rDoc) != w.Body.String() {
    t.Error("Result Error - Status: ", w.Code)
  }

  /* cleanup */
  os.RemoveAll(tmpDir)
  err = session.DB("bebber_test").DropDatabase()
  if err != nil {
    t.Error(err)
  }
}
