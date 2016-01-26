package docs

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"gopkg.in/gorp.v1"

	"github.com/gin-gonic/gin"
	"github.com/tochti/docMa-handler/common"
	"github.com/tochti/gin-gum/gumrest"
	"github.com/tochti/gin-gum/gumtest"
	"github.com/tochti/gin-gum/gumwrap"
)

func initDB(t *testing.T) *gorp.DbMap {
	db := common.InitTestDB(t, AddTables)

	return db
}

func Test_CreateDoc(t *testing.T) {
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
	r.GET("/:id", gumwrap.Gorp(ReadOneDocHandler, db))
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
	r.PUT("/:id", gumwrap.Gorp(UpdateDocHandler, db))
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
	r.PATCH("/:id/name", gumwrap.Gorp(UpdateDocNameHandler, db))
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

func Test_CreateDocNumber(t *testing.T) {
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

func Test_CreateDocNumber_MissingDocID(t *testing.T) {
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

func Test_CreateDocNumber_MissingNumber(t *testing.T) {
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

func Test_DeleteDocNumber(t *testing.T) {
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
	r.DELETE("/:id/:number", gumwrap.Gorp(DeleteDocNumberHandler, db))
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
	r.GET("/:id", gumwrap.Gorp(ReadOneDocAccountDataHandler, db))
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
	r.PUT("/:id", gumwrap.Gorp(UpdateDocAccountDataHandler, db))
	resp := gumtest.NewRouter(r).ServeHTTP("PUT", "/2", body)
	expectResp := gumtest.JSONResponse{http.StatusOK, accountData}
	if err := gumtest.EqualJSONResponse(expectResp, resp); err != nil {
		t.Fatal(err)
	}
}
