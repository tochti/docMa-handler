package accountingData

import (
	"strings"
	"testing"
	"time"

	"github.com/tochti/docMa-handler/common"
	"github.com/tochti/gin-gum/gumtest"
)

func Test_FindAccountingDataByDocNumbers(t *testing.T) {
	db := common.InitTestDB(t, AddTables, AddTables)

	tmpD := gumtest.SimpleNow()
	accData := AccountingData{
		ID:               1,
		DocDate:          tmpD,
		DateOfEntry:      tmpD,
		DocNumberRange:   "DNR",
		DocNumber:        "123",
		PostingText:      "PT",
		AmountPosted:     1.1,
		DebitAccount:     1400,
		CreditAccount:    1500,
		TaxCode:          1,
		CostUnit1:        "CU1",
		CostUnit2:        "CU2",
		AmountPostedEuro: 1.2,
		Currency:         "EUR",
	}

	if err := db.Insert(&accData); err != nil {
		t.Fatal(err)
	}

	r, err := FindAccountingDataByDocNumbers(db, []string{"DNR123"})
	if err != nil {
		t.Fatal(r)
	}

	expect := []AccountingData{accData}
	if len(expect) != len(r) {
		t.Fatalf("Expect len %v was %v", len(expect), len(r))
	}

	for i, e := range expect {
		if e.ID != r[i].ID {
			t.Fatalf("Expect %v was %v", expect, r)
		}
	}

}

func Test_FindAccountingDataByDocNumbers_Multip(t *testing.T) {
	db := common.InitTestDB(t, AddTables, AddTables)

	tmpD := gumtest.SimpleNow()
	accData := AccountingData{
		ID:               1,
		DocDate:          tmpD,
		DateOfEntry:      tmpD,
		DocNumberRange:   "DNR",
		DocNumber:        "123",
		PostingText:      "PT",
		AmountPosted:     1.1,
		DebitAccount:     1400,
		CreditAccount:    1500,
		TaxCode:          1,
		CostUnit1:        "CU1",
		CostUnit2:        "CU2",
		AmountPostedEuro: 1.2,
		Currency:         "EUR",
	}

	accData1 := AccountingData{
		ID:               2,
		DocDate:          tmpD,
		DateOfEntry:      tmpD,
		DocNumberRange:   "",
		DocNumber:        "123",
		PostingText:      "PT",
		AmountPosted:     1.1,
		DebitAccount:     1400,
		CreditAccount:    1500,
		TaxCode:          1,
		CostUnit1:        "CU1",
		CostUnit2:        "CU2",
		AmountPostedEuro: 1.2,
		Currency:         "EUR",
	}

	accData2 := AccountingData{
		ID:               2,
		DocDate:          tmpD,
		DateOfEntry:      tmpD,
		DocNumberRange:   "",
		DocNumber:        "124",
		PostingText:      "PT",
		AmountPosted:     1.1,
		DebitAccount:     1400,
		CreditAccount:    1500,
		TaxCode:          1,
		CostUnit1:        "CU1",
		CostUnit2:        "CU2",
		AmountPostedEuro: 1.2,
		Currency:         "EUR",
	}

	if err := db.Insert(&accData, &accData1, &accData2); err != nil {
		t.Fatal(err)
	}

	r, err := FindAccountingDataByDocNumbers(db, []string{"123", "DNR123"})
	if err != nil {
		t.Fatal(err)
	}

	expect := []AccountingData{accData, accData1}
	if len(expect) != len(r) {
		t.Fatalf("Expect len %v was %v", len(expect), len(r))
	}

	for i, e := range expect {
		if e.ID != r[i].ID {
			t.Fatalf("Expect %v was %v", expect, r)
		}
	}

}

func Test_FindAccountingDataByAccountNumber(t *testing.T) {
	db := common.InitTestDB(t, AddTables, AddTables)

	tmpD := gumtest.SimpleNow()
	accData := AccountingData{
		ID:               1,
		DocDate:          tmpD,
		DateOfEntry:      tmpD,
		DocNumberRange:   "DNR",
		DocNumber:        "123",
		PostingText:      "PT",
		AmountPosted:     1.1,
		DebitAccount:     1400,
		CreditAccount:    1500,
		TaxCode:          1,
		CostUnit1:        "CU1",
		CostUnit2:        "CU2",
		AmountPostedEuro: 1.2,
		Currency:         "EUR",
	}

	if err := db.Insert(&accData); err != nil {
		t.Fatal(err)
	}

	from := gumtest.SimpleNow().Add(-48 * time.Hour)
	to := gumtest.SimpleNow().Add(+48 * time.Hour)
	r, err := FindAccountingDataByAccountNumber(db, 1400, from, to)
	if err != nil {
		t.Fatal(err)
	}

	expect := []AccountingData{accData}
	if len(expect) != len(r) {
		t.Fatalf("Expect len %v was %v", len(expect), len(r))
	}

	for i, e := range expect {
		if e.ID != r[i].ID {
			t.Fatalf("Expect %v was %v", expect, r)
		}
	}

}

func Test_SplitDocNumber_OK(t *testing.T) {
	rang, number, err := SplitDocNumber("BB123")
	if err != nil {
		t.Fatal("Expect nil was", err)
	}
	if rang != "BB" {
		t.Fatal("Expect BB was", rang)
	}
	if number != "123" {
		t.Fatal("Expect 123 was", number)
	}

	rang, number, err = SplitDocNumber("987")
	if err != nil {
		t.Fatal("Expect nil was", err)
	}
	if rang != "" {
		t.Fatal("Expect empty string was", rang)
	}
	if number != "987" {
		t.Fatal("Expect 987 was", number)
	}
}

func Test_SplitDocNumber_Fail(t *testing.T) {
	_, _, err := SplitDocNumber("BB")

	expectErrMsg := "Invalid docnumber!"
	if err != nil {
		if strings.Contains(err.Error(), expectErrMsg) == false {
			t.Fatal("Expect", expectErrMsg, "was nil")
		}
	} else {
		t.Fatal("Expect", expectErrMsg, "was nil")
	}
}
