package bebber

import (
  "time"
  "gopkg.in/mgo.v2/bson"
)

const (
  VarsColl = "Vars"
  DocsColl = "Docs"
  UsersColl= "Users"
  SessionsColl = "Sessions"
  AccProcessColl = "AccProcess"
  XSRFCookieName = "XSRF-TOKEN"
  TokenHeaderField = "X-XSRF-TOKEN"
  DocNumberProposalName = "DocNumberProposal"
)

type Label string

type DocNumberProposal int
type VarDocNumberProposal struct {
  Name string
  Proposal DocNumberProposal
}

//
//  Requests and Responses
//

type DocMakeRequest Doc
type DocChangeRequest Doc

type DocRenameRequest struct {
  Name string
  NewName string
}

type DocAppendLabelsRequest struct {
  Name string
  Labels []Label
}

type DocAppendDocNumbersRequest struct {
  Name string
  DocNumbers []string
}

type AccProcessMakeRequest AccProcess


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

type MongoDBSuccessResponse struct {
  Status string
  DocID string
}

type AccProcessReadResponse struct {
  Status string
  AccProcess []AccProcessDocRef
}

type AccProcessMakeResponse struct {
  Status string
  DocID string
}

type AccProcessFindByDocNumberResponse struct {
  Status string
  AccProcessList []AccProcess
}

type AccProcessFindByAccNumberResponse struct {
  Status string
  AccProcessList []AccProcess
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
//  Accountant
//

type AccProcess struct {
  ID bson.ObjectId `bson:"_id,omitempty"`
  DocDate time.Time
  DateOfEntry time.Time
  DocNumberRange string
  DocNumber string
  PostingText string
  AmountPosted float64
  DebitAcc int
  CreditAcc int
  TaxCode int
  CostUnit1 string
  CostUnit2 string
  AmountPostedEuro float64
  Currency string
}

type AccProcessDocRef struct {
  AccProcess AccProcess
  Doc Doc
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
  DocNumbers []string
  DocPeriod DocPeriod
  AccNumber int
}

type DocPeriod struct {
  From time.Time
  To time.Time
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

