package bebber

import (
  "encoding/json"
  "gopkg.in/mgo.v2"
  "gopkg.in/mgo.v2/bson"
)

func Search(str string, db *mgo.Database) *[]bson.M {
  searchQuery := bson.M{}
  json.Unmarshal([]byte(str), &searchQuery)

  filesColl := db.C(FilesCollection)

  result := []bson.M{}
  query := filesColl.Find(searchQuery)
  query.All(&result)
  return &result
}
