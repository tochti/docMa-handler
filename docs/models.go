package docs

import "time"

var (
	DocsTable           = "docs"
	DocNumbersTable     = "doc_numbers"
	DocAccountDataTable = "account_data"
)

type Doc struct {
	ID            int64     `db:"id" json:"id"`
	Name          string    `db:"name" json:"name" valid:"required"`
	Barcode       string    `db:"barcode" json:"barcode"`
	DateOfScan    time.Time `db:"date_of_scan" json:"date_of_scan"`
	DateOfReceipt time.Time `db:"date_of_receipt" json:"date_of_receipt"`
	Note          string    `db:"note" json:"note"`
}

type DocAccountData struct {
	DocID         int64     `db:"doc_id" json:"doc_id" valid:"required,gt=0"`
	PeriodFrom    time.Time `db:"period_from" json:"period_from"`
	PeriodTo      time.Time `db:"period_to" json:"period_to"`
	AccountNumber int       `db:"account_number" json:"account_number"`
}

type DocNumber struct {
	DocID  int64  `db:"doc_id" json:"doc_id" valid:"required,gt=0"`
	Number string `db:"number" json:"number" valid:"required"`
}
