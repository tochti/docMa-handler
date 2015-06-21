package bebber

import (
  "os"
  "fmt"
  "time"
  "path"
  "bytes"
  "errors"
  "strings"
  "strconv"
  "net/http"
  "math/rand"
  "crypto/sha1"

  "gopkg.in/mgo.v2"
  "gopkg.in/mgo.v2/bson"
  "github.com/gin-gonic/gin"
)

const (
  FilesCollection = "files"
)

type ErrorResponse struct {
  Status string
  Msg string
}

type LoadDirResponse struct {
  Status string
  Dir []FileDoc
}

type AddTagsRequest struct {
  Filename string
  Tags []string
}

type LoadAccFilesResponse struct {
  Status string
  Msg string
  AccFiles []AccFile
}

func LoadFile(c *gin.Context) {
  filename := strings.Trim(c.Params.ByName("filename"), "\"")

  session, err := mgo.Dial(GetSettings("BEBBER_DB_SERVER"))
  if err != nil {
    c.JSON(http.StatusOK, ErrorResponse{"fail", err.Error()})
  }
  defer session.Close()
  db := session.DB(GetSettings("BEBBER_DB_NAME"))

  usTmp, err := c.Get("session")
  userSession := usTmp.(UserSession)
  if err != nil {
    c.JSON(http.StatusOK, ErrorResponse{"fail", err.Error()})
    return
  }
  user := User{}
  err = user.Load(userSession.User, db.C(UsersCollection))
  if err != nil {
    c.JSON(http.StatusOK, ErrorResponse{"fail", err.Error()})
  }

  boxpath := GetSettings("BEBBER_FILES")
  filepath := path.Join(boxpath, filename)
  c.File(filepath)
}

func AddTags(c *gin.Context) {
  jsonReq := AddTagsRequest{}
  err := ParseJSONRequest(c, &jsonReq)
  if err != nil {
    errMsg := "Couldn't parse request - "+ err.Error()
    res := ErrorResponse{"fail", errMsg}
    c.JSON(http.StatusOK, res)
    return
  }

  if jsonReq.Filename == "" {
    res := ErrorResponse{"fail", "No Filename"}
    c.JSON(http.StatusOK, res)
    return
  }

  session, err := mgo.Dial(GetSettings("BEBBER_DB_SERVER"))
  collection := session.DB(GetSettings("BEBBER_DB_NAME")).C(FilesCollection)

  updateDoc := FileDoc{Filename: jsonReq.Filename}
  err = collection.Find(bson.M{"filename": jsonReq.Filename}).One(&updateDoc)
  if err != nil  && err.Error() != "not found" {
    errMsg := "Db error - "+ err.Error()
    c.JSON(http.StatusOK, ErrorResponse{"fail", errMsg})
    return
  }

  err = CreateUpdateDoc(jsonReq.Tags, &updateDoc)
  if err != nil {
    errMsg := "Cannot update file "+ jsonReq.Filename +" - "+ err.Error()
    c.JSON(http.StatusOK, ErrorResponse{"fail", errMsg})
    return
  }
  newDoc := FileDoc{}
  change := mgo.Change{
              Update: updateDoc,
              Upsert: true,
              ReturnNew: true,
            }
  info, err := collection.Find(bson.M{"filename": jsonReq.Filename}).
                          Apply(change, &newDoc)
  if err != nil {
    errMsg := "Cannot update file "+ jsonReq.Filename +" - "+ err.Error()
    c.JSON(http.StatusOK, ErrorResponse{"fail", errMsg})
    return
  }

  if info.Updated != 1 && info.UpsertedId == nil {
    errMsg := "Expected to update 1 document, was "+ string(info.Updated)
    c.JSON(http.StatusOK, ErrorResponse{"fail", errMsg})
    return
  }

  c.JSON(http.StatusOK, SuccessResponse{})
}

func CreateUpdateDoc(tags []string, doc *FileDoc) error {
  for _, tag := range tags {
    typ, err := SpotTagType(tag)
    if err != nil {
      return errors.New(err.Error())
    }
    switch typ {
    case "SimpleTag":
      doc.SimpleTags = append(doc.SimpleTags, SimpleTag{tag})
    case "ValueTag":
      doc.ValueTags = append(doc.ValueTags, ParseValueTag(tag))
    case "RangeTag":
      tags, err := ParseRangeTag(tag)
      if err != nil {
        return err
      }
      doc.RangeTags = append(doc.RangeTags, *tags)
    }
  }

  return nil
}

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

func Login(c *gin.Context) {
  loginData := LoginData{}
  err := ParseJSONRequest(c, &loginData)
  if err != nil {
    c.JSON(http.StatusOK, ErrorResponse{"fail", err.Error()})
    return
  }

  session, err := mgo.Dial(GetSettings("BEBBER_DB_SERVER"))
  if err != nil {
    c.JSON(http.StatusOK, ErrorResponse{"fail", err.Error()})
    return
  }
  db := session.DB(GetSettings("BEBBER_DB_NAME"))

  seed := rand.New(rand.NewSource(time.Now().UnixNano()))
  sha1Pass := fmt.Sprintf("%x", sha1.Sum([]byte(loginData.Password)))

  usersC := db.C(UsersCollection)
  users := usersC.Find(bson.M{"username": loginData.Username,
                     "password": sha1Pass})
  n, err := users.Count()
  if err != nil {
    c.JSON(http.StatusOK, ErrorResponse{"fail", err.Error()})
    return
  }
  if n != 1 {
    c.JSON(http.StatusOK, ErrorResponse{"fail", "Wrong username or password"})
    return
  }

  tmp := strconv.Itoa(seed.Int())
  token := fmt.Sprintf("%x", sha1.Sum([]byte(tmp)))
  expires := time.Now().AddDate(0,0,2)

  sessionsC := db.C(SessionsCollection)
  userSession := UserSession{User: loginData.Username,
                             Token: token, Expires: expires}
  err = sessionsC.Insert(userSession)
  if err != nil {
    c.JSON(http.StatusOK, ErrorResponse{"fail", err.Error()})
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
  usersColl := db.C(UsersCollection)
  users := usersColl.Find(bson.M{"username": username})
  n, err := users.Count()
  if err != nil {
    c.JSON(http.StatusOK, ErrorResponse{"fail", err.Error()})
    return
  }
  if n != 1 {
    //TODO: Info an Admin das es user mit dem selben Username mehrmals gibt
    c.JSON(http.StatusOK, ErrorResponse{"fail", "Cannot find user"})
    return
  }

  user := User{}
  err = users.One(&user)
  if err != nil {
    c.JSON(http.StatusOK, ErrorResponse{"fail", err.Error()})
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
  err = doc.Find(db)
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
    MakeFailResponse(c, err.Error())
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

func DocRemoveLabelsHandler(c *gin.Context, g Globals) {
  appendRequest := DocRemoveLabelsRequest{}
  err := ParseJSONRequest(c, &appendRequest)
  if err != nil {
    MakeFailResponse(c, err.Error())
    return
  }

  session := g.MongoDB.Session.Copy()
  defer session.Close()

  db := session.DB(g.Config["MONGODB_DBNAME"])

  doc := Doc{Name: appendRequest.Name}
  err = doc.RemoveLabels(appendRequest.Labels, db)
  if err != nil {
    MakeFailResponse(c, err.Error())
    return
  }

  c.JSON(http.StatusOK, SuccessResponse{"success"})
}
