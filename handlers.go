package bebber

import (
  "os"
  "fmt"
  "time"
  "path"
  "bytes"
  "strings"
  "strconv"
  "net/http"
  "math/rand"
  "crypto/sha1"

  "gopkg.in/mgo.v2/bson"
  "github.com/gin-gonic/gin"
)


func ReadDocFileHandler(c *gin.Context, g Globals) {
  filename := strings.Trim(c.Params.ByName("name"), "\"")
  filepath := path.Join(g.Config["FILES_DIR"], filename)
  c.File(filepath)
}

/*
func LoadAccFiles(c *gin.Context) {
  session, err := mgo.Dial(GetSettings("BEBBER_DB_SERVER"))
  if err != nil {
    errMsg := "Db error - "+ err.Error()
    c.JSON(http.StatusOK, ErrorResponse{"fail", errMsg})
    return
  }
  defer session.Close()

  collection := session.DB(GetSettings("BEBBER_DB_NAME")).C(FilesCollection)

  accData := []AccData{}
  err = ReadAccFile(GetSettings("BEBBER_ACC_FILE"), &accData)
  if err != nil {
    c.JSON(http.StatusOK, ErrorResponse{"fail", err.Error()})
    return
  }

  validCSV := false
  r, err := c.Get("validCSV")
  if err == nil {
    validCSV = r.(bool)
  }
  accFiles, err := JoinAccFile(accData, collection, validCSV)
  if err != nil {
    c.JSON(http.StatusOK, ErrorResponse{"fail", err.Error()})
    return
  }

  res := LoadAccFilesResponse{
          Status: "success",
          AccFiles: accFiles,
         }

  c.JSON(http.StatusOK, res)
}
*/

func LoginHandler(c *gin.Context, g Globals) {
  loginData := LoginData{}
  err := ParseJSONRequest(c, &loginData)
  if err != nil {
    MakeFailResponse(c, err.Error())
    return
  }
  session := g.MongoDB.Session.Copy()
  defer session.Close()
  db := session.DB(g.Config["MONGODB_DBNAME"])

  seed := rand.New(rand.NewSource(time.Now().UnixNano()))
  sha1Pass := fmt.Sprintf("%x", sha1.Sum([]byte(loginData.Password)))

  usersColl := db.C(UsersColl)
  users := usersColl.Find(bson.M{"username": loginData.Username,
                     "password": sha1Pass})
  n, err := users.Count()
  if err != nil {
    MakeFailResponse(c, err.Error())
    return
  }
  if n != 1 {
    MakeFailResponse(c, "Wrong username or password")
    return
  }

  tmp := strconv.Itoa(seed.Int())
  token := fmt.Sprintf("%x", sha1.Sum([]byte(tmp)))
  expires := time.Now().AddDate(0,0,2)

  sessionsColl := db.C(SessionsColl)
  userSession := UserSession{User: loginData.Username,
                             Token: token, Expires: expires}
  err = sessionsColl.Insert(userSession)
  if err != nil {
    MakeFailResponse(c, err.Error())
    return
  }

  cookie := http.Cookie{Name: XSRFCookieName, Value: token, Expires: expires}
  http.SetCookie(c.Writer, &cookie)
  c.JSON(http.StatusOK, SuccessResponse{"success"})
}

func UserHandler(c *gin.Context, globals Globals) {
  config := globals.Config

  session := globals.MongoDB.Session.Copy()
  defer session.Close()

  db := session.DB(config["MONGODB_DBNAME"])
  username := c.Params.ByName("name")
  usersColl := db.C(UsersColl)
  users := usersColl.Find(bson.M{"username": username})
  n, err := users.Count()
  if err != nil {
    MakeFailResponse(c, err.Error())
    return
  }
  if n != 1 {
    MakeFailResponse(c, "Cannot find user")
    return
  }

  user := User{}
  err = users.One(&user)
  if err != nil {
    MakeFailResponse(c, err.Error())
    return
  }

  user.Password = ""
  c.JSON(http.StatusOK, user)
}

func SearchHandler(c *gin.Context, g Globals) {
  session := g.MongoDB.Session.Copy()
  defer session.Close()
  db := session.DB(g.Config["MONGODB_DBNAME"])

  buf := new(bytes.Buffer)
  buf.ReadFrom(c.Request.Body)
  body := buf.String()
  result := Search(body, db)
  c.JSON(http.StatusOK, result)
}

