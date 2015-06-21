package bebber

import (
  "os"
  "fmt"
  "path"
  "time"
  "bytes"
  "strings"
  "testing"
  "net/http"
  "io/ioutil"
  "crypto/sha1"
  "encoding/json"

  "gopkg.in/mgo.v2"
  "gopkg.in/mgo.v2/bson"
  "github.com/gin-gonic/gin"
)

func TestLoadAccFile(t *testing.T) {
  /* setup */
  tmpDir, err := ioutil.TempDir(testDir, "accdata")
  defer os.RemoveAll(tmpDir)

  // Create expected invoices 
  invo1 := AccData{
    Belegdatum: time.Date(2013,time.September,29, 0,0,0,0,time.UTC),
    Buchungsdatum: time.Date(2013, time.September,29, 0,0,0,0,time.UTC),
    Belegnummernkreis: "B",
    Belegnummer: "6",
  }
  invo2 := AccData{
    Belegdatum: time.Date(2013,time.September,29, 0,0,0,0,time.UTC),
    Buchungsdatum: time.Date(2013, time.September,29, 0,0,0,0,time.UTC),
    Belegnummernkreis: "B",
    Belegnummer: "8",
  }
  invo3 := AccData{
    Belegdatum: time.Date(2013,time.September,29, 0,0,0,0,time.UTC),
    Buchungsdatum: time.Date(2013, time.September,29, 0,0,0,0,time.UTC),
    Belegnummernkreis: "S",
    Belegnummer: "9",
    Sollkonto: 1210,
    Habenkonto: 0,
  }
  invo4 := AccData{
    Belegdatum: time.Date(2013,time.September,29, 0,0,0,0,time.UTC),
    Buchungsdatum: time.Date(2013, time.September,29, 0,0,0,0,time.UTC),
    Belegnummernkreis: "S",
    Belegnummer: "10",
    Sollkonto: 0,
    Habenkonto: 1211,
  }

  //acd := []AccData{invo1, invo2, stat1, stat2}
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
    ValueTags: []ValueTag{ValueTag{"Belegnummer", "B6"}},
  }
  f2 := FileDoc{
    Filename: "i2.pdf",
    ValueTags: []ValueTag{ValueTag{"Belegnummer", "B8"}},
  }
  f3 := FileDoc{
    Filename: "inone.pdf",
    ValueTags: []ValueTag{ValueTag{"Belegnummer", "19"}},
  }
  sD := time.Date(2013,time.September,29, 0,0,0,0,time.UTC)
  eD := time.Date(2013,time.September,29, 0,0,0,0,time.UTC)
  rT1 := RangeTag{"Belegzeitraum", sD, eD}
  f4 := FileDoc{
    Filename: "s1.pdf",
    ValueTags: []ValueTag{ValueTag{"Kontonummer", "1210"}},
    RangeTags: []RangeTag{rT1},
  }
  sD = time.Date(2014,time.September,29, 0,0,0,0,time.UTC)
  eD = time.Date(2014,time.September,29, 0,0,0,0,time.UTC)
  rT2 := RangeTag{"Belegzeitraum", sD, eD}
  f5 := FileDoc{
    Filename: "s2.pdf",
    ValueTags: []ValueTag{ValueTag{"Kontonummer", "1211"}},
    RangeTags: []RangeTag{rT2},
  }
  sD = time.Date(2014,time.April,20, 0,0,0,0,time.UTC)
  eD = time.Date(2014,time.April,24, 0,0,0,0,time.UTC)
  rT3 := RangeTag{"Belegzeitraum", sD, eD}
  f6 := FileDoc{
    Filename: "snone.pdf",
    ValueTags: []ValueTag{ValueTag{"Kontonummer", "10001"}},
    RangeTags: []RangeTag{rT3},
  }
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

  f1.Filename = path.Join(tmpDir, f1.Filename)
  f2.Filename = path.Join(tmpDir, f2.Filename)
  f3.Filename = path.Join(tmpDir, f3.Filename)
  f4.Filename = path.Join(tmpDir, f4.Filename)
  f5.Filename = path.Join(tmpDir, f5.Filename)
  f6.Filename = path.Join(tmpDir, f6.Filename)
  f7.Filename = path.Join(tmpDir, f7.Filename)

  accFiles := []AccFile{
    AccFile{invo1, f1},
    AccFile{invo2, f2},
    AccFile{invo3, f4},
    AccFile{invo4, f5},
  }

  eRes := LoadAccFilesResponse{
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
  r.POST("/", LoadAccFiles)
  body := bytes.NewBufferString("")
  res := PerformRequest(r, "POST", "/", body)

  /* Compare */
  // Fix == should be !=
  if string(eresult) == strings.TrimSpace(res.Body.String()) {
    t.Fatal("Expect -->", string(eresult), "<--\nwas\n-->", strings.TrimSpace(res.Body.String()), "<--")
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
  tmp := RangeTag{
            "rT1",
            time.Date(2014, time.April, 6, 0, 0, 0, 0, time.UTC),
            time.Date(2014, time.April, 7, 0, 0, 0, 0, time.UTC),
          }
  doc := FileDoc{Filename: "test.txt",
                 SimpleTags: []SimpleTag{SimpleTag{"sT1"}},
                 ValueTags: []ValueTag{ValueTag{"vT1", "v1"}},
                 RangeTags: []RangeTag{tmp},
                }

  collection.Insert(doc)

  /* perfom request */
  s := gin.New()
  s.POST("/", AddTags)
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
  rD := FileDoc{}
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

func Test_UserAuth_OK(t *testing.T) {
  os.Setenv("BEBBER_DB_SERVER", "127.0.0.1")
  os.Setenv("BEBBER_DB_NAME", "bebber_test")
  session, err := mgo.Dial("127.0.0.1")
  if err != nil {
    t.Fatal(err)
  }
  defer session.Close()
  sha1Pass := fmt.Sprintf("%x", sha1.Sum([]byte("test")))

  user := bson.M{"username": "loveMaster_999", "password": sha1Pass}
  usersC := session.DB("bebber_test").C("users")
  usersC.Insert(user)
  defer session.DB("bebber_test").DropDatabase()

  h := gin.New()
  h.POST("/", Login)
  body := bytes.NewBufferString(`{"Username":"loveMaster_999","Password":"test"}`)
  resp := PerformRequest(h, "POST", "/", body)

  if resp.Code != 200 {
    t.Fatal("Response code should be 200 was", resp.Code)
  }

  _, ok := resp.HeaderMap["Set-Cookie"]

  if ok == false {
    t.Fatal("Expect a cookie in header but was", resp.HeaderMap,"\n",resp)
  }
  if strings.Contains(resp.HeaderMap["Set-Cookie"][0], "XSRF-TOKEN") == false {
    t.Fatal("Expect XSRF-TOKEN field in header but was", resp.HeaderMap["Set-Cookie"][0])
  }
}

func Test_UserAuth_Fail(t *testing.T) {
  os.Setenv("BEBBER_DB_SERVER", "127.0.0.1")
  os.Setenv("BEBBER_DB_NAME", "bebber_test")
  session, err := mgo.Dial("127.0.0.1")
  if err != nil {
    t.Fatal(err)
  }
  defer session.Close()
  sha1Pass := fmt.Sprintf("%x", sha1.Sum([]byte("test")))
  user := bson.M{"Username": "loveMaster_999", "Password": sha1Pass}
  usersC := session.DB("bebber_test").C("users")
  usersC.Insert(user)
  defer session.DB("bebber_test").DropDatabase()

  h := gin.New()
  h.POST("/", Login)
  body := bytes.NewBufferString(`{"username":"loveMaster_999","password":"wrong"}`)
  resp := PerformRequest(h, "POST", "/", body)

  if resp.Code != 200 {
    t.Fatal("Response code should be 200 was", resp.Code)
  }

  expectJsonResp := `{"Status":"fail","Msg":"Wrong username or password"}` + "\n"
  if resp.Body.String() !=  expectJsonResp {
    t.Fatal("Expect", expectJsonResp, "was", resp.Body.String())
  }

}

func Test_UserHandler_OK(t *testing.T) {
  globals := MakeTestGlobals(t)
  defer globals.MongoDB.Session.Close()
  db := globals.MongoDB.Session.DB(TestDBName)
  defer db.DropDatabase()


  sha1Pass := fmt.Sprintf("%x", sha1.Sum([]byte("test")))
  user := User{Username: "hitman", Password: sha1Pass}
  user1 := User{Username: "catwomen", Password: sha1Pass}
  user2 := User{Username: "loveMaster_999", Password: sha1Pass}
  usersColl := db.C(UsersCollection)
  err := usersColl.Insert(user, user1, user2)
  if err != nil {
    t.Fatal(err.Error())
  }

  handler := gin.New()
  userHandler := MakeGlobalsHandler(UserHandler, globals)
  handler.GET("/:name", userHandler)
  request := TestRequest{
                Body: "",
                Handler: handler,
              }
  response := request.Send("GET", "/hitman")

  expect := `{"Username":"hitman","Password":""}`
  if strings.Contains(response.Body.String(), expect) {
    t.Fatal("Expect", expect, "was", response.Body.String())
  }

}

func Test_SearchHandler_OK(t *testing.T) {
  globals := MakeTestGlobals(t)
  defer globals.MongoDB.Session.Close()

  db := globals.MongoDB.Session.DB(TestDBName)
  defer db.DropDatabase()

  doc := bson.M{"Blue": "House"}
  testColl := db.C(FilesCollection)
  testColl.Insert(doc)

  handler := gin.New()
  handler.POST("/", MakeGlobalsHandler(SearchHandler, globals))
  request := TestRequest{
    Body: `{"Blue": "House"}`,
    Header: http.Header{},
    Handler: handler,
  }
  response := request.SendWithToken("POST", "/", "123")

  expect := `[{"Blue": "House"}]`
  if strings.Contains(response.Body.String(), expect) {
    t.Fatal("Expect", expect, "was", response.Body.String())
  }

}

// Make a new Doc in DB, simple case
func Test_DocMakeHandler_MakeNewDocOK(t *testing.T) {
  globals := MakeTestGlobals(t)
  defer globals.MongoDB.Session.Close()

  date := time.Date(2015, 1, 1, 0,0,0,0, time.UTC)
  docInfos := DocInfos {
                DateOfScan: date,
                DateOfReceipt: date,
              }
  docRequest := Doc{
              Name: "darkmoon.txt",
              Infos: docInfos,
              Barcode: "darkmoon",
              Note: "There was a man...",
              Labels: []Label{Label("old"), Label("story")},
              AccountData: DocAccountData{DocNumber: "0815"},
            }

  docRequestJSON, err := json.Marshal(docRequest)
  if err != nil {
    t.Fatal(err.Error())
  }

  handler := gin.New()
  handler.POST("/Doc", MakeGlobalsHandler(DocMakeHandler, globals))
  request := TestRequest{
                Body: string(docRequestJSON),
                Header: http.Header{},
                Handler: handler,
              }
  response := request.Send("POST", "/Doc")

  if strings.Contains(response.Body.String(), `"Status":"success"`) == false {
    t.Fatal("Expect to contains \"Status\":\"success\" was",
            response.Body.String())
  }

  successResponse := MongoDBSuccessResponse{}
  err = json.Unmarshal([]byte(response.Body.String()), &successResponse)
  if err != nil {
    t.Fatal(err.Error())
  }
  db := globals.MongoDB.Session.Copy().DB(TestDBName)
  defer db.DropDatabase()

  doc := Doc{}
  err = db.C(DocsColl).Find(bson.M{"_id": bson.ObjectIdHex(successResponse.DocID)}).One(&doc)
  if err != nil {
    t.Fatal(err.Error())
  }

  docJSON, err := json.Marshal(doc)
  if err != nil {
    t.Fatal(err.Error())
  }

  docRequest.ID = bson.ObjectIdHex(successResponse.DocID)
  docRequestJSON, err = json.Marshal(docRequest)
  if err != nil {
    t.Fatal(err.Error())
  }

  if string(docRequestJSON) != string(docJSON) {
    t.Fatal("Expect", string(docRequestJSON), "was", string(docJSON))
  }

}


// Try to make a doc although a doc with the same name already exists
func Test_DocMakeHandler_AlreadyExistsFail(t *testing.T) {
  globals := MakeTestGlobals(t)
  defer globals.MongoDB.Session.Close()

  db := globals.MongoDB.Session.Copy().DB(TestDBName)
  defer db.DropDatabase()

  docInfos := DocInfos {
                DateOfScan: time.Now(),
                DateOfReceipt: time.Now(),
              }
  doc := Doc{Name: "EarlyBird.txt", Infos: docInfos}
  err := db.C(DocsColl).Insert(doc)
  if err != nil {
    t.Fatal(err.Error())
  }
  docJSON, err := json.Marshal(doc)
  if err != nil {
    t.Fatal(err.Error())
  }
  handler := gin.New()
  handler.POST("/Doc", MakeGlobalsHandler(DocMakeHandler, globals))
  request := TestRequest{
                Body: string(docJSON),
                Header: http.Header{},
                Handler: handler,
              }

  response := request.Send("POST", "/Doc")

  failMsg := "Document already exists!"
  if strings.Contains(response.Body.String(), failMsg) == false {
    t.Fatal("Expect Msg", failMsg, "was", response)
  }
}

// Try to make a doc send crazy request 
func Test_DocMakeHandler_CrazyFail(t *testing.T) {
  globals := MakeTestGlobals(t)
  defer globals.MongoDB.Session.Close()
  defer globals.MongoDB.Session.Copy().DB(TestDBName).DropDatabase()

  handler := gin.New()
  handler.POST("/Doc", MakeGlobalsHandler(DocMakeHandler, globals))
  testDate := time.Date(2015,1,1,0,0,0,0,time.UTC)
  dateStr := testDate.Format(time.RFC3339Nano)
  fmt.Println(dateStr)
  request := TestRequest{
                Body: `{"Name": "chunk", "Infos": {"DateOfReceipt": "`+ dateStr +`", "DateOfScan": "`+ dateStr +`"}, "Fail": "yes!"}`,
                Header: http.Header{},
                Handler: handler,
              }

  response := request.Send("POST", "/Doc")

  failMsg := "Missing a name!"
  if strings.Contains(response.Body.String(), failMsg) == false {
    t.Fatal("Expect Msg", failMsg, "was", response)
  }
}

// Try to make a doc without a doc name 
func Test_DocMakeHandler_WithoutNameFail(t *testing.T) {
  globals := MakeTestGlobals(t)
  defer globals.MongoDB.Session.Close()

  doc := Doc{}
  docJSON, err := json.Marshal(doc)
  if err != nil {
    t.Fatal(err.Error())
  }
  handler := gin.New()
  handler.POST("/Doc", MakeGlobalsHandler(DocMakeHandler, globals))
  request := TestRequest{
                Body: string(docJSON),
                Header: http.Header{},
                Handler: handler,
              }

  response := request.Send("POST", "/Doc")

  failMsg := "Missing a name!"
  if strings.Contains(response.Body.String(), failMsg) == false {
    t.Fatal("Expect Msg", failMsg, "was", response)
  }
}

// Try to make a doc without a info field
func Test_DocMakeHandler_WithoutInfosFail(t *testing.T) {
  globals := MakeTestGlobals(t)
  defer globals.MongoDB.Session.Close()

  doc := Doc{Name: "NoInfo.pdf"}
  docJSON, err := json.Marshal(doc)
  if err != nil {
    t.Fatal(err.Error())
  }
  handler := gin.New()
  handler.POST("/Doc", MakeGlobalsHandler(DocMakeHandler, globals))
  request := TestRequest{
                Body: string(docJSON),
                Header: http.Header{},
                Handler: handler,
              }

  response := request.Send("POST", "/Doc")

  failMsg := "Missing the infos field!"
  if strings.Contains(response.Body.String(), failMsg) == false {
    t.Fatal("Expect Msg", failMsg, "was", response)
  }
}

// Change a existing doc. Append not existing fields
func Test_DocChangeHandler_OK(t *testing.T) {
  globals := MakeTestGlobals(t)
  defer globals.MongoDB.Session.Close()

  db := globals.MongoDB.Session.Copy().DB(TestDBName)
  defer db.DropDatabase()
  docsColl := db.C(DocsColl)

  docID := bson.NewObjectId()
  docTmp := Doc{ID: docID, Name: "Touchme.txt"}
  err := docsColl.Insert(docTmp)
  if err != nil {
    t.Fatal(err.Error())
  }

  changeRequest := DocChangeRequest{
                    Name: "Touchme.txt",
                    Barcode: "Touchme",
                    Labels: []Label{Label("oohh"), Label("aahh")},
                    Note: "Don't stop!",
                    AccountData: DocAccountData{DocNumber: "VIP"},
                  }
  changeRequestJSON, err := json.Marshal(changeRequest)
  if err != nil {
    t.Fatal(err.Error())
  }

  handler := gin.New()
  handler.PATCH("/Doc", MakeGlobalsHandler(DocChangeHandler, globals))
  request := TestRequest{
                Body: string(changeRequestJSON),
                Header: http.Header{},
                Handler: handler,
              }

  response := request.Send("PATCH", "/Doc")

  if strings.Contains(response.Body.String(), "success") == false {
    t.Fatal("Expect success response was", response)
  }

  doc := Doc{Name: "Touchme.txt"}
  err = doc.Find(db)
  if err != nil {
    t.Fatal(err.Error())
  }

  doc.ID = docID
  docJSON, err := json.Marshal(doc)
  if err != nil {
    t.Fatal(err.Error())
  }

  changeRequest.ID = docID
  changeRequestJSON, err = json.Marshal(changeRequest)
  if err != nil {
    t.Fatal(err.Error())
  }

  if string(changeRequestJSON) != string(docJSON) {
    t.Fatal("Expect", string(changeRequestJSON), "was", string(docJSON))
  }
}

// Read Doc
func Test_DocReadHandler_OK(t *testing.T) {
  globals := MakeTestGlobals(t)
  session := globals.MongoDB.Session.Copy()
  defer session.Close()

  db := session.DB(TestDBName)
  defer db.DropDatabase()
  docsColl := db.C(DocsColl)

  docID := bson.NewObjectId()
  expectDoc := Doc{ID: docID, Name: "Seek.pdf", Labels: []Label{}}

  err := docsColl.Insert(expectDoc)
  if err != nil {
    t.Fatal(err.Error())
  }

  handler := gin.New()
  handler.GET("/Doc/:name", MakeGlobalsHandler(DocReadHandler, globals))

  request := TestRequest{
                Body: "",
                Header: http.Header{},
                Handler: handler,
              }
  response := request.Send("GET", "/Doc/Seek.pdf")

  docReadResponse := DocReadResponse{Status: "success", Doc: expectDoc}
  docReadResponseJSON, err := json.Marshal(docReadResponse)
  if err != nil {
    t.Fatal(err.Error())
  }

  if string(docReadResponseJSON)+"\n" != response.Body.String() {
    t.Fatal("Expect", string(docReadResponseJSON), "was", response.Body.String())
  }
}

// Remove Doc
func Test_DocRemoveHandler_OK(t *testing.T) {
  globals := MakeTestGlobals(t)
  session := globals.MongoDB.Session.Copy()
  defer session.Close()

  tmpDir, err := ioutil.TempDir(testDir, "remove")
  defer os.RemoveAll(tmpDir)
  globals.Config["FILES_DIR"] = tmpDir

  tmpFile, err := ioutil.TempFile(tmpDir, "remove")
  fileBase := path.Base(tmpFile.Name())

  db := session.DB(TestDBName)
  defer db.DropDatabase()
  docsColl := db.C(DocsColl)

  docID := bson.NewObjectId()
  expectDoc := Doc{ID: docID, Name: fileBase, Labels: []Label{}}

  err = docsColl.Insert(expectDoc)
  if err != nil {
    t.Fatal(err.Error())
  }

  handler := gin.New()
  handler.DELETE("/Doc/:name", MakeGlobalsHandler(DocRemoveHandler, globals))

  request := TestRequest{
                Body: "",
                Header: http.Header{},
                Handler: handler,
              }
  url := "/Doc/"+ fileBase
  response := request.Send("DELETE", url)

  if strings.Contains(response.Body.String(), "success") == false {
    t.Fatal("Expect successs response was", response)
  }

  err = expectDoc.Find(db)
  if err == nil {
    t.Fatal("Expect Cannot find document error was", err)
  }
  if strings.Contains(err.Error(), "Cannot find document") == false {
    t.Fatal("Expect Cannot find document error was", err)
  }

  tmpFilePath := path.Join(tmpDir, fileBase)
  if _, err := os.Stat(tmpFilePath); os.IsNotExist(err) == false {
    t.Fatal("Expect file", tmpFilePath, "to be delete")
  }

}
