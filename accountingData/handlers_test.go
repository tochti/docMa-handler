package accountingData

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/tochti/docMa-handler/common"
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

func initDB(t *testing.T) *gorp.DbMap {
	return common.InitTestDB(t, AddTables)
}