func DocMakeHandler(c *gin.Context, g Globals) {
  session := g.MongoDB.Session.Copy()
  defer session.Close()

  db := session.DB(g.Config["MONGODB_DBNAME"])
  docsColl := db.C(DocsColl)

  requestBody := DocMakeRequest{}
  err := ParseJSONRequest(c, &requestBody)
  if err != nil {
    MakeFailResponse(c, err.Error())
    return
  }

  if requestBody.Name == "" {
    MakeFailResponse(c, "Missing a name!")
    return
  }

  if requestBody.Infos.DateOfScan.IsZero() ||
     requestBody.Infos.DateOfReceipt.IsZero() {
    MakeFailResponse(c, "Missing the infos field!")
    return
  }

  doc := Doc{Name: requestBody.Name}
  err = doc.Find(db)
  if err == nil {
    MakeFailResponse(c, "Document already exists!")
    return
  }

  docID := bson.NewObjectId()
  requestBody.ID = docID
  err = docsColl.Insert(requestBody)
  if err != nil {
    MakeFailResponse(c, err.Error())
    return
  }

  c.JSON(http.StatusOK,
         MongoDBSuccessResponse{Status: "success", DocID: docID.Hex()})
}

func DocChangeHandler(c *gin.Context, g Globals) {
  changeRequest := DocChangeRequest{}
  ParseJSONRequest(c, &changeRequest)

  session := g.MongoDB.Session.Copy()
  defer session.Close()

  db := session.DB(g.Config["MONGODB_DBNAME"])

  doc := Doc{Name: changeRequest.Name}
  err := doc.Find(db)
  if err != nil {
    MakeFailResponse(c, err.Error())
    return
  }

  changeDoc := Doc(changeRequest)
  err = doc.Change(changeDoc, db)
  if err != nil {
    MakeFailResponse(c, err.Error())
    return
  }

  c.JSON(http.StatusOK, SuccessResponse{"success"})
}

func DocReadHandler(c *gin.Context, g Globals) {
  name := c.Params.ByName("name")
  session := g.MongoDB.Session.Copy()
  defer session.Close()

  db := session.DB(g.Config["MONGODB_DBNAME"])
  doc := Doc{Name: name}
  err := doc.Find(db)
  if err != nil {
    MakeFailResponse(c, err.Error())
    return
  }

  c.JSON(http.StatusOK, DocReadResponse{Status: "success", Doc: doc})
}

func DocRemoveHandler(c *gin.Context, g Globals) {
  name := c.Params.ByName("name")
  session := g.MongoDB.Session.Copy()
  defer session.Close()

  db := session.DB(g.Config["MONGODB_DBNAME"])
  doc := Doc{Name: name}
  err := doc.Remove(db)
  if err != nil {
    MakeFailResponse(c, err.Error())
    return
  }

  err = os.Remove(path.Join(g.Config["FILES_DIR"], name))
  if err != nil {
type DocAppendLabelsRequest struct {
  Name string
  Labels []Label
}
    MakeFailResponse(c, err.Error())
  }

  c.JSON(http.StatusOK, SuccessResponse{"success"})
}

func DocRenameHandler(c *gin.Context, g Globals) {
  renameRequest := DocRenameRequest{}
  err := ParseJSONRequest(c, &renameRequest)
  if err != nil {
    MakeFailResponse(c, err.Error())
    return
  }

  session := g.MongoDB.Session.Copy()
  defer session.Close()

  db := session.DB(g.Config["MONGODB_DBNAME"])

  doc := Doc{Name: renameRequest.Name}
  changeDoc := Doc{Name: renameRequest.NewName}
  err = doc.Change(changeDoc, db)
  if err != nil {
    MakeFailResponse(c, err.Error())
    return
  }

  docs := g.Config["FILES_DIR"]
  oldPath := path.Join(docs, renameRequest.Name)
  newPath := path.Join(docs, renameRequest.NewName)
  err = os.Rename(oldPath, newPath)
  if err != nil {
    MakeFailResponse(c, err.Error())
    return
  }

  c.JSON(http.StatusOK, SuccessResponse{"success"})
}

func DocAppendLabelsHandler(c *gin.Context, g Globals) {
  appendRequest := DocAppendLabelsRequest{}
  err := ParseJSONRequest(c, &appendRequest)
  if err != nil {
    MakeFailResponse(c, err.Error())
    return
  }

  session := g.MongoDB.Session.Copy()
  defer session.Close()

  db := session.DB(g.Config["MONGODB_DBNAME"])

  doc := Doc{Name: appendRequest.Name}
  err = doc.AppendLabels(appendRequest.Labels, db)
  if err != nil {
    MakeFailResponse(c, err.Error())
    return
  }

  c.JSON(http.StatusOK, SuccessResponse{"success"})
}

