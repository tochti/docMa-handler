package docs

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"gopkg.in/gorp.v1"

	"github.com/gin-gonic/gin"
	"github.com/tochti/docMa-handler/accountingData"
	"github.com/tochti/docMa-handler/common"
	"github.com/tochti/docMa-handler/labels"
	"github.com/tochti/gin-gum/gumrest"
	"github.com/tochti/gin-gum/gumtest"
	"github.com/tochti/gin-gum/gumwrap"
)

func Test_CreateDocHandler(t *testing.T) {
	db := initDB(t)

	doc := Doc{
		ID:            1,
		Name:          "darkmoon.txt",
		Barcode:       "darkmoon",
		DateOfScan:    gumtest.SimpleNow(),
		DateOfReceipt: gumtest.SimpleNow(),
		Note:          "There was a man...",
	}

	body, err := json.Marshal(doc)
	if err != nil {
		t.Fatal(err)
	}

	r := gin.New()
	r.POST("/", gumwrap.Gorp(CreateDocHandler, db))
	resp := gumtest.NewRouter(r).ServeHTTP("POST", "/", string(body))
	expectResp := gumtest.JSONResponse{http.StatusCreated, doc}

	if err := gumtest.EqualJSONResponse(expectResp, resp); err != nil {
		t.Fatal(err)
	}

}

func Test_CreateDocHandler_MissingName(t *testing.T) {
	db := initDB(t)

	doc := Doc{
		ID:            1,
		Barcode:       "darkmoon",
		DateOfScan:    gumtest.SimpleNow(),
		DateOfReceipt: gumtest.SimpleNow(),
		Note:          "There was a man...",
	}

	body, err := json.Marshal(doc)
	if err != nil {
		t.Fatal(err)
	}

	r := gin.New()
	r.POST("/", gumwrap.Gorp(CreateDocHandler, db))
	resp := gumtest.NewRouter(r).ServeHTTP("POST", "/", string(body))

	expectResp := gumtest.JSONResponse{
		http.StatusBadRequest,
		gumrest.ErrorMessage{
			Message: "Key: 'Doc.Name' Error:Field validation for 'Name' failed on the 'required' tag",
		},
	}
	if err := gumtest.EqualJSONResponse(expectResp, resp); err != nil {
		t.Fatal(err)
	}
}

func Test_CreateDocHandler_AlreadyExists(t *testing.T) {
	db := initDB(t)

	doc := Doc{
		ID:            1,
		Name:          "darkmoon.txt",
		Barcode:       "darkmoon",
		DateOfScan:    gumtest.SimpleNow(),
		DateOfReceipt: gumtest.SimpleNow(),
		Note:          "There was a man...",
	}

	err := db.Insert(&doc)
	if err != nil {
		t.Fatal(err)
	}

	body, err := json.Marshal(doc)
	if err != nil {
		t.Fatal(err)
	}

	r := gin.New()
	r.POST("/", gumwrap.Gorp(CreateDocHandler, db))
	resp := gumtest.NewRouter(r).ServeHTTP("POST", "/", string(body))

	expectResp := gumtest.JSONResponse{
		http.StatusBadRequest,
		gumrest.ErrorMessage{
			Message: "Error 1062: Duplicate entry 'darkmoon.txt' for key 'name'",
		},
	}
	if err := gumtest.EqualJSONResponse(expectResp, resp); err != nil {
		t.Fatal(err)
	}
}

func Test_ReadOneDocHandler(t *testing.T) {
	db := initDB(t)

	doc := Doc{
		ID:            1,
		Name:          "darkmoon.txt",
		Barcode:       "darkmoon",
		DateOfScan:    gumtest.SimpleNow(),
		DateOfReceipt: gumtest.SimpleNow(),
		Note:          "There was a man...",
	}

	err := db.Insert(&doc)
	if err != nil {
		t.Fatal(err)
	}

	r := gin.New()
	r.GET("/:docID", gumwrap.Gorp(ReadOneDocHandler, db))
	resp := gumtest.NewRouter(r).ServeHTTP("GET", "/1", "")

	expectResp := gumtest.JSONResponse{
		http.StatusOK,
		doc,
	}
	if err := gumtest.EqualJSONResponse(expectResp, resp); err != nil {
		t.Fatal(err)
	}
}

