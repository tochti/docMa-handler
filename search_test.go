package bebber

import (
  "strings"
  "testing"
  "encoding/json"
  "gopkg.in/mgo.v2"
  "gopkg.in/mgo.v2/bson"
)

func Test_Search(t *testing.T) {
  doc := bson.M{"Song": bson.M{"RingOf": "Fire"}}
  session, err := mgo.Dial(TestDBServer)
  if err != nil {
    t.Fatal(err.Error())
  }
  defer session.Close()

  db := session.DB(TestDBName)
  defer db.DropDatabase()

  filesColl := db.C("files")
  err = filesColl.Insert(doc)
  if err != nil {
    t.Fatal(err.Error())
  }

  searchStr := `{"Song.RingOf": "Fire"}`
  tmp := Search(searchStr, db)
  result := *tmp

  if len(result) != 1 {
    t.Fatal("Expect len 1 was", len(result))
  }

  JSONBytes, err := json.Marshal(result)
  if err != nil {
    t.Fatal(err.Error())
  }

  expectJSONStr := `{"Song":{"RingOf":"Fire"}`
  if strings.Contains(expectJSONStr, string(JSONBytes)) {
    t.Fatal("Expect to contain", expectJSONStr, "was", string(JSONBytes))
  }
}