func DocRemoveLabelHandler(c *gin.Context, g Globals) {
  name := c.Params.ByName("name")
  label := c.Params.ByName("label")

  session := g.MongoDB.Session.Copy()
  defer session.Close()

  db := session.DB(g.Config["MONGODB_DBNAME"])

  doc := Doc{Name: name}
  labels := []Label{Label(label)}
  err := doc.RemoveLabels(labels, db)
  if err != nil {
    MakeFailResponse(c, err.Error())
    return
  }

  c.JSON(http.StatusOK, SuccessResponse{"success"})
}

func DocAppendDocNumbersHandler(c *gin.Context, g Globals) {
  appendRequest := DocAppendDocNumbersRequest{}
  err := ParseJSONRequest(c, &appendRequest)
  if err != nil {
    MakeFailResponse(c, err.Error())
    return
  }

  session := g.MongoDB.Session.Copy()
  defer session.Close()

  db := session.DB(g.Config["MONGODB_DBNAME"])

  doc := Doc{Name: appendRequest.Name}
  err = doc.AppendDocNumbers(appendRequest.DocNumbers, db)
  if err != nil {
    MakeFailResponse(c, err.Error())
    return
  }

  c.JSON(http.StatusOK, SuccessResponse{"success"})
}

func DocRemoveDocNumberHandler(c *gin.Context, g Globals) {
  name := c.Params.ByName("name")
  number := c.Params.ByName("number")

  session := g.MongoDB.Session.Copy()
  defer session.Close()

  db := session.DB(g.Config["MONGODB_DBNAME"])

  doc := Doc{Name: name}
  docNumbers := []string{number}
  err := doc.RemoveDocNumbers(docNumbers, db)
  if err != nil {
    MakeFailResponse(c, err.Error())
    return
  }

  c.JSON(http.StatusOK, SuccessResponse{"success"})
}

func AccProcessMakeHandler(c *gin.Context, g Globals) {
  requestBody := AccProcessMakeRequest{}
  err := ParseJSONRequest(c, &requestBody)
  if err != nil {
    MakeFailResponse(c, err.Error())
    return
  }

  session := g.MongoDB.Session.Copy()
  defer session.Close()

  db := session.DB(g.Config["MONGODB_DBNAME"])
  accProcessColl := db.C(AccProcessColl)

  docID := bson.NewObjectId()
  requestBody.ID = docID
  err = accProcessColl.Insert(requestBody)
  if err != nil {
    MakeFailResponse(c, err.Error())
    return
  }

  response := AccProcessMakeResponse{"success", docID.Hex()}
  c.JSON(http.StatusOK, response)
}

func AccProcessFindByDocNumberHandler(c *gin.Context, g Globals) {
  number := c.Params.ByName("number")
  session := g.MongoDB.Session.Copy()
  defer session.Close()

  db := session.DB(g.Config["MONGODB_DBNAME"])

  accProcessList, err := FindAccProcessByDocNumbers(db, []string{number})
  if err != nil {
    MakeFailResponse(c, err.Error())
    return
  }

  response := AccProcessFindByDocNumberResponse{
    "success",
    accProcessList,
  }
  c.JSON(http.StatusOK, response)
}

func AccProcessFindByAccNumberHandler(c *gin.Context, g Globals) {
  from := c.Params.ByName("from")
  fromDate, err := ParseDate(from)
  if err != nil {
    MakeFailResponse(c, err.Error())
    return
  }

  to := c.Params.ByName("to")
  toDate, err := ParseDate(to)
  if err != nil {
    MakeFailResponse(c, err.Error())
    return
  }

  number := c.Params.ByName("number")
  accNumber, err := strconv.Atoi(number)
  if err != nil {
    MakeFailResponse(c, err.Error())
    return
  }

  session := g.MongoDB.Session.Copy()
  defer session.Close()

  db := session.DB(g.Config["MONGODB_DBNAME"])

  accProList, err := FindAccProcessByAccNumber(db, accNumber, fromDate, toDate)
  if err != nil {
    MakeFailResponse(c, err.Error())
    return
  }

  response := AccProcessFindByDocNumberResponse{
    "success",
    accProList,
  }
  c.JSON(http.StatusOK, response)
}
