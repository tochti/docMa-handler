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
  session, err := mgo.Dial(TestDBHost)
  if err != nil {
    t.Fatal(err.Error())
  }
  defer session.Close()

  db := session.DB(TestDBName)
  defer db.DropDatabase()

  docsColl := db.C(DocsColl)
  err = docsColl.Insert(doc)
  if err != nil {
    t.Fatal(err.Error())
  }

  searchStr := `{"Song.RingOf": "Fire"}`
  result := Search(searchStr, db)

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
