package bebber

import (
  "os"
  "io"
  "errors"
  "strconv"
  "encoding/csv"
)

func ReadAccRecordsFile(fName string) ([]AccRecord, error) {
  f, err := os.Open(fName)
  if err != nil  {
    return []AccRecord{}, err
  }
  reader := csv.NewReader(f)
  reader.Comma = ';'
  reader.FieldsPerRecord = 13
  // Skip Headline
  reader.Read()
  accRecords := []AccRecord{}
  for {
    record, err := UnmarshalAccRecord(reader)
    if err == io.EOF {
      break
    } else if err != nil {
      return []AccRecord{}, err
    } else {
      // Part of a statement skip this
      if (record.DocDate.IsZero() == true) &&
         (record.DateOfEntry.IsZero() == true) &&
         (record.DocNumberRange == "") &&
         (record.DocNumber == "") {
           continue
      }
      accRecords = append(accRecords, record)
    }
  }
  return accRecords, nil
}

func UnmarshalAccRecord(reader *csv.Reader) (AccRecord, error) {
  s, err := reader.Read()
  if err != nil {
    return AccRecord{}, err
  }

  accRecord := AccRecord{}
  /* 
  Sind die ersten vier Felder leer ist der Eintrag ein Teil eines
  Kontoauszugs das heißt die in diesem if-Block zugewiesenen Felder können nicht 
  zugewiesen werden, dies ist jedoch kein Fehler alle 
  restliche vorhanden Daten werden zugewisen was damit passiert 
  muss die aufrufende Funktion bestimmen. 
  */
  if (s[0] == "") && (s[1] == "") && (s[2] == "") && (s[3] == "") {
    accRecord.DocDate = GetZeroDate()
    accRecord.DateOfEntry = GetZeroDate()
    accRecord.DocNumberRange = ""
    accRecord.DocNumber = ""
  } else {
    date, err := ParseGermanDate(s[0], ".")
    if err != nil {
      return AccRecord{},errors.New("Cannot make DocDate")
    }
    accRecord.DocDate = date

    date, err = ParseGermanDate(s[1], ".")
    if err != nil {
      return AccRecord{},errors.New("Cannot make DateOfEntry")
    }
    accRecord.DateOfEntry = date

    accRecord.DocNumberRange = s[2]
    accRecord.DocNumber = s[3]
  }

  accRecord.PostingText = s[4]

  fl, err := ParseFloatComma(s[5])
  if err != nil {
    errMsg := "Posted amount have to be a float - "+ err.Error()
    return AccRecord{}, errors.New(errMsg)
  }
  accRecord.AmountPosted = fl

  in, err := ParseAccInt(s[6])
  if err != nil {
    errMsg := "Debit account have to be a integer - "+ err.Error()
    return AccRecord{}, errors.New(errMsg)
  }
  accRecord.DebitAcc = in

  in, err = ParseAccInt(s[7])
  if err != nil {
    errMsg := "Credit account have to be a integer - "+ err.Error()
    return AccRecord{}, errors.New(errMsg)
  }
  accRecord.CreditAcc = in

  in, err = ParseAccInt(s[8])
  if err != nil {
    errMsg := "Tax code have to be a integer - "+ err.Error()
    return AccRecord{}, errors.New(errMsg)
  }
  accRecord.TaxCode = in

  accRecord.CostUnit1 = s[9]
  accRecord.CostUnit2 = s[10]

  fl, err = ParseFloatComma(s[11])
  if err != nil {
    errMsg := "Amount posted have to be a float - "+ err.Error()
    return AccRecord{}, errors.New(errMsg)
  }
  accRecord.AmountPostedEuro = fl

  accRecord.Currency = s[12]

  return accRecord, nil
}

/*
func JoinAccFile(accRecords []AccRecord, db *mgo.Database, validCSV bool) ([]AccDocRef, error) {

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
*/

func ParseAccInt(s string) (int, error) {
  if s == "" {
    return -1, nil
  }

  in, err := strconv.Atoi(s)
  if err != nil {
    return 0, err
  }

  return in, nil
}

func (ar AccRecord) IsEmpty() bool {
  if (ar.DocDate.IsZero()) &&
    (ar.DateOfEntry.IsZero()) &&
    (ar.DocNumberRange == "") &&
    (ar.DocNumber == "") &&
    (ar.PostingText == "") &&
    (ar.AmountPosted == 0) &&
    (ar.DebitAcc == 0) &&
    (ar.CreditAcc == 0) &&
    (ar.TaxCode == 0) &&
    (ar.CostUnit1 == "") &&
    (ar.CostUnit2 == "") &&
    (ar.AmountPostedEuro == 0.0) &&
    (ar.Currency == "") {
      return true
  } else {
    return false
  }
}

