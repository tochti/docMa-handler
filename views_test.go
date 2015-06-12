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

func TestLoadBoxRouteOk(t *testing.T) {
  /* Config */
  os.Setenv("BEBBER_DB_NAME", "bebber_test")
  os.Setenv("BEBBER_DB_SERVER", "127.0.0.1")
  os.Setenv("BEBBER_FILES", testDir)

  m := int(0777)
  mode := os.FileMode(m)
  tmpDir, err := ioutil.TempDir(testDir, "loaddir")
  tmpDirBase := path.Base(tmpDir)
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
  defer session.Close()
  db := session.DB("bebber_test")
  defer db.DropDatabase()

  sT := time.Date(2014, time.April, 1, 0, 0, 0, 0, time.UTC)
  eT := time.Date(2014, time.April, 2, 0, 0, 0, 0, time.UTC)
  doc1 := FileDoc{
    "test1.txt",
    []SimpleTag{SimpleTag{"sTag1"}},
    []RangeTag{RangeTag{"rTag1", sT, eT}},
    []ValueTag{ValueTag{"vTag1", "value1"}},
  }
  doc2 := FileDoc{
    "test2.txt",
    []SimpleTag{SimpleTag{"sTag1"}},
    []RangeTag{},
    []ValueTag{},
  }
  doc3 := FileDoc{
    "notinlist.txt",
    []SimpleTag{},
    []RangeTag{},
    []ValueTag{},
  }

  c := db.C(FilesCollection)
  err = c.Insert(doc1, doc2, doc3)
  if err != nil {
    t.Error(err)
  }

  CreateTestUserSession("bommel", "123", db, t)

  user := User{Username:"bommel",
              Password:"",
              Dirs: map[string]string{"testbox": tmpDirBase}}
  err = user.Save(db.C(UsersCollection))
  if err != nil {
    t.Fatal(err.Error())
  }

  /* Perform Request */
	handler := gin.New()
  handler.GET("/:boxname", Auth(), LoadBox)
  req := TestRequest{
    Handler: handler,
    Body: "",
    Header: http.Header{},
  }
  w := req.DialToken("GET", "/testbox", "123")

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

func TestUserAuthOK(t *testing.T) {
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

func TestUserAuthFail(t *testing.T) {
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

func TestGetUser(t *testing.T) {
  os.Setenv("BEBBER_DB_SERVER", "127.0.0.1")
  os.Setenv("BEBBER_DB_NAME", "bebber_test")
  session, err := mgo.Dial("127.0.0.1")
  if err != nil {
    t.Fatal(err)
  }
  defer session.Close()
  sha1Pass := fmt.Sprintf("%x", sha1.Sum([]byte("test")))
  dirs := map[string]string{"inbox": "/to/dir"}
  user := User{Username: "hitman", Password: sha1Pass, Dirs:dirs}
  user1 := User{Username: "catwomen", Password: sha1Pass, Dirs:dirs}
  user2 := User{Username: "loveMaster_999", Password: sha1Pass, Dirs:dirs}
  usersC := session.DB("bebber_test").C("users")
  usersC.Insert(user, user1, user2)
  defer session.DB("bebber_test").DropDatabase()

  r := gin.New()
  r.GET("/:name", GetUser)
  body := bytes.NewBufferString("")
  resp := PerformRequest(r, "GET", "/hitman", body)

  expect := `{"Username":"hitman","Password":"","Dirs":{"inbox":"/to/dir"}}`+"\n"
  if resp.Body.String() != expect {
    t.Fatal("Expect", expect, "was", resp.Body.String())
  }

}

func TestMoveFileOk(t *testing.T) {
  /* setup Test */
  fromBox, err := ioutil.TempDir(testDir, "box")
  if err != nil {
    t.Fatal(err.Error())
  }
  defer os.RemoveAll(fromBox)

  toBox, err := ioutil.TempDir(testDir, "box")
  if err != nil  {
    t.Fatal(err.Error())
  }
  defer os.RemoveAll(toBox)

  moveFile, err := ioutil.TempFile(fromBox, "moveme")
  if err != nil {
    t.Fatal(err.Error())
  }
  moveFile.Close()

  session, err := mgo.Dial("127.0.0.1")
  if err != nil {
    t.Fatal(err.Error())
  }
  defer session.Close()
  defer session.DB("bebber_test").DropDatabase()

  sessionsC := session.DB("bebber_test").C(SessionsCollection)
  expires := time.Now().AddDate(0,0,1)
  userSession := UserSession{Token: "123", User: "kloß", Expires: expires}
  err = sessionsC.Insert(userSession)
  if err != nil {
    t.Fatal(err.Error())
  }

  usersC := session.DB("bebber_test").C(UsersCollection)
  dirs := map[string]string{"from":fromBox,"to":toBox}
  user := User{Username:"kloß",Password:"",Dirs:dirs}
  err = usersC.Insert(user)
  if err != nil {
    t.Fatal(err.Error())
  }

  /* Perform Request */
  serv := gin.New()
  serv.POST("/", Auth(), MoveFile)
  bodyStr := `{"FromBox":"from","ToBox":"to","File":"`+path.Base(moveFile.Name())+`"}`
  body := bytes.NewBufferString(bodyStr)
  header := http.Header{}
  header.Add("X-XSRF-TOKEN", "123")
  resp := PerformRequestHeader(serv, "POST", "/", body, &header)

  if resp.Code != 200 {
    t.Fatal("Expect 200 was", resp.Code)
  }

  err = moveFile.Sync()
  if os.IsExist(err) == true {
    t.Fatal(moveFile.Name, " should be moved", resp.Body.String())
  }

  _, err = os.Open(path.Join(toBox, path.Base(moveFile.Name())))
  if err != nil {
    t.Fatal(err.Error(), resp.Body.String())
  }

}