func Test_UpdateDocHandler(t *testing.T) {
	db := initDB(t)

	doc := Doc{
		ID:            1,
		Name:          "darkmoon.txt",
		Barcode:       "darkmoon",
		DateOfScan:    gumtest.SimpleNow(),
		DateOfReceipt: gumtest.SimpleNow(),
		Note:          "There was a man...",
	}

	err := db.Insert(&doc)
	if err != nil {
		t.Fatal(err)
	}

	body := `{
		"name": "fungi.txt", 
		"barcode": "fungi",
		"date_of_scan": "2012-04-23T18:00:00Z",
		"date_of_receipt": "2012-04-23T18:01:00Z"
	}`

	r := gin.New()
	r.PUT("/:docID", gumwrap.Gorp(UpdateDocHandler, db))
	resp := gumtest.NewRouter(r).ServeHTTP("PUT", "/1", body)

	doc = Doc{
		ID:            1,
		Name:          "fungi.txt",
		Barcode:       "fungi",
		DateOfScan:    time.Date(2012, time.April, 23, 18, 0, 0, 0, time.UTC),
		DateOfReceipt: time.Date(2012, time.April, 23, 18, 1, 0, 0, time.UTC),
		Note:          "",
	}
	expectResp := gumtest.JSONResponse{
		http.StatusOK,
		doc,
	}
	if err := gumtest.EqualJSONResponse(expectResp, resp); err != nil {
		t.Fatal(err)
	}
}

func Test_UpdateDocNameHandler(t *testing.T) {
	db := initDB(t)

	doc := Doc{
		ID:            1,
		Name:          "darkmoon.txt",
		Barcode:       "darkmoon",
		DateOfScan:    gumtest.SimpleNow(),
		DateOfReceipt: gumtest.SimpleNow(),
		Note:          "There was a man...",
	}

	err := db.Insert(&doc)
	if err != nil {
		t.Fatal(err)
	}

	body := `{"name": "fungi.txt"}`

	r := gin.New()
	r.PATCH("/:docID/name", gumwrap.Gorp(UpdateDocNameHandler, db))
	resp := gumtest.NewRouter(r).ServeHTTP("PATCH", "/1/name", body)

	doc.Name = "fungi.txt"
	expectResp := gumtest.JSONResponse{
		http.StatusOK,
		nil,
	}
	if err := gumtest.EqualJSONResponse(expectResp, resp); err != nil {
		t.Fatal(err)
	}
}

func Test_CreateDocNumberHandler(t *testing.T) {
	db := initDB(t)

	doc := DocNumber{
		DocID:  1,
		Number: "1",
	}

	body, err := json.Marshal(doc)
	if err != nil {
		t.Fatal(err)
	}

	r := gin.New()
	r.POST("/", gumwrap.Gorp(CreateDocNumberHandler, db))
	resp := gumtest.NewRouter(r).ServeHTTP("POST", "/", string(body))
	expectResp := gumtest.JSONResponse{http.StatusCreated, doc}
	if err := gumtest.EqualJSONResponse(expectResp, resp); err != nil {
		t.Fatal(err)
	}

}

func Test_CreateDocNumberHandler_MissingDocID(t *testing.T) {
	db := initDB(t)

	body := `{"number":"1"}`

	r := gin.New()
	r.POST("/", gumwrap.Gorp(CreateDocNumberHandler, db))
	resp := gumtest.NewRouter(r).ServeHTTP("POST", "/", body)
	expectResp := gumtest.JSONResponse{
		http.StatusBadRequest,
		gumrest.ErrorMessage{
			Message: "Key: 'DocNumber.DocID' Error:Field validation for 'DocID' failed on the 'required' tag",
		},
	}
	if err := gumtest.EqualJSONResponse(expectResp, resp); err != nil {
		t.Fatal(err)
	}

}

