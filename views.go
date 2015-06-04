package bebber

import (
  "fmt"
  "time"
  "bytes"
  "errors"
  "strconv"
  "net/http"
  "math/rand"
  "crypto/sha1"
  "io/ioutil"
  "encoding/json"
  "github.com/gin-gonic/gin"
  "gopkg.in/mgo.v2"
  "gopkg.in/mgo.v2/bson"
)

const (
  DbFileCollection = "files"
  UsersCollection = "users"
)

type ErrorResponse struct {
  Status string
  Msg string
}

type SuccessResponse struct {
  Status string
  Msg string
}

type LoadDirRequest struct {
  Dir string
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

type LoginData struct {
  Username string
  Password string
}

type UserSession struct {
  Token string
  User string
  Expires time.Time
}

type User struct {
  Username string
  Password string
  Dirs map[string]string
}

func LoadDir(c *gin.Context) {
  session, err := mgo.Dial(GetSettings("BEBBER_DB_SERVER"))
  if err != nil {
    errMsg := "Db error - "+ err.Error()
    c.JSON(http.StatusOK, ErrorResponse{"fail", errMsg})
    return
  }
  defer session.Close()

  collection := session.DB(GetSettings("BEBBER_DB_NAME")).C(DbFileCollection)


  buf := new(bytes.Buffer)
  buf.ReadFrom(c.Request.Body)


  var dir LoadDirRequest
  err = json.Unmarshal(buf.Bytes(), &dir)
  if err != nil {
    errMsg := "Parsing error - "+ err.Error()
    c.JSON(http.StatusOK, ErrorResponse{"fail", errMsg})
    return
  }

  files, err := ioutil.ReadDir(dir.Dir)
  if err != nil {
    errMsg := "Dir error - "+ err.Error()
    c.JSON(http.StatusOK, ErrorResponse{"fail", errMsg})
    return
  }
  incFiles := []string{}
  for i := range files {
    if files[i].IsDir() == false {
      incFiles = append(incFiles, files[i].Name())
    }
  }

  filter := bson.M{
    "filename": bson.M{"$in": incFiles},
  }

  var result []FileDoc
  iter := collection.Find(filter).Iter()
  err = iter.All(&result)
  if err != nil {
    errMsg := "Db Error - "+ err.Error()
    c.JSON(http.StatusOK, ErrorResponse{"fail", errMsg})
    return
  }

  tmp := make([]string, len(result))
  for i := range result {
    tmp[i] = result[i].Filename
  }

  sub := SubList(incFiles, tmp)

  for i := range sub {
    e := FileDoc{
      sub[i],
      []SimpleTag{},
      []RangeTag{},
      []ValueTag{},
    }
    result = append(result, e)
  }

  res := LoadDirResponse{"success", result}

  c.JSON(http.StatusOK, res)

}

func AddTags(c *gin.Context) {
  jsonReq := AddTagsRequest{}
  err := ParseJsonRequest(c, &jsonReq)
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
  collection := session.DB(GetSettings("BEBBER_DB_NAME")).C(DbFileCollection)

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

  c.JSON(http.StatusOK, SuccessResponse{"success", ""})
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

  collection := session.DB(GetSettings("BEBBER_DB_NAME")).C(DbFileCollection)

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
  err := ParseJsonRequest(c, &loginData)
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
  users := usersC.Find(bson.M{"Username": loginData.Username,
                     "Password": sha1Pass})
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

  cookie := http.Cookie{Name: "X-XSRF-TOKEN", Value: token, Expires: expires}
  http.SetCookie(c.Writer, &cookie)
  c.JSON(http.StatusOK, SuccessResponse{Status: "success"})
}
