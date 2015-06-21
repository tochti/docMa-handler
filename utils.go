package bebber

import (
  "io"
  "os"
  "fmt"
  "time"
  "bytes"
  "errors"
  "strconv"
  "strings"
  "net/http"
  "crypto/sha1"
  "encoding/csv"
  "encoding/json"

  "gopkg.in/mgo.v2"
  "gopkg.in/mgo.v2/bson"
  "github.com/gin-gonic/gin"
)

const (
  RangeSep = ".."
  TagKeyValueSep = ":"
)

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

type FileDocs struct {
  List []FileDoc
}

type AccData struct {
  Belegdatum time.Time
  Buchungsdatum time.Time
  Belegnummernkreis string
  Belegnummer string
  Buchungstext string
  Buchungsbetrag float64
  Sollkonto int64
  Habenkonto int64
  Steuerschlüssel int64
  Kostenstelle1 string
  Kostenstelle2 string
  BuchungsbetragEuro float64
  Waehrung string
}

type AccFile struct {
  AccData AccData
  FileDoc FileDoc
}

func Month(m int64) (time.Month, error) {
  if (m < 1) || (m > 12) {
    return time.April, errors.New("Month out of range")
  }
  months := []time.Month{
      time.January,
      time.February,
      time.March,
      time.April,
      time.May,
      time.June,
      time.July,
      time.August,
      time.September,
      time.October,
      time.November,
      time.December,
    }

  return months[m-1], nil
}

func MakeFailResponse(c *gin.Context, msg string) {
  c.JSON(http.StatusOK, FailResponse{Status: "fail", Msg: msg})
}

func GetSettings(k string) string {
  v := os.Getenv(k)
  if v == "" {
    fmt.Println(k + " env is missing")
    os.Exit(2)
  }

  return v
}

func SubList(a, b []string) []string {
  m := make(map[string]bool)
  for _, v := range b {
    m[v] = false
  }

  var res []string
  for _, v := range a {
    _, ok := m[v]
    if ok == false {
      res = append(res, v)
    }
  }

  return res
}

func SpotTagType(tag string) (string, error) {
  if strings.Contains(tag, TagKeyValueSep) == false{
    return "SimpleTag", nil
  }

  if strings.Index(tag, TagKeyValueSep)+1 == len(tag) {
    return "", errors.New("Missing Value")
  }

  tmp := strings.Split(tag, TagKeyValueSep)
  value := strings.Join(tmp[1:], TagKeyValueSep)

  if strings.Contains(value, RangeSep) {
    tmp := strings.SplitAfter(value, RangeSep)
    // If ..EndDate
    if tmp[0] == RangeSep {
      if len(tmp[1]) == 8 {
        return "RangeTag", nil
      } else {
        return "RangeTag", errors.New("Error in range")
      }
    }
    // If StartDate..
    if tmp[1] == "" {
      if len(tmp[0]) == (8 + len(RangeSep)) {
        return "RangeTag", nil
      } else {
        return "RangeTag", errors.New("Error in range")
      }
    }

    sDate := strings.Replace(tmp[0], RangeSep, "", 1)
    eDate := tmp[1]
    if len(sDate) != 8 || len(eDate) != 8 {
      return "RangeTag", errors.New("Error in range")
    }
    return "RangeTag", nil
  } else {
    return "ValueTag", nil
  }

}

func ParseValueTag(tag string) ValueTag {
  tmp := strings.Split(tag, TagKeyValueSep)
  key := tmp[0]
  value := strings.Join(tmp[1:], TagKeyValueSep)
  if value[0] == '"' && value[len(value)-1] == '"' {
    value = value[1:len(value)-1]
  }
  return ValueTag{key, value}
}

