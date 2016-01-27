package accountingData

import "time"

var (
	AccountingDataTable = "accounting_data"
)

type AccountingData struct {
	ID               int64     `db:"id" json:"id"`
	DocDate          time.Time `db:"doc_date" json:"doc_date"`
	DateOfEntry      time.Time `db:"date_of_entry" json:"date_of_entry"`
	DocNumberRange   string    `db:"doc_number_range" json:"doc_number_range"`
	DocNumber        string    `db:"doc_number" json:"doc_number"`
	PostingText      string    `db:"posting_text" json:"posting_text"`
	AmountPosted     float64   `db:"amount_posted" json:"amount_posted"`
	DebitAccount     int       `db:"debit_account" json:"debit_account"`
	CreditAccount    int       `db:"credit_account" json:"credit_account"`
	TaxCode          int       `db:"tax_code" json:"tax_code"`
	CostUnit1        string    `db:"cost_unit1" json:"cost_unit1"`
	CostUnit2        string    `db:"cost_unit2" json:"cost_unit2"`
	AmountPostedEuro float64   `db:"amount_posted_euro", json:"amount_posted_euro"`
	Currency         string    `db:"currency", json:"currency"`
}
