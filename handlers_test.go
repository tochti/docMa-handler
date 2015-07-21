package bebber

import (
  "os"
  "fmt"
  "path"
  "time"
  "strings"
  "testing"
  "net/http"
  "io/ioutil"
  "crypto/sha1"
  "encoding/json"

  "gopkg.in/mgo.v2/bson"
  "github.com/gin-gonic/gin"
)
/*
func Test_LoadAccFile(t *testing.T) {
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
  err = session.DB(TestDBName).DropDatabase()
  if err != nil {
    t.Fatal(err.Error())
  }

  c := session.DB(TestDBName).C("files")

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

  r := gin.New()
  r.POST("/", LoadAccFiles)
  body := bytes.NewBufferString("")
  res := PerformRequest(r, "POST", "/", body)

  // Fix == should be !=
  if string(eresult) == strings.TrimSpace(res.Body.String()) {
    t.Fatal("Expect -->", string(eresult), "<--\nwas\n-->", strings.TrimSpace(res.Body.String()), "<--")
  }

}
*/

func Test_LoginHandler_OK(t *testing.T) {
  globals := MakeTestGlobals(t)
  session := globals.MongoDB.Session.Copy()
  defer session.Close()

  db := session.DB(TestDBName)
  defer db.DropDatabase()

  sha1Pass := fmt.Sprintf("%x", sha1.Sum([]byte("test")))
  user := bson.M{"username": "loveMaster_999", "password": sha1Pass}
  usersColl := db.C(UsersColl)
  usersColl.Insert(user)

  h := gin.New()
  h.POST("/", MakeGlobalsHandler(LoginHandler, globals))
  request := TestRequest{
                Body: `{"Username":"loveMaster_999","Password":"test"}`,
                Header: http.Header{},
                Handler: h,
              }
  resp := request.Send("POST", "/")

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

func Test_LoginHandler_Fail(t *testing.T) {
  globals := MakeTestGlobals(t)
  session := globals.MongoDB.Session.Copy()
  defer session.Close()
  db := session.DB(TestDBName)
  defer db.DropDatabase()

  sha1Pass := fmt.Sprintf("%x", sha1.Sum([]byte("test")))
  user := bson.M{"Username": "loveMaster_999", "Password": sha1Pass}
  usersColl := db.C(UsersColl)
  usersColl.Insert(user)

  h := gin.New()
  h.POST("/", MakeGlobalsHandler(LoginHandler, globals))
  request := TestRequest{
                Body: `{"username":"loveMaster_999","password":"wrong"}`,
                Header: http.Header{},
                Handler: h,
              }

  resp := request.Send("POST", "/")

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
  usersColl := db.C(UsersColl)
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
  if strings.Contains(response.Body.String(), expect) == false {
    t.Fatal("Expect", expect, "was", response.Body.String())
  }

}

func Test_SearchHandler_OK(t *testing.T) {
  globals := MakeTestGlobals(t)
  defer globals.MongoDB.Session.Close()

  db := globals.MongoDB.Session.DB(TestDBName)
  defer db.DropDatabase()

  doc := bson.M{"Blue": "House"}
  testColl := db.C(DocsColl)
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
    AccountData: DocAccountData{
      DocNumbers: []string{"0815"},
    },
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
  docTmp := Doc{
    ID: docID,
    Name: "Touchme.txt",
    Barcode: "Codey",
    Note: "Nutty",
  }
  err := docsColl.Insert(docTmp)
  if err != nil {
    t.Fatal(err.Error())
  }

  handler := gin.New()
  handler.PATCH("/Doc", MakeGlobalsHandler(DocChangeHandler, globals))
  changeRequest := `{
    "Name": "Touchme.txt",
    "Note": "Mutti",
    "Labels": ["label1"],
    "AccountData":{"DocNumbers": ["123"]}
  }`
  request := TestRequest{
                Body: changeRequest,
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

  expectDoc := Doc{
    ID: docID,
    Name: "Touchme.txt",
    Barcode: "Codey",
    Note: "Mutti",
    AccountData: DocAccountData{
      DocNumbers: []string{"123"},
    },
    Labels: []Label{"label1"},
  }
  expectDocJSON, err := json.Marshal(expectDoc)
  if err != nil {
    t.Fatal(err.Error())
  }

  if string(expectDocJSON) != string(docJSON) {
    t.Fatal("Expect", string(expectDocJSON), "was", string(docJSON))
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
  expectDoc := Doc{
    ID: docID,
    Name: "Seek.pdf",
    Labels: []Label{},
    AccountData: DocAccountData{
      DocNumbers: []string{},
    },
  }

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

// Rename Doc
func Test_DocRenameHandler_OK(t *testing.T) {
  globals := MakeTestGlobals(t)
  session := globals.MongoDB.Session.Copy()
  defer session.Close()

  tmpDir, err := ioutil.TempDir(testDir, "rename")
  defer os.RemoveAll(tmpDir)
  globals.Config["FILES_DIR"] = tmpDir

  tmpFile, err := ioutil.TempFile(tmpDir, "rename")
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
  handler.PATCH("/Doc/Rename", MakeGlobalsHandler(DocRenameHandler, globals))

  request := TestRequest{
                Body: `{"Name":"`+ fileBase +`", "NewName":"Simple" }`,
                Header: http.Header{},
                Handler: handler,
              }
  response := request.Send("PATCH", "/Doc/Rename")

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

  if _, err := os.Stat(tmpFile.Name()); os.IsNotExist(err) == false {
    t.Fatal("Expect file", tmpFile.Name(), "to be delete")
  }

  newFilePath := path.Join(tmpDir, "Simple")
  if _, err := os.Stat(newFilePath); os.IsNotExist(err) {
    t.Fatal("Expect file", newFilePath, "to exist")
  }

  doc := Doc{Name: "Simple"}
  err = doc.Find(db)
  if err != nil {
    t.Fatal(err.Error())
  }

}
// Append labels to a existing list of labels
func Test_DocAppendLabelsHandler_OK(t *testing.T) {
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

  handler := gin.New()
  handler.PATCH("/Doc/Labels",
                MakeGlobalsHandler(DocAppendLabelsHandler, globals))
  appendLabelsRequest := `{"Name": "Hoocker", "Labels":["label2", "label3"]}`
  request := TestRequest{
                Body: appendLabelsRequest,
                Header: http.Header{},
                Handler: handler,
              }

  response := request.Send("PATCH", "/Doc/Labels")

  if strings.Contains(response.Body.String(), "success") == false {
    t.Fatal("Expect success response was", response)
  }

  doc := Doc{Name: "Hoocker"}
  err = doc.Find(db)
  if err != nil {
    t.Fatal(err.Error())
  }

  if len(doc.Labels) != 3 {
    t.Fatal("Expect 3 labels was", doc.Labels)
  }
}

// Remove labels
func Test_DocRemoveLabelHandler_OK(t *testing.T) {
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

  handler := gin.New()
  handler.DELETE("/DocLabels/:name/:label",
                MakeGlobalsHandler(DocRemoveLabelHandler, globals))
  request := TestRequest{
                Body: "",
                Header: http.Header{},
                Handler: handler,
              }

  response := request.Send("DELETE", "/DocLabels/Hoocker/l2")

  if strings.Contains(response.Body.String(), "success") == false {
    t.Fatal("Expect success response was", response)
  }

  doc := Doc{Name: "Hoocker"}
  err = doc.Find(db)
  if err != nil {
    t.Fatal(err.Error())
  }

  if len(doc.Labels) != 2 {
    t.Fatal("Expect 1 labels was", doc.Labels)
  }
  if doc.Labels[0] != "l1" || doc.Labels[1] != "l3" {
    t.Fatal("Expect l1 or l3 was", doc.Labels[0], doc.Labels[1])
  }
}

// Append docnumbers to a existing list of labels
func Test_DocAppendDocNumbersHandler_OK(t *testing.T) {
  globals := MakeTestGlobals(t)
  session := globals.MongoDB.Session.Copy()
  defer session.Close()

  db := session.DB(TestDBName)
  defer db.DropDatabase()

  docTmp := Doc{
    Name: "Hoocker",
    AccountData: DocAccountData{
      DocNumbers: []string{},
    },
  }
  docsColl := db.C(DocsColl)
  err := docsColl.Insert(docTmp)

  handler := gin.New()
  handler.PATCH("/DocNumbers",
                MakeGlobalsHandler(DocAppendDocNumbersHandler, globals))
  appendLabelsRequest := `{
    "Name": "Hoocker",
    "DocNumbers": ["1","2"]
  }`
  request := TestRequest{
      Body: appendLabelsRequest,
      Header: http.Header{},
      Handler: handler,
    }

  response := request.Send("PATCH", "/DocNumbers")

  if strings.Contains(response.Body.String(), "success") == false {
    t.Fatal("Expect success response was", response)
  }

  doc := Doc{Name: "Hoocker"}
  err = doc.Find(db)
  if err != nil {
    t.Fatal(err.Error())
  }

  if len(doc.AccountData.DocNumbers) != 2 {
    t.Fatal("Expect 2 doc numbers was", doc.AccountData.DocNumbers)
  }
}

// Remove docnumber
func Test_DocRemoveDocNumberHandler_OK(t *testing.T) {
  globals := MakeTestGlobals(t)
  session := globals.MongoDB.Session.Copy()
  defer session.Close()

  db := session.DB(TestDBName)
  defer db.DropDatabase()

  docTmp := Doc{
    Name: "Hoocker",
    AccountData: DocAccountData{
      DocNumbers: []string{"1", "2", "3"},
    },
  }
  docsColl := db.C(DocsColl)
  err := docsColl.Insert(docTmp)

  handler := gin.New()
  handler.DELETE("/DocNumber/:name/:number",
    MakeGlobalsHandler(DocRemoveDocNumberHandler, globals))
  request := TestRequest{
    Body: "",
    Header: http.Header{},
    Handler: handler,
  }

  response := request.Send("DELETE", "/DocNumber/Hoocker/2")

  if strings.Contains(response.Body.String(), "success") == false {
    t.Fatal("Expect success response was", response)
  }

  doc := Doc{Name: "Hoocker"}
  err = doc.Find(db)
  if err != nil {
    t.Fatal(err.Error())
  }

  if len(doc.AccountData.DocNumbers) != 2 {
    t.Fatal("Expect 1 labels was", doc.Labels)
  }
  if (doc.AccountData.DocNumbers[0] != "1") ||
     (doc.AccountData.DocNumbers[1] != "3" ) {
    t.Fatal("Expect 1 or 3 was", doc.Labels[0], doc.Labels[1])
  }
}

func Test_AccProcessMakeHandler_OK(t *testing.T) {
  globals := MakeTestGlobals(t)
  session := globals.MongoDB.Session.Copy()
  defer session.Close()

  db := session.DB(TestDBName)
  defer db.DropDatabase()

  handler := gin.New()
  handler.POST("/AccProcess",
    MakeGlobalsHandler(AccProcessMakeHandler, globals))
  requestBody := AccProcessMakeRequest{}
  requestBodyJSON, err := json.Marshal(requestBody)
  if err != nil {
    t.Fatal(err.Error())
  }

  request := TestRequest{
    Body: string(requestBodyJSON),
    Header: http.Header{},
    Handler: handler,
  }

  response := request.Send("POST", "/AccProcess")

  responseBody := AccProcessMakeResponse{}

  tmp := []byte(response.Body.String())
  err = json.Unmarshal(tmp, &responseBody)
  if err != nil {
    t.Fatal(err.Error())
  }

  if responseBody.Status != "success" {
    t.Fatal("Expect success was", responseBody.Status)
  }

  accProcessColl := db.C(AccProcessColl)
  docID := bson.ObjectIdHex(responseBody.DocID)
  query := accProcessColl.Find(bson.M{"_id": docID})
  n, err := query.Count()
  if err != nil {
    t.Fatal(err.Error())
  }
  if n != 1 {
    t.Fatal("Expect 1 was", n)
  }

}

func Test_AccProcessMakeHandler_AlreadyExistsFail(t *testing.T) {
  globals := MakeTestGlobals(t)
  session := globals.MongoDB.Session.Copy()
  defer session.Close()

  db := session.DB(TestDBName)
  defer db.DropDatabase()

  handler := gin.New()
  handler.POST("/AccProcess",
    MakeGlobalsHandler(AccProcessMakeHandler, globals))
  requestBody := AccProcessMakeRequest{}
  requestBodyJSON, err := json.Marshal(requestBody)
  if err != nil {
    t.Fatal(err.Error())
  }

  request := TestRequest{
    Body: string(requestBodyJSON),
    Header: http.Header{},
    Handler: handler,
  }

  response := request.Send("POST", "/AccProcess")
  if response.Code != 200 {
    t.Fatal("Expect status 200 was", response.Code)
  }

  response = request.Send("POST", "/AccProcess")

  responseBody := FailResponse{}

  tmp := []byte(response.Body.String())
  err = json.Unmarshal(tmp, &responseBody)
  if err != nil {
    t.Fatal(err.Error())
  }

  if responseBody.Status != "fail" {
    t.Fatal("Expect fail was", responseBody.Status)
  }

  errMsg := "Account process already exists!"
  if strings.Contains(responseBody.Msg, errMsg) == false {
    t.Fatal("Expect", errMsg, "was", responseBody.Msg)
  }
}

func Test_FindAccProcessByDocNumberHandler_OK(t *testing.T) {
  globals := MakeTestGlobals(t)
  session := globals.MongoDB.Session.Copy()
  defer session.Close()

  db := session.DB(TestDBName)
  defer db.DropDatabase()

  err := ImportAccProcess(db, path.Join(testDir, "export.csv"))
  if err != nil {
    t.Fatal(err.Error())
  }

  handler := gin.New()
  handler.GET("/AccProcess/FindByDocNumber/:number",
    MakeGlobalsHandler(AccProcessFindByDocNumberHandler, globals))

  request := TestRequest{
    Body: "",
    Header: http.Header{},
    Handler: handler,
  }

  response := request.Send("GET", "/AccProcess/FindByDocNumber/13")

  if response.Code != 200 {
    t.Fatal("Expect 200 was", response.Code)
  }

  result := AccProcessFindByDocNumberResponse{}
  err = json.Unmarshal([]byte(response.Body.String()), &result)
  if err != nil {
    t.Fatal(err.Error())
  }

  if len(result.AccProcessList) != 1 {
    t.Fatal("Expect len of 1 was", len(result.AccProcessList))
  }

  expectDocNumber := "13"
  if result.AccProcessList[0].DocNumber != expectDocNumber {
    t.Fatal("Expect", expectDocNumber, "was", result.AccProcessList[0].DocNumber)
  }

  expectPostingText := "Miete 09.2013"
  if result.AccProcessList[0].PostingText != expectPostingText {
    t.Fatal("Expect", expectPostingText, "was", result.AccProcessList[0].PostingText)
  }
}

func Test_FindAccProcessByAccNumberHandler_OK(t *testing.T) {
  globals := MakeTestGlobals(t)
  session := globals.MongoDB.Session.Copy()
  defer session.Close()

  db := session.DB(TestDBName)
  defer db.DropDatabase()

  err := ImportAccProcess(db, path.Join(testDir, "export.csv"))
  if err != nil {
    t.Fatal(err.Error())
  }

  handler := gin.New()
  handler.GET("/AccProcess/FindByAccNumber/:from/:to/:number",
    MakeGlobalsHandler(AccProcessFindByAccNumberHandler, globals))

  request := TestRequest{
    Body: "",
    Header: http.Header{},
    Handler: handler,
  }

  response := request.Send("GET",
    "/AccProcess/FindByAccNumber/31082013/01092013/71002")

  if response.Code != 200 {
    t.Fatal("Expect 200 was", response.Code)
  }

  result := AccProcessFindByAccNumberResponse{}
  err = json.Unmarshal([]byte(response.Body.String()), &result)
  if err != nil {
    t.Fatal(err.Error())
  }

  if len(result.AccProcessList) != 1 {
    t.Fatal("Expect len of 1 was", len(result.AccProcessList))
  }

  expectDocNumber := "13"
  if result.AccProcessList[0].DocNumber != expectDocNumber {
    t.Fatal("Expect", expectDocNumber, "was", result.AccProcessList[0].DocNumber)
  }

  expectPostingText := "Miete 09.2013"
  if result.AccProcessList[0].PostingText != expectPostingText {
    t.Fatal("Expect", expectPostingText, "was", result.AccProcessList[0].PostingText)
  }
}

func Test_DocNumberProposalCurrHandler_OK(t *testing.T) {
  globals := MakeTestGlobals(t)
  session := globals.MongoDB.Session.Copy()
  defer session.Close()
  db := session.DB(TestDBName)
  defer db.DropDatabase()

  expectProposal := 1234
  docNumberProposal := DocNumberProposal(expectProposal)
  err := docNumberProposal.Save(db)
  if err != nil {
    t.Fatal(err.Error())
  }

  handler := gin.New()
  handler.GET(
    "/",
    MakeGlobalsHandler(DocNumberProposalCurrHandler, globals),
  )
  header := http.Header{}

  request := TestRequest{
    Body: "",
    Handler: handler,
    Header: header,
  }
  response := request.Send("GET", "/")

  body := response.Body.String()
  bodyJSON := DocNumberProposalCurrResponse{}
  err = json.Unmarshal([]byte(body), &bodyJSON)
  if err != nil {
    t.Fatal(err.Error())
  }

  if "success" != bodyJSON.Status {
    t.Fatal("Expect success was", bodyJSON.Status)
  }

  if expectProposal != bodyJSON.Proposal {
    t.Fatal("Expect", expectProposal, "was", bodyJSON.Proposal)
  }

}

func Test_DocNumberProposalCurrHandler_NoneFail(t *testing.T) {
  globals := MakeTestGlobals(t)

  handler := gin.New()
  handler.GET(
    "/",
    MakeGlobalsHandler(DocNumberProposalCurrHandler, globals),
  )
  header := http.Header{}

  request := TestRequest{
    Body: "",
    Handler: handler,
    Header: header,
  }
  response := request.Send("GET", "/")

  expectBody := `{"Status":"fail","Msg":"Cannot find DocNumberProposal"}` +
    "\n"
  body := response.Body.String()

  if expectBody != body {
    t.Fatal("Expect", expectBody, "was", body)
  }

}

func Test_DocNumberProposalChangeHandler_OK(t *testing.T) {
  globals := MakeTestGlobals(t)
  session := globals.MongoDB.Session.Copy()
  defer session.Close()
  db := session.DB(TestDBName)
  defer db.DropDatabase()

  docNumberProposal := DocNumberProposal(123)
  err := docNumberProposal.Save(db)
  if err != nil {
    t.Fatal(err.Error())
  }

  handler := gin.New()
  handler.PUT(
    "/",
    MakeGlobalsHandler(DocNumberProposalChangeHandler, globals),
  )
  header := http.Header{}

  expectProposal := 1234
  requestBody, err := json.Marshal(
    DocNumberProposalChangeRequest{expectProposal},
  )
  if err != nil {
    t.Fatal(err.Error())
  }
  request := TestRequest{
    Body: string(requestBody),
    Handler: handler,
    Header: header,
  }
  response := request.Send("PUT", "/")

  body := response.Body.String()
  bodyJSON := SuccessResponse{}
  err = json.Unmarshal([]byte(body), &bodyJSON)
  if err != nil {
    t.Fatal(err.Error())
  }

  if "success" != bodyJSON.Status {
    t.Fatal("Expect success was", bodyJSON.Status)
  }

  no, err := docNumberProposal.Curr(db)
  if err != nil {
    t.Fatal(err.Error())
  }

  if expectProposal != int(no) {
    t.Fatal("Expect", expectProposal, "was", int(no))
  }
}

func Test_DocNumberProposalChangeHandler_WrongTypeFail(t *testing.T) {
  globals := MakeTestGlobals(t)
  session := globals.MongoDB.Session.Copy()
  defer session.Close()
  db := session.DB(TestDBName)
  defer db.DropDatabase()

  docNumberProposal := DocNumberProposal(123)
  err := docNumberProposal.Save(db)
  if err != nil {
    t.Fatal(err.Error())
  }

  handler := gin.New()
  handler.PUT(
    "/",
    MakeGlobalsHandler(DocNumberProposalChangeHandler, globals),
  )
  header := http.Header{}

  request := TestRequest{
    Body: `{"Proposal": "DE1234"}`,
    Handler: handler,
    Header: header,
  }
  response := request.Send("PUT", "/")

  body := response.Body.String()
  bodyJSON := FailResponse{}
  err = json.Unmarshal([]byte(body), &bodyJSON)
  if err != nil {
    t.Fatal(err.Error())
  }

  if "fail" != bodyJSON.Status {
    t.Fatal("Expect fail was", bodyJSON.Status)
  }

  errMsg := "json: cannot unmarshal string into Go value of type int"
  if errMsg != bodyJSON.Msg {
    t.Fatal("Expect", errMsg, "was", bodyJSON.Msg)
  }
}
func Test_DocNumberProposalNextHandler_OK(t *testing.T) {
  globals := MakeTestGlobals(t)
  session := globals.MongoDB.Session.Copy()
  defer session.Close()
  db := session.DB(TestDBName)
  defer db.DropDatabase()

  docNumberProposal := DocNumberProposal(123)
  err := docNumberProposal.Save(db)
  if err != nil {
    t.Fatal(err.Error())
  }

  handler := gin.New()
  handler.GET(
    "/",
    MakeGlobalsHandler(DocNumberProposalNextHandler, globals),
  )
  header := http.Header{}

  request := TestRequest{
    Body: "",
    Handler: handler,
    Header: header,
  }
  response := request.Send("GET", "/")

  body := response.Body.String()
  bodyJSON := DocNumberProposalNextResponse{}
  err = json.Unmarshal([]byte(body), &bodyJSON)
  if err != nil {
    t.Fatal(err.Error())
  }

  if "success" != bodyJSON.Status {
    t.Fatal("Expect success was", bodyJSON.Status)
  }

  expectProposal := 124
  if expectProposal != bodyJSON.Proposal {
    t.Fatal("Expect", expectProposal, "was", bodyJSON.Proposal)
  }
}