func ParseRangeTag(tag string) (*RangeTag, error) {
  tmp := strings.Split(tag, TagKeyValueSep)
  key := tmp[0]
  val := strings.SplitAfter(tmp[1], RangeSep)

  var sDate time.Time
  var eDate time.Time
  var err error
  if val[0] == RangeSep {
    eDate, err = ParseDate(val[1])
    if err != nil {
      return nil, err
    }
  } else if val[1] == "" {
    sDate, err = ParseDate(strings.Replace(val[0], RangeSep, "", 1))
    if err != nil {
      return nil, err
    }
  } else {
    sDate, err = ParseDate(strings.Replace(val[0], RangeSep, "", 1))
    if err != nil {
      return nil, err
    }
    eDate, err = ParseDate(val[1])
    if err != nil {
      return nil, err
    }
  }

  return &RangeTag{Tag: key, Start: sDate, End: eDate}, nil
}

func ParseDate(d string) (time.Time, error) {
    year, err := strconv.ParseInt(d[4:8], 10, 0)
    month, err := strconv.ParseInt(d[2:4], 10, 0)
    day, err := strconv.ParseInt(d[0:2], 10, 0)
    if err != nil {
      return GetZeroDate(), err
    }
    m, err := Month(int64(month))
    if err != nil {
      return GetZeroDate(), err
    }
    return time.Date(int(year), m, int(day), 0, 0, 0, 0, time.UTC), nil

}

func ParseJSONRequest(c *gin.Context, s interface{}) error {
  buf := new(bytes.Buffer)
  buf.ReadFrom(c.Request.Body)

  err := json.Unmarshal(buf.Bytes(), s)
  return err
}

func ReadAccFile(fName string, ad *[]AccData) (error) {
  f, err := os.Open(fName)
  if err != nil  {
    return err
  }
  reader := csv.NewReader(f)
  reader.Comma = ';'
  reader.FieldsPerRecord = 13
  // Skip Headline
  reader.Read()
  for {
    r := AccData{}
    err := UnmarshalAccData(reader, &r)
    if err == io.EOF {
      break
    } else if err != nil {
      return err
    } else {
      if (r.Belegdatum.IsZero() == true) &&
         (r.Buchungsdatum.IsZero() == true) &&
         (r.Belegnummernkreis == "") &&
         (r.Belegnummer == "") {
           continue
      }
      *ad = append(*ad, r)
    }
  }
  return nil
}

func UnmarshalAccData(reader *csv.Reader, data *AccData) error {
  s, err := reader.Read()
  if err != nil {
    return err
  }

  if (s[0] == "") && (s[1] == "") && (s[2] == "") && (s[3] == "") {
    data.Belegdatum = GetZeroDate()
    data.Buchungsdatum = GetZeroDate()
    data.Belegnummernkreis = ""
    data.Belegnummer = ""
  } else {
    /* 
    Sind die ersten vier Felder leer ist der Eintrag ein Teil einer
    Rechnung data.h die in diesem if-Block zugewiesenen Felder können nicht 
    zugewiesen werden, dies ist jedoch kein Fehler alle 
    restliche vorhanden Daten werden zugewisen was damit passiert 
    muss die aufrufende Funktion bestimmen. 
    */
    date, err := ParseGermanDate(s[0], ".")
    if err != nil {
      return errors.New("Cannot create Belegdatum")
    }
    data.Belegdatum = date

    date, err = ParseGermanDate(s[1], ".")
    if err != nil {
      return errors.New("Cannot create Buchungsdatum")
    }
    data.Buchungsdatum = date

    data.Belegnummernkreis = s[2]
    data.Belegnummer = s[3]
  }

  data.Buchungstext = s[4]
  fl, err := ParseFloatComma(s[5])
  if err != nil {
    return errors.New("Buchungstext have to be a float - "+ err.Error())
  }
  data.Buchungsbetrag = fl

  in, err := ParseAccInt(s[6])
  if err != nil {
    return errors.New("Sollkonto have to be a integer - "+ err.Error())
  }
  data.Sollkonto = in

  in, err = ParseAccInt(s[7])
  if err != nil {
    return errors.New("Habenkonto have to be a integer - "+ err.Error())
  }
  data.Habenkonto = in

  in, err = strconv.ParseInt(s[8], 10, 32)
  if err != nil {
    return errors.New("Steuerschlüssel have to be a integer - "+ err.Error())
  }
  data.Steuerschlüssel = in
  data.Kostenstelle1 = s[9]
  data.Kostenstelle2 = s[10]

  fl, err = ParseFloatComma(s[11])
  if err != nil {
    return errors.New("Buchungstext have to be a float - "+ err.Error())
  }
  data.BuchungsbetragEuro = fl
  data.Waehrung = s[12]

  return nil
}

