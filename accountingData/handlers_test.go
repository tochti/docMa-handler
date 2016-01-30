package accountingData

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tochti/docMa-handler/common"
	"github.com/tochti/docMa-handler/docs"
	"github.com/tochti/gin-gum/gumtest"
	"github.com/tochti/gin-gum/gumwrap"
	"gopkg.in/gorp.v1"
)

func Test_CreateAccountingDataHandler(t *testing.T) {
	db := initDB(t)

	tmpD := gumtest.SimpleNow()
	accountingData := AccountingData{
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

	body, err := json.Marshal(accountingData)
	if err != nil {
		t.Fatal(err)
	}

	r := gin.New()
	r.POST("/", gumwrap.Gorp(CreateAccountingDataHandler, db))
	resp := gumtest.NewRouter(r).ServeHTTP("POST", "/", string(body))
	expectResp := gumtest.JSONResponse{http.StatusCreated, accountingData}
	if err := gumtest.EqualJSONResponse(expectResp, resp); err != nil {
		t.Fatal(err)
	}
}

func Test_MergeAccountingData(t *testing.T) {
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

	a1 := []AccountingData{accData}
	a2 := []AccountingData{accData}

	r := mergeAccountingData(a1, a2)
	if len(r) != 1 {
		t.Fatal("Expect len 1 was", len(r))
	}

	if r[0].ID != 1 {
		t.Fatal("Expect 1 was", r[0].ID)
	}
}

func Test_FindAllAccountingDataOfDocHandler(t *testing.T) {
	db := common.InitTestDB(t, AddTables, docs.AddTables)

	docNumber := docs.DocNumber{
		DocID:  1,
		Number: "DNR123",
	}

	tmpD := gumtest.SimpleNow()
	accountData := docs.DocAccountData{
		DocID:         1,
		PeriodFrom:    tmpD,
		PeriodTo:      tmpD,
		AccountNumber: 1400,
	}

	accData1 := AccountingData{
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

	accData2 := AccountingData{
		ID:               2,
		DocDate:          tmpD,
		DateOfEntry:      tmpD,
		DocNumberRange:   "X",
		DocNumber:        "1",
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

	if err := db.Insert(&accountData, &docNumber, &accData1, &accData2); err != nil {
		t.Fatal(err)
	}

	r := gin.New()
	r.GET("/:id", gumwrap.Gorp(FindAllAccountingDataOfDocHandler, db))
	resp := gumtest.NewRouter(r).ServeHTTP("GET", "/1", "")

	expectResp := gumtest.JSONResponse{
		http.StatusOK,
		[]AccountingData{
			accData1,
			accData2,
		},
	}
	if err := gumtest.EqualJSONResponse(expectResp, resp); err != nil {
		t.Fatal(err)
	}
}

func Test_FindAllAccountingDataOfDocHandler_NoDocNumbers(t *testing.T) {
	db := common.InitTestDB(t, AddTables, docs.AddTables)

	tmpD := gumtest.SimpleNow()
	accountData := docs.DocAccountData{
		DocID:         1,
		PeriodFrom:    tmpD,
		PeriodTo:      tmpD,
		AccountNumber: 1400,
	}

	accData1 := AccountingData{
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

	accData2 := AccountingData{
		ID:               2,
		DocDate:          tmpD.Add(48 * time.Hour),
		DateOfEntry:      tmpD.Add(48 * time.Hour),
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

	if err := db.Insert(&accountData, &accData1, &accData2); err != nil {
		t.Fatal(err)
	}

	r := gin.New()
	r.GET("/:id", gumwrap.Gorp(FindAllAccountingDataOfDocHandler, db))
	resp := gumtest.NewRouter(r).ServeHTTP("GET", "/1", "")

	expectResp := gumtest.JSONResponse{
		http.StatusOK,
		[]AccountingData{
			accData1,
		},
	}
	if err := gumtest.EqualJSONResponse(expectResp, resp); err != nil {
		t.Fatal(err)
	}
}

func Test_FindAllAccountingDataOfDocHandler_NoAccountData(t *testing.T) {
	db := common.InitTestDB(t, AddTables, docs.AddTables)

	docNumber := docs.DocNumber{
		DocID:  1,
		Number: "DNR123",
	}

	tmpD := gumtest.SimpleNow()
	accData1 := AccountingData{
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

	if err := db.Insert(&docNumber, &accData1); err != nil {
		t.Fatal(err)
	}

	r := gin.New()
	r.GET("/:id", gumwrap.Gorp(FindAllAccountingDataOfDocHandler, db))
	resp := gumtest.NewRouter(r).ServeHTTP("GET", "/1", "")

	expectResp := gumtest.JSONResponse{
		http.StatusOK,
		[]AccountingData{
			accData1,
		},
	}
	if err := gumtest.EqualJSONResponse(expectResp, resp); err != nil {
		t.Fatal(err)
	}
}

func initDB(t *testing.T) *gorp.DbMap {
	return common.InitTestDB(t, AddTables)
}