func Test_CreateDocNumberHandler_MissingNumber(t *testing.T) {
	db := initDB(t)

	body := `{"doc_id": 1}`

	r := gin.New()
	r.POST("/", gumwrap.Gorp(CreateDocNumberHandler, db))
	resp := gumtest.NewRouter(r).ServeHTTP("POST", "/", body)
	expectResp := gumtest.JSONResponse{
		http.StatusBadRequest,
		gumrest.ErrorMessage{
			Message: "Key: 'DocNumber.Number' Error:Field validation for 'Number' failed on the 'required' tag",
		},
	}
	if err := gumtest.EqualJSONResponse(expectResp, resp); err != nil {
		t.Fatal(err)
	}

}

func Test_ReadAllDocNumbers(t *testing.T) {
	db := initDB(t)

	docNumber := DocNumber{
		DocID:  1,
		Number: "v",
	}
	err := db.Insert(&docNumber)
	if err != nil {
		t.Fatal(err)
	}

	r := gin.New()
	r.GET("/:docID", gumwrap.Gorp(ReadAllDocNumbersHandler, db))
	resp := gumtest.NewRouter(r).ServeHTTP("GET", "/1", "")
	expectResp := gumtest.JSONResponse{
		http.StatusOK,
		[]DocNumber{docNumber},
	}
	if err := gumtest.EqualJSONResponse(expectResp, resp); err != nil {
		t.Fatal(err)
	}
}

func Test_DeleteDocNumberHandler(t *testing.T) {
	db := initDB(t)

	docNumber := DocNumber{
		DocID:  2,
		Number: "1",
	}

	err := db.Insert(&docNumber)
	if err != nil {
		t.Fatal(err)
	}

	r := gin.New()
	r.DELETE("/:docID/:docNumber", gumwrap.Gorp(DeleteDocNumberHandler, db))
	resp := gumtest.NewRouter(r).ServeHTTP("DELETE", "/2/1", "")
	expectResp := gumtest.JSONResponse{http.StatusOK, nil}

	if err := gumtest.EqualJSONResponse(expectResp, resp); err != nil {
		t.Fatal(err)
	}

}

func Test_CreateDocAccountDataHandler(t *testing.T) {
	db := initDB(t)

	accountData := DocAccountData{
		DocID:         2,
		PeriodFrom:    gumtest.SimpleNow(),
		PeriodTo:      gumtest.SimpleNow(),
		AccountNumber: 2,
	}

	body, err := json.Marshal(accountData)
	if err != nil {
		t.Fatal(err)
	}

	r := gin.New()
	r.POST("/", gumwrap.Gorp(CreateDocAccountDataHandler, db))
	resp := gumtest.NewRouter(r).ServeHTTP("POST", "/", string(body))
	expectResp := gumtest.JSONResponse{http.StatusCreated, accountData}
	if err := gumtest.EqualJSONResponse(expectResp, resp); err != nil {
		t.Fatal(err)
	}

}

func Test_CreateDocAccountDataHandler_MissingDocID(t *testing.T) {
	db := initDB(t)

	accountData := DocAccountData{
		PeriodFrom:    gumtest.SimpleNow(),
		PeriodTo:      gumtest.SimpleNow(),
		AccountNumber: 2,
	}

	body, err := json.Marshal(accountData)
	if err != nil {
		t.Fatal(err)
	}

	r := gin.New()
	r.POST("/", gumwrap.Gorp(CreateDocAccountDataHandler, db))
	resp := gumtest.NewRouter(r).ServeHTTP("POST", "/", string(body))
	expectResp := gumtest.JSONResponse{
		http.StatusBadRequest,
		gumrest.ErrorMessage{
			Message: "Key: 'DocAccountData.DocID' Error:Field validation for 'DocID' failed on the 'required' tag",
		},
	}
	if err := gumtest.EqualJSONResponse(expectResp, resp); err != nil {
		t.Fatal(err)
	}
}