func JoinAccFile(data []AccData, collection *mgo.Collection, validCSV bool) ([]AccFile, error) {

  fItems := []bson.M{}
  var tmp bson.M
  for i := range data {
    // Create mgo find query for each account dataset
    hKonto := strconv.FormatInt(data[i].Habenkonto, 10)
    sKonto := strconv.FormatInt(data[i].Sollkonto, 10)
    no := data[i].Belegnummernkreis + data[i].Belegnummer

    tmp = bson.M{"$or": []bson.M{

      // Find invoices
      bson.M{
        "valuetags": bson.M{
          "$elemMatch": bson.M{
            "tag": "Belegnummer",
            "value": no,
          },
        },
      },

      // Find statments
      bson.M{"$and": []bson.M{

        bson.M{
          "rangetags": bson.M{
            "$elemMatch": bson.M{
              "tag": "Belegzeitraum",
              "start": bson.M{"$lte": data[i].Belegdatum},
              "end": bson.M{"$gte": data[i].Belegdatum},
            },
          },
        },

        bson.M{
          "valuetags": bson.M{
            "$elemMatch": bson.M{
              "tag": "Kontonummer",
              "value": bson.M{"$in": []string{
                hKonto,
                sKonto,
              }},
            },
          },
        },
      }},

    }}

    fItems = append(fItems, tmp)
  }

  tmpResult := FileDocsNew([]FileDoc{})
  filter := bson.M{"$or": fItems}
  iter := collection.Find(filter).Iter()
  err := iter.All(&tmpResult.List)
  if err != nil {
    return nil, err
  }

  result := []AccFile{}
  for i, r := range data {
    q := FileDoc{
        ValueTags: []ValueTag{
            ValueTag{"Belegnummer", r.Belegnummernkreis + r.Belegnummer},
          },
        }
    docs := tmpResult.FindFile(q)

    if len(docs.List) == 0 {
      continue
    } else if len(docs.List) > 1 {
      docsJson, _ := json.Marshal(docs.List)
      errMsg := string(docsJson) +" have the same Belegnummer "+ r.Belegnummer
      return nil, errors.New(errMsg)
    } else if len(docs.List) == 1 {
      tmp := AccFile{data[i], docs.List[0]}
      result = append(result, tmp)
      data[i] = AccData{}
    }
  }
  for i, r := range data {
    docs := tmpResult.FindStat(r.Belegdatum, r.Sollkonto, r.Habenkonto)
    if len(docs.List) == 0 {
      continue
    }
    tmp := AccFile{data[i], docs.List[0]}
    result = append(result, tmp)
    data[i] = AccData{}
  }

  if validCSV == true {
    fmt.Println("Prüfe Buchhaltungsdaten")
    valid := true
    for _,r := range data {
      if r.Empty() == false {
        date := DateToString(r.Belegdatum)
        fmt.Println("\t E: ", date, r.Belegnummernkreis,
                    r.Belegnummer, r.Buchungstext, r.Sollkonto,
                    r.Habenkonto, r.Buchungsbetrag)
        valid = false
      }
    }
    if valid {
      fmt.Println("\tAlles OK!")
    }
  }

  return result, nil
}

