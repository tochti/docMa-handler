package bebber

import (
  "fmt"
  "errors"
  "strconv"

  "gopkg.in/mgo.v2"
  "gopkg.in/mgo.v2/bson"
)

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

  docTmp := Doc{}
  err = query.One(&docTmp)
  if err != nil {
    return err
  }

  *d = docTmp

  return nil
}

func (d *Doc) Change(changeDoc Doc, db *mgo.Database) error {
  doc := *d
  docsColl := db.C(DocsColl)

  if changeDoc.Infos.IsEmpty() == false {
      return errors.New("Not allowed to change infos!")
  }

  setMap := bson.M{}
  if changeDoc.Name != "" {
    setMap["name"] = changeDoc.Name
  }
  if changeDoc.Barcode != "" {
    setMap["barcode"] = changeDoc.Barcode
  }
  if changeDoc.Note != "" {
    setMap["note"] = changeDoc.Note
  }
  if changeDoc.AccountData.IsEmpty() == false {
    if len(changeDoc.AccountData.DocNumbers) != 0 {
      setMap["accountdata.docnumbers"] = changeDoc.AccountData.DocNumbers
    }
    if changeDoc.AccountData.AccNumber != 0 {
      setMap["accountdata.accnumber"] = changeDoc.AccountData.AccNumber
    }
    if changeDoc.AccountData.DocPeriod.IsEmpty() == false {
      setMap["accountdata.docperiod"] = changeDoc.AccountData.DocPeriod
    }
  }
  if len(changeDoc.Labels) != 0 {
    setMap["labels"] = changeDoc.Labels
  }
  cD := bson.M{"$set": setMap}

  change := mgo.Change{
              Update: cD,
              ReturnNew: true,
            }
  returnDoc := Doc{}
  _, err := docsColl.Find(bson.M{"name": doc.Name}).Apply(change, &returnDoc)
  if err != nil {
    fmt.Println(err.Error())
    return err
  }

  *d = returnDoc

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

func (d *Doc) AppendLabels(labels []Label, db *mgo.Database) error {
  doc := *d
  docsColl := db.C(DocsColl)
  changeRequest := bson.M{"$addToSet": bson.M{
                             "labels": bson.M{"$each": labels},
                          }}
  change := mgo.Change{
              Update: changeRequest,
              ReturnNew: true,
            }
  returnDoc := Doc{}
  _, err := docsColl.Find(bson.M{"name": doc.Name}).Apply(change, &returnDoc)
  if err != nil {
    return err
  }
  *d = returnDoc
  return nil
}

func (d *Doc) RemoveLabels(labels []Label, db *mgo.Database) error {
  doc := *d
  docsColl := db.C(DocsColl)
  changeRequest := bson.M{"$pullAll": bson.M{"labels": labels}}
  change := mgo.Change{
              Update: changeRequest,
              ReturnNew: true,
            }
  returnDoc := Doc{}
  _, err := docsColl.Find(bson.M{"name": doc.Name}).Apply(change, &returnDoc)
  if err != nil {
    return err
  }
  *d = returnDoc
  return nil
}

func (d *Doc) AppendDocNumbers(docNumbers []string, db *mgo.Database) error {
  doc := *d
  docsColl := db.C(DocsColl)
  changeRequest := bson.M{
    "$addToSet": bson.M{
      "accountdata.docnumbers": bson.M{
        "$each": docNumbers,
      },
    },
  }
  change := mgo.Change{
              Update: changeRequest,
              ReturnNew: true,
            }
  returnDoc := Doc{}
  _, err := docsColl.Find(bson.M{"name": doc.Name}).Apply(change, &returnDoc)
  if err != nil {
    return err
  }
  *d = returnDoc
  return nil
}

func (d *Doc) RemoveDocNumbers(docNumbers []string, db *mgo.Database) error {
  doc := *d
  docsColl := db.C(DocsColl)
  changeRequest := bson.M{
    "$pullAll": bson.M{
      "accountdata.docnumbers": docNumbers,
    },
  }
  change := mgo.Change{
              Update: changeRequest,
              ReturnNew: true,
            }
  returnDoc := Doc{}
  _, err := docsColl.Find(bson.M{"name": doc.Name}).Apply(change, &returnDoc)
  if err != nil {
    return err
  }
  *d = returnDoc
  return nil
}

func (docP DocPeriod) IsEmpty() bool {
  if (docP.From.IsZero()) &&
    (docP.To.IsZero()) {
    return true
  } else {
    return false
  }

}

func (docAD DocAccountData) IsEmpty() bool {
  if (len(docAD.DocNumbers) == 0) &&
    (docAD.DocPeriod.From.IsZero()) &&
    (docAD.DocPeriod.To.IsZero()) &&
    (docAD.AccNumber == 0) {
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
