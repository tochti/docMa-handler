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
  _"regexp"
  "encoding/csv"
  "encoding/json"

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

func Month(m int64) time.Month {
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

  return months[m-1]
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

func ParseRangeTag(tag string) RangeTag {
  tmp := strings.Split(tag, TagKeyValueSep)
  key := tmp[0]
  val := strings.SplitAfter(tmp[1], RangeSep)

  var sDate time.Time
  var eDate time.Time
  if val[0] == RangeSep {
    eDate = ParseDate(val[1])
  } else if val[1] == "" {
    sDate = ParseDate(strings.Replace(val[0], RangeSep, "", 1))
  } else {
  sDate = ParseDate(strings.Replace(val[0], RangeSep, "", 1))
    eDate = ParseDate(val[1])
  }

  return RangeTag{Tag: key, Start: sDate, End: eDate}
}

func ParseDate(d string) time.Time {
    year, _ := strconv.ParseInt(d[4:8], 10, 0)
    month, _ := strconv.ParseInt(d[2:4], 10, 0)
    day, _ := strconv.ParseInt(d[0:2], 10, 0)
    return time.Date(int(year), Month(int64(month)), int(day), 0, 0, 0, 0, time.UTC)

}

func ParseJsonRequest(c *gin.Context, s interface{}) error {
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

func UnmarshalAccData(r *csv.Reader, d *AccData) error {
  s, err := r.Read()
  if err != nil {
    return err
  }

  if (s[0] == "") && (s[1] == "") && (s[2] == "") && (s[3] == "") {
    d.Belegdatum = GetZeroDate()
    d.Buchungsdatum = GetZeroDate()
    d.Belegnummernkreis = ""
    d.Belegnummer = ""
  } else {
    /* 
    Sind die ersten vier Felder leer ist der Eintrag ein Teil einer
    Rechnung d.h die in diesem if-Block zugewiesenen Felder können nicht 
    zugewiesen werden, dies ist jedoch kein Fehler alle 
    restliche vorhanden Daten werden zugewisen was damit passiert 
    muss die aufrufende Funktion bestimmen. 
    */
    date, err := ParseGermanDate(s[0], ".")
    if err != nil {
      return errors.New("Cannot create Belegdatum")
    }
    d.Belegdatum = date

    date, err = ParseGermanDate(s[1], ".")
    if err != nil {
      return errors.New("Cannot create Buchungsdatum")
    }
    d.Buchungsdatum = date

    d.Belegnummernkreis = s[2]
    d.Belegnummer = s[3]
  }

  d.Buchungstext = s[4]
  fl, err := ParseFloatComma(s[5])
  if err != nil {
    return errors.New("Buchungstext have to be a float - "+ err.Error())
  }
  d.Buchungsbetrag = fl

  in, err := ParseAccInt(s[6])
  if err != nil {
    return errors.New("Sollkonto have to be a integer - "+ err.Error())
  }
  d.Sollkonto = in

  in, err = ParseAccInt(s[7])
  if err != nil {
    return errors.New("Habenkonto have to be a integer - "+ err.Error())
  }
  d.Habenkonto = in

  in, err = strconv.ParseInt(s[8], 10, 32)
  if err != nil {
    return errors.New("Steuerschlüssel have to be a integer - "+ err.Error())
  }
  d.Steuerschlüssel = in
  d.Kostenstelle1 = s[9]
  d.Kostenstelle2 = s[10]

  fl, err = ParseFloatComma(s[11])
  if err != nil {
    return errors.New("Buchungstext have to be a float - "+ err.Error())
  }
  d.BuchungsbetragEuro = fl
  d.Waehrung = s[12]

  return nil
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

  return time.Date(int(ytmp), Month(mtmp), int(dtmp), 0, 0, 0, 0, time.UTC), nil
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