func (ad AccData) Empty() bool {
  if (ad.Belegdatum.IsZero()) &&
    (ad.Buchungsdatum.IsZero()) &&
    (ad.Belegnummernkreis == "") &&
    (ad.Belegnummer == "") &&
    (ad.Buchungstext == "") &&
    (ad.Buchungsbetrag == 0) &&
    (ad.Sollkonto == 0) &&
    (ad.Habenkonto == 0) &&
    (ad.Steuerschlüssel == 0) &&
    (ad.Kostenstelle1 == "") &&
    (ad.Kostenstelle2 == "") &&
    (ad.BuchungsbetragEuro == 0.0) &&
    (ad.Waehrung == "") {
      return true
  } else {
    return false
  }
}

func FileDocsNew(docs []FileDoc) FileDocs {
  return FileDocs{docs}
}

func (fd FileDocs) FindStat(belegdatum time.Time, sollkonto int64, habenkonto int64) FileDocs {

  sKonto := strconv.FormatInt(sollkonto, 10)
  hKonto := strconv.FormatInt(habenkonto, 10)

  tmp := []FileDoc{}
  for i, f := range fd.List {
    findCount := 0
    for _, t := range f.RangeTags {
      if (t.Tag == "Belegzeitraum") &&
         ((t.Start.Equal(belegdatum) || t.End.Equal(belegdatum)) ||
         (t.Start.Before(belegdatum)) && (t.End.After(belegdatum))) {
           findCount += 1
           break
      }
    }
    for _, t := range f.ValueTags {
      if (t.Tag == "Kontonummer") &&
         ((t.Value == sKonto) || (t.Value == hKonto)) {
           findCount += 1
           break
      }
    }

    if findCount == 2 {
      tmp = append(tmp, fd.List[i])
    }

  }

  return FileDocsNew(tmp)

}

func (fd FileDocs) FindFile(query FileDoc) FileDocs {
  resDocs := []FileDoc{}
  for _, fileDoc := range fd.List {
    if (query.Filename != "") && (fileDoc.Filename != query.Filename) {
      continue
    }
    if len(query.SimpleTags) != 0 {
      findCount := 0
      for _, t1 := range query.SimpleTags {
        for _, t2 := range fileDoc.SimpleTags {
          if (t1.Tag == t2.Tag) {
            findCount += 1
          }
        }
      }

      if findCount != len(query.SimpleTags) {
        continue
      }

    }
    if len(query.ValueTags) != 0 {
      findCount := 0
      for _, t1 := range query.ValueTags {
        for _, t2 := range fileDoc.ValueTags {
          if (t1.Tag == t2.Tag) && (t1.Value == t2.Value) {
            findCount += 1
          }
        }
      }

      if findCount != len(query.ValueTags) {
        continue
      }
    }
    if len(query.RangeTags) != 0 {
      findCount := 0
      for _, t1 := range query.RangeTags {
        for _, t2 := range fileDoc.RangeTags {
          if (t1.Tag == t2.Tag) &&
             (t1.Start == t2.Start) &&
             (t1.End == t1.End) {
            findCount += 1
          }
        }
      }

      if findCount != len(query.ValueTags) {
        continue
      }
    }

    resDocs = append(resDocs, fileDoc)
  }

  return FileDocsNew(resDocs)
}

func ParseGermanDate(d string, sep string) (time.Time, error) {
  tmp := strings.Split(d, sep)

  dtmp, err := strconv.ParseInt(tmp[0], 10, 0)
  if err != nil {
    return GetZeroDate(), err
  }
  mtmp, err := strconv.ParseInt(tmp[1], 10, 0)
  if err != nil {
    return GetZeroDate(), err
  }
  ytmp, err := strconv.ParseInt(tmp[2], 10, 0)
  if err != nil {
    return GetZeroDate(), err
  }
  m, err := Month(mtmp)
  if err != nil {
    return GetZeroDate(), err
  }
  return time.Date(int(ytmp), m, int(dtmp), 0, 0, 0, 0, time.UTC), nil
}

