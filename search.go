package bebber

import (
  "time"
  "errors"
  "reflect"
  "strings"
  "encoding/json"
  "gopkg.in/mgo.v2"
  "gopkg.in/mgo.v2/bson"
)

var (
  typeOfString = reflect.TypeOf("")
  typeOfDateOfScan = reflect.TypeOf(map[string]interface{}{})
)

func SearchDocs(str string, db *mgo.Database) ([]Doc, error) {
  searchQuery, err := MakeSearchQuery(str)
  if err != nil {
    return []Doc{}, err
  }

  docsColl := db.C(DocsColl)

  result := []Doc{}
  query := docsColl.Find(searchQuery)
  err = query.All(&result)
  if err != nil {
    return []Doc{}, err
  }
  return result, nil
}

func MakeSearchQuery(str string) (bson.M, error) {
  searchQueryTmp := bson.M{}
  err := json.Unmarshal([]byte(str), &searchQueryTmp)
  if err != nil {
    return bson.M{}, err
  }

  searchQueryNew := bson.M{}
  labelsStr, ok := searchQueryTmp["Labels"]
  if ok && reflect.TypeOf(labelsStr) == typeOfString {
    labels := strings.Split(labelsStr.(string), ",")
    searchQueryNew["labels"] = bson.M{"$in": labels}
  }

  docNumbersStr, ok := searchQueryTmp["DocNumbers"]
  if ok && reflect.TypeOf(docNumbersStr) == typeOfString {
    docNumbers := strings.Split(docNumbersStr.(string), ",")
    searchQueryNew["accountdata.docnumbers"] = bson.M{"$in": docNumbers}
  }

  dateOfScanTmp, ok := searchQueryTmp["DateOfScan"]
  if ok && reflect.TypeOf(dateOfScanTmp) == typeOfDateOfScan {
    dateOfScan := dateOfScanTmp.(map[string]interface{})
    tmp := map[string]time.Time{}
    from, ok := dateOfScan["From"]
    if ok && reflect.TypeOf(from) == typeOfString {
      d, err := time.Parse(time.RFC3339, from.(string))
      if err != nil {
        return bson.M{}, err
      }
      tmp["$gte"] = d
    }
    to, ok := dateOfScan["To"]
    if ok && reflect.TypeOf(to) == typeOfString {
      d, err := time.Parse(time.RFC3339, to.(string))
      if err != nil {
        return bson.M{}, err
      }
      tmp["$lte"] = d
    }
    if len(tmp) != 0 {
      searchQueryNew["infos.dateofscan"] = tmp
    }
  }

  if len(searchQueryNew) == 0 {
    return bson.M{}, errors.New("Empty search query")
  }

  return searchQueryNew, nil
}
