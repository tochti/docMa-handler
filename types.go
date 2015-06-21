package bebber

import (
  "time"
  "gopkg.in/mgo.v2/bson"
)

const (
  DocsColl = "Docs"
  UsersCollection = "Users"
  XSRFCookieName = "XSRF-TOKEN"
)

type Label string

//
//  Requests and Responses
//

type DocMakeRequest Doc
type DocChangeRequest Doc

type MongoDBSuccessResponse struct {
  Status string
  DocID string
}

type DocAppendLabelsRequest struct {
  Name string
  Labels []Label
}

type DocRemoveLabelsRequest struct {
  Name string
  Labels []Label
}

type SuccessResponse struct {
  Status string
}

type FailResponse struct {
  Status string
  Msg string
}

type DocMakeResponse struct {
  Status string
  Id bson.ObjectId
}

type DocReadResponse struct {
  Status string
  Doc Doc
}

//
// Authentication
//

type User struct {
  Username string
  Password string
}

type LoginData struct {
  Username string
  Password string
}

type UserSession struct {
  Token string
  User string
  Expires time.Time
}

//
//  Doc
//

type DocInfos struct {
  DateOfScan time.Time
  DateOfReceipt time.Time
}

type DocNote string

type DocAccountData struct {
  DocDate time.Time
  DateOfEntry time.Time
  DocNumberRange string
  DocNumber string
  PostingText string
  AmountPosted float64
  DebitAcc int64
  CreditAcc int64
  TaxCode int64
  CostUnit1 string
  CostUnit2 string
  AmountPostedEuro float64
  Currency string
}

type Doc struct {
  ID bson.ObjectId `bson:"_id,omitempty"`
  Name string
  Barcode string
  Infos DocInfos
  Note DocNote
  AccountData DocAccountData
  Labels []Label
}

