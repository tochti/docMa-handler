package accountingData

import (
	"errors"
	"fmt"
	"log"
	"reflect"
	"regexp"
	"time"

	"gopkg.in/gorp.v1"
)

func AddTables(db *gorp.DbMap) {
	db.AddTableWithName(AccountingData{}, AccountingDataTable).SetKeys(true, "id")
}

func FindAccountingDataByDocNumbers(db *gorp.DbMap, docNumbers []string) ([]AccountingData, error) {

	filters := `
		WHERE (
			accountingData.doc_number=?
			AND accountingData.doc_number_range=?
		)
	`
	for x := 1; x < len(docNumbers); x++ {
		q := `
			OR(
				accountingData.doc_number=?
				AND accountingData.doc_number_range=?
			)
		`
		filters = fmt.Sprintf("%v OR %v", filters, q)
	}

	q := Q(`
		SELECT *
		FROM %v as accountingData
		%v
	`, AccountingDataTable, filters)

	flatDocNumbers := []string{}
	for _, n := range docNumbers {
		rang, number, err := SplitDocNumber(n)
		if err != nil {
			return []AccountingData{}, nil
		}
		flatDocNumbers = append(flatDocNumbers, number, rang)
	}

	l := []AccountingData{}
	if _, err := db.Select(&l, q, IfaceSlice(flatDocNumbers)...); err != nil {
		return []AccountingData{}, err
	}

	return l, nil
}

func FindAccountingDataByAccountNumber(db *gorp.DbMap, accountNumber int, from time.Time, to time.Time) ([]AccountingData, error) {
	q := Q(`
		SELECT * 
		FROM %v as accountingData
		WHERE ( 
			(doc_date BETWEEN ? AND ?) 
			AND 
			(credit_account=? OR debit_account=?) 
		)
	`, AccountingDataTable)

	l := []AccountingData{}
	if _, err := db.Select(&l, q, from, to, accountNumber, accountNumber); err != nil {
		return []AccountingData{}, err
	}

	return l, nil
}

// Make []"any type" to []interface{}
func IfaceSlice(slice interface{}) []interface{} {
	s := reflect.ValueOf(slice)
	if s.Kind() != reflect.Slice {
		log.Fatal("InterfaceSlice() given a non-slice type")
	}

	ret := make([]interface{}, s.Len())

	for i := 0; i < s.Len(); i++ {
		ret[i] = s.Index(i).Interface()
	}

	return ret
}

func SplitDocNumber(docNumber string) (string, string, error) {
	reStr := "([[:alpha:]]*)(\\d+)"
	re, err := regexp.Compile(reStr)
	if err != nil {
		return "", "", err
	}
	results := re.FindStringSubmatch(docNumber)
	l := len(results)
	if (l != 2) && (l != 3) {
		err := errors.New("Invalid docnumber!")
		return "", "", err
	}

	rang := ""
	number := ""

	switch l {
	case 2:
		number = results[1]
	case 3:
		rang = results[1]
		number = results[2]
	}

	return rang, number, nil
}

func Q(q string, p ...interface{}) string {
	return fmt.Sprintf(q, p...)
}
