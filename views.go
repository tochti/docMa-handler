package bebber

import (
  "time"
  "bytes"
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

type SimpleTag struct {
  Tag string
}

type RangeTag struct {
  Tag string
  Start time.Time
  End time.Time
}

type ValueTag struct {
  Tag string
  Value string
}

type FileDoc struct {
  Filename string
  SimpleTags []SimpleTag
  RangeTags []RangeTag
  ValueTags []ValueTag
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

