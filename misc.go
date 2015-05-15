package bebber

import (
  "os"
  "fmt"
  "time"
  "errors"
  "strconv"
  "strings"
  _"regexp"
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


func GetSettings(k string) string {
  v := os.Getenv(k)
  if v == "" {
    fmt.Sprintf("%v env is missing", k)
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
    year, _ := strconv.ParseInt(d[4:8], 10, 0)
    month, _ := strconv.ParseInt(d[2:4], 10, 0)
    day, _ := strconv.ParseInt(d[0:2], 10, 0)
    return time.Date(int(year), months[int(month)-1], int(day), 0, 0, 0, 0, time.UTC)

}