func Test_ReadOneDocAccountDataHandler(t *testing.T) {
	db := initDB(t)

	accountData := DocAccountData{
		DocID:         2,
		PeriodFrom:    gumtest.SimpleNow(),
		PeriodTo:      gumtest.SimpleNow(),
		AccountNumber: 2,
	}

	err := db.Insert(&accountData)
	if err != nil {
		t.Fatal(err)
	}

	r := gin.New()
	r.GET("/:docID", gumwrap.Gorp(ReadOneDocAccountDataHandler, db))
	resp := gumtest.NewRouter(r).ServeHTTP("GET", "/2", "")
	expectResp := gumtest.JSONResponse{http.StatusOK, accountData}
	if err := gumtest.EqualJSONResponse(expectResp, resp); err != nil {
		t.Fatal(err)
	}

}

func Test_UpdateDocAccountDataHandler(t *testing.T) {
	db := initDB(t)

	accountData := DocAccountData{
		DocID:         2,
		PeriodFrom:    time.Date(2012, time.April, 23, 18, 0, 0, 0, time.UTC),
		PeriodTo:      time.Date(2012, time.April, 23, 18, 1, 0, 0, time.UTC),
		AccountNumber: 2,
	}

	err := db.Insert(&accountData)
	if err != nil {
		t.Fatal(err)
	}

	body := `{
		"period_from": "2012-04-23T18:00:00Z",
		"period_to": "2012-04-23T18:01:00Z",
		"account_number": 2
	}`

	r := gin.New()
	r.PUT("/:docID", gumwrap.Gorp(UpdateDocAccountDataHandler, db))
	resp := gumtest.NewRouter(r).ServeHTTP("PUT", "/2", body)
	expectResp := gumtest.JSONResponse{http.StatusOK, accountData}
	if err := gumtest.EqualJSONResponse(expectResp, resp); err != nil {
		t.Fatal(err)
	}
}

func Test_FindAllLabelsOfDocHandler(t *testing.T) {
	db := common.InitTestDB(t, AddTables, labels.AddTables)

	label := labels.Label{
		ID:   1,
		Name: "label",
	}

	doc := Doc{
		ID:   1,
		Name: "karl.pdf",
	}

	docsLabels := DocsLabels{
		DocID:   doc.ID,
		LabelID: label.ID,
	}

	if err := db.Insert(&label); err != nil {
		t.Fatal(err)
	}

	if err := db.Insert(&doc); err != nil {
		t.Fatal(err)
	}

	if err := db.Insert(&docsLabels); err != nil {
		t.Fatal(err)
	}

	r := gin.New()
	r.GET("/docs/:docID/labels", gumwrap.Gorp(FindAllLabelsOfDocHandler, db))
	resp := gumtest.NewRouter(r).ServeHTTP("GET", "/docs/1/labels", "")
	expectResp := gumtest.JSONResponse{http.StatusOK, []labels.Label{label}}
	if err := gumtest.EqualJSONResponse(expectResp, resp); err != nil {
		t.Fatal(err)
	}
}

