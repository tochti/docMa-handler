package bebber

import (
  _"time"
  "bytes"
  "errors"
  _"strings"
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
  status string
  msg string
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
    panic(err)
  }
  defer session.Close()

  collection := session.DB(GetSettings("BEBBER_DB_NAME")).C(DbFileCollection)


  buf := new(bytes.Buffer)
  buf.ReadFrom(c.Request.Body)


  var dir LoadDirRequest
  err = json.Unmarshal(buf.Bytes(), &dir)
  if err != nil {
    c.JSON(http.StatusOK, ErrorResponse{"fail", "dir param is missing"})
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
    c.JSON(http.StatusOK, ErrorResponse{"fail", "DB Problems"})
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
  var jsonReq AddTagsRequest
  err := ParseJsonRequest(c, jsonReq)
  if err != nil {
    c.JSON(http.StatusOK, ErrorResponse{
                          "fail",
                          "Couldn't parse request - "+ err.Error(),
                        })
  }

  session, err := mgo.Dial(GetSettings("BEBBER_DB_SERVER"))
  collection := session.DB(GetSettings("BEBBER_DB_NAME")).C(DbFileCollection)

  doc := FileDoc{}
  err = collection.Find(bson.M{"filename": jsonReq.Filename}).One(&doc)
  if err != nil {
    c.JSON(http.StatusOK, ErrorResponse{
                              "fail",
                              "Db error - "+ err.Error(),
                          })
  }

  _, err = CreateUpdateDoc(jsonReq.Filename, jsonReq.Tags)
}

func CreateUpdateDoc(file string, tags []string) (*FileDoc, error) {
  if len(tags) == 0 {
    return &FileDoc{Filename: file}, nil
  }

  var sTags []SimpleTag
  var vTags []ValueTag
  var rTags []RangeTag
  for _, tag := range tags {
    typ, err := SpotTagType(tag)
    if err != nil {
      return nil, errors.New(err.Error())
    }
    switch typ {
    case "SimpleTag":
      sTags = append(sTags, SimpleTag{tag})
    case "ValueTag":
      vTags = append(vTags, ParseValueTag(tag))
    case "RangeTag":
      rTags = append(rTags, ParseRangeTag(tag))
    }
  }

  doc := FileDoc{
    Filename: file,
    SimpleTags: sTags,
    ValueTags: vTags,
    RangeTags: rTags,
  }

  return &doc, nil
}
