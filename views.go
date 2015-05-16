package bebber

import (
  "bytes"
  "errors"
  "net/http"
  "io/ioutil"
  "encoding/json"
  "github.com/gin-gonic/gin"
  "gopkg.in/mgo.v2"
  "gopkg.in/mgo.v2/bson"
)

const (
  DbFileCollection = "files"
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
  incFiles := make([]string, len(files))
  for i := range files {
    if files[i].IsDir() == false {
      incFiles[i] = files[i].Name()
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

  session, err := mgo.Dial(GetSettings("BEBBER_DB_SERVER"))
  collection := session.DB(GetSettings("BEBBER_DB_NAME")).C(DbFileCollection)

  updateDoc := FileDoc{}
  err = collection.Find(bson.M{"filename": jsonReq.Filename}).One(&updateDoc)
  if err != nil {
    errMsg := "Db error -"+ err.Error()
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
  if info.Updated != 1 {
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
      doc.RangeTags = append(doc.RangeTags, ParseRangeTag(tag))
    }
  }

  return nil
}