func GetZeroDate() time.Time {
  return time.Date(1,time.January,1, 0,0,0,0, time.UTC)
}

func ParseFloatComma(s string) (float64, error) {
  fStr := strings.Replace(s, ".", "", -1)
  fStr = strings.Replace(fStr, ",", ".", -1)

  f, err := strconv.ParseFloat(fStr, 64)
  if err != nil {
    return 0, err
  } else {
    return f, nil
  }

}

func ParseAccInt(s string) (int64, error) {
  if s == "" {
    return -1, nil
  }

  in, err := strconv.ParseInt(s, 10, 64)
  if err != nil {
    return 0, err
  }

  return in, nil
}

func DateToString(t time.Time) (string) {
  y, m, d := t.Date()
  date := fmt.Sprintf("%02d.%02d.%d", d, int(m), y)
  return date
}

func (user *User) Load(username string, collection *mgo.Collection) error {
  u := *user
  query := collection.Find(bson.M{"username": username})

  n, err := query.Count()
  if err != nil {
    return err
  }
  if n != 1 {
    errMsg := fmt.Sprintf("Cannot find user %v", username)
    return errors.New(errMsg)
  }

  err = query.One(&u)
  if err != nil {
    return err
  }
  *user = u
  return nil

}

func (user *User) Save(col *mgo.Collection) error {
  u := *user
  u.Password = fmt.Sprintf("%x", sha1.Sum([]byte("tt")))
  err := col.Insert(u)
  if err != nil {
    return err
  }

  return nil

}

func (d *Doc) Find(db *mgo.Database) error {
  doc := *d
  docsColl := db.C(DocsColl)

  query := docsColl.Find(bson.M{"name": doc.Name})
  n, err := query.Count()
  if err != nil {
    return err
  }

  if n == 0 {
    return errors.New("Cannot find document "+ doc.Name)
  }

  if n > 1 {
    return errors.New("Found "+ strconv.Itoa(n) +" documents "+ doc.Name)
  }

  err = query.One(&doc)
  if err != nil {
    return err
  }

  *d = doc

  return nil
}

func (d *Doc) Change(changeDoc Doc, db *mgo.Database) error {
  doc := *d
  docsColl := db.C(DocsColl)

  if changeDoc.Infos.IsEmpty() == false {
      return errors.New("Not allowed to change infos!")
  }

  err := docsColl.Update(bson.M{"name": doc.Name}, changeDoc)
  if err != nil {
    return err
  }

  doc.Barcode = changeDoc.Barcode
  doc.AccountData = changeDoc.AccountData
  doc.Note = changeDoc.Note
  doc.Labels = changeDoc.Labels

  *d = doc

  return nil
}

func (d *Doc) Remove(db *mgo.Database) error {
  doc := *d
  docsColl := db.C(DocsColl)
  err := docsColl.Remove(bson.M{"name": doc.Name})
  if err != nil {
    return err
  }

  return nil
}

func (ad DocAccountData) IsEmpty() bool {
  if (ad.DocDate.IsZero()) &&
    (ad.DateOfEntry.IsZero()) &&
    (ad.DocNumberRange == "") &&
    (ad.DocNumber == "") &&
    (ad.PostingText == "") &&
    (ad.AmountPosted == 0) &&
    (ad.DebitAcc == 0) &&
    (ad.CreditAcc == 0) &&
    (ad.TaxCode == 0) &&
    (ad.CostUnit1 == "") &&
    (ad.CostUnit2 == "") &&
    (ad.AmountPostedEuro == 0.0) &&
    (ad.Currency == "") {
      return true
  } else {
    return false
  }
}

func (infos DocInfos) IsEmpty() bool {
  if (infos.DateOfScan.IsZero()) &&
     (infos.DateOfReceipt.IsZero()) {
    return true
  } else {
    return false
  }
}