func Test_MergeAccountingData(t *testing.T) {
	tmpD := gumtest.SimpleNow()
	accData := accountingData.AccountingData{
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

	a1 := []accountingData.AccountingData{accData}
	a2 := []accountingData.AccountingData{accData}

	r := mergeAccountingData(a1, a2)
	if len(r) != 1 {
		t.Fatal("Expect len 1 was", len(r))
	}

	if r[0].ID != 1 {
		t.Fatal("Expect 1 was", r[0].ID)
	}
}

func Test_FindAllAccountingDataOfDocHandler(t *testing.T) {
	db := common.InitTestDB(t, AddTables, accountingData.AddTables)

	docNumber := DocNumber{
		DocID:  1,
		Number: "DNR123",
	}

	tmpD := gumtest.SimpleNow()
	accountData := DocAccountData{
		DocID:         1,
		PeriodFrom:    tmpD,
		PeriodTo:      tmpD,
		AccountNumber: 1400,
	}

	accData1 := accountingData.AccountingData{
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

	accData2 := accountingData.AccountingData{
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
	r.GET("/:docID", gumwrap.Gorp(FindAllAccountingDataOfDocHandler, db))
	resp := gumtest.NewRouter(r).ServeHTTP("GET", "/1", "")

	expectResp := gumtest.JSONResponse{
		http.StatusOK,
		[]accountingData.AccountingData{
			accData1,
			accData2,
		},
	}
	if err := gumtest.EqualJSONResponse(expectResp, resp); err != nil {
		t.Fatal(err)
	}
}

func Test_FindAllAccountingDataOfDocHandler_NoDocNumbers(t *testing.T) {
	db := common.InitTestDB(t, AddTables, accountingData.AddTables)

	tmpD := gumtest.SimpleNow()
	accountData := DocAccountData{
		DocID:         1,
		PeriodFrom:    tmpD,
		PeriodTo:      tmpD,
		AccountNumber: 1400,
	}

	accData1 := accountingData.AccountingData{
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

	accData2 := accountingData.AccountingData{
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
	r.GET("/:docID", gumwrap.Gorp(FindAllAccountingDataOfDocHandler, db))
	resp := gumtest.NewRouter(r).ServeHTTP("GET", "/1", "")

	expectResp := gumtest.JSONResponse{
		http.StatusOK,
		[]accountingData.AccountingData{
			accData1,
		},
	}
	if err := gumtest.EqualJSONResponse(expectResp, resp); err != nil {
		t.Fatal(err)
	}
}

func Test_FindAllAccountingDataOfDocHandler_NoAccountData(t *testing.T) {
	db := common.InitTestDB(t, AddTables, accountingData.AddTables)

	docNumber := DocNumber{
		DocID:  1,
		Number: "DNR123",
	}

	tmpD := gumtest.SimpleNow()
	accData1 := accountingData.AccountingData{
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
	r.GET("/:docID", gumwrap.Gorp(FindAllAccountingDataOfDocHandler, db))
	resp := gumtest.NewRouter(r).ServeHTTP("GET", "/1", "")

	expectResp := gumtest.JSONResponse{
		http.StatusOK,
		[]accountingData.AccountingData{
			accData1,
		},
	}
	if err := gumtest.EqualJSONResponse(expectResp, resp); err != nil {
		t.Fatal(err)
	}
}

func Test_JoinLabelHandler(t *testing.T) {
	db := common.InitTestDB(t, AddTables, labels.AddTables)

	docsLabels := DocsLabels{
		DocID:   1,
		LabelID: 1,
	}

	body, err := json.Marshal(docsLabels)
	if err != nil {
		t.Fatal(err)
	}

	r := gin.New()
	r.POST("/", gumwrap.Gorp(JoinLabelHandler, db))
	resp := gumtest.NewRouter(r).ServeHTTP("POST", "/", string(body))

	expectResp := gumtest.JSONResponse{
		http.StatusCreated,
		docsLabels,
	}
	if err := gumtest.EqualJSONResponse(expectResp, resp); err != nil {
		t.Fatal(err)
	}
}

func Test_DetachLabelHandler(t *testing.T) {
	db := common.InitTestDB(t, AddTables, labels.AddTables)

	a := DocsLabels{
		DocID:   1,
		LabelID: 1,
	}
	if err := db.Insert(&a); err != nil {
		t.Fatal(err)
	}

	r := gin.New()
	r.DELETE("/:docID/:labelID", gumwrap.Gorp(DetachLabelHandler, db))
	resp := gumtest.NewRouter(r).ServeHTTP("DELETE", "/1/1", "")

	expectResp := gumtest.JSONResponse{
		http.StatusOK,
		nil,
	}
	if err := gumtest.EqualJSONResponse(expectResp, resp); err != nil {
		t.Fatal(err)
	}
}

func Test_FindDocsWithLabelHandler(t *testing.T) {
	db := common.InitTestDB(t, AddTables, labels.AddTables)

	doc := Doc{
		ID:   1,
		Name: "test1.pdf",
	}

	label := labels.Label{
		ID:   1,
		Name: "label",
	}

	docsLabels := DocsLabels{
		DocID:   1,
		LabelID: 1,
	}
	if err := db.Insert(&doc, &label, &docsLabels); err != nil {
		t.Fatal(err)
	}

	r := gin.New()
	r.GET("/labels/:labelID/docs", gumwrap.Gorp(FindDocsWithLabelHandler, db))
	resp := gumtest.NewRouter(r).ServeHTTP("GET", "/labels/1/docs", "")

	expectResp := gumtest.JSONResponse{
		http.StatusOK,
		[]Doc{doc},
	}
	if err := gumtest.EqualJSONResponse(expectResp, resp); err != nil {
		t.Fatal(err)
	}
}

func Test_SearchDocsHandler(t *testing.T) {
	db := common.InitTestDB(t, AddTables, labels.AddTables)

	doc := Doc{
		ID:   1,
		Name: "test1.pdf",
	}

	label := labels.Label{
		ID:   1,
		Name: "label",
	}

	docsLabels := DocsLabels{
		DocID:   1,
		LabelID: 1,
	}
	if err := db.Insert(&doc, &label, &docsLabels); err != nil {
		t.Fatal(err)
	}

	body := `
	{
		"labels": "label"
	}
	`

	r := gin.New()
	r.POST("/", gumwrap.Gorp(SearchDocsHandler, db))
	resp := gumtest.NewRouter(r).ServeHTTP("POST", "/", body)

	expectResp := gumtest.JSONResponse{
		http.StatusOK,
		[]Doc{doc},
	}
	if err := gumtest.EqualJSONResponse(expectResp, resp); err != nil {
		t.Fatal(err)
	}
}

func Test_ReadIntParam(t *testing.T) {
	passed := false
	h := func(c *gin.Context) {
		i, err := ReadIntParam(c, "docID")
		if err != nil {
			t.Fatal(err)
		}

		if i != int(1) {
			t.Fatal("Expect 1 was", 1)
		}

		passed = true
	}

	r := gin.New()
	r.GET("/:docID", h)
	gumtest.NewRouter(r).ServeHTTP("GET", "/1", "")

	if !passed {
		t.Fatal("Didn't passed all tests")
	}
}

func Test_ReadIntParam_CannotFoundParam(t *testing.T) {
	passed := false
	h := func(c *gin.Context) {
		_, err := ReadIntParam(c, "id")
		if err == nil {
			t.Fatal("Expect err not to be nil")
		}

		passed = true
	}

	r := gin.New()
	r.GET("/:none", h)
	resp := gumtest.NewRouter(r).ServeHTTP("GET", "/1", "")
	expectResp := gumtest.JSONResponse{
		http.StatusBadRequest,
		gumrest.ErrorMessage{
			Message: `strconv.ParseInt: parsing "": invalid syntax`,
		},
	}
	if err := gumtest.EqualJSONResponse(expectResp, resp); err != nil {
		t.Fatal(err)
	}

	if !passed {
		t.Fatal("Didn't passed all tests")
	}
}

func initDB(t *testing.T) *gorp.DbMap {
	db := common.InitTestDB(t, AddTables)

	return db
}
