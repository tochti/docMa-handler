package accountingData

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/tochti/docMa-handler/docs"
	"github.com/tochti/docMa-handler/valid"
	"github.com/tochti/gin-gum/gumrest"
	"gopkg.in/gorp.v1"
)

func CreateAccountingDataHandler(ginCtx *gin.Context, db *gorp.DbMap) {
	accountingData := AccountingData{}
	if err := ginCtx.BindJSON(&accountingData); err != nil {
		gumrest.ErrorResponse(ginCtx, http.StatusBadRequest, err)
		return
	}

	if err := valid.Struct(&accountingData); err != nil {
		gumrest.ErrorResponse(ginCtx, http.StatusBadRequest, err)
		return
	}

	if err := db.Insert(&accountingData); err != nil {
		gumrest.ErrorResponse(ginCtx, http.StatusBadRequest, err)
		return
	}

	ginCtx.JSON(http.StatusCreated, accountingData)

}

func FindAllAccountingDataOfDocHandler(ginCtx *gin.Context, db *gorp.DbMap) {
	id, err := ReadDocID(ginCtx)
	if err != nil {
		return
	}

	docNumbers := []docs.DocNumber{}
	_, err = db.Select(&docNumbers, Q("SELECT * FROM %v WHERE doc_id=?", docs.DocNumbersTable), id)
	if err != nil {
		gumrest.ErrorResponse(ginCtx, http.StatusBadRequest, err)
		return
	}

	var r1 []AccountingData
	if len(docNumbers) > 0 {
		var err error
		r1, err = FindAccountingDataByDocNumbers(db, docNumbers)
		if err != nil {
			gumrest.ErrorResponse(ginCtx, http.StatusBadRequest, err)
			return
		}
	}

	accountData := docs.DocAccountData{}
	err = db.SelectOne(&accountData, Q("SELECT * FROM %v WHERE doc_id=?", docs.DocAccountDataTable), id)
	if err != nil {
		if err != sql.ErrNoRows {
			gumrest.ErrorResponse(ginCtx, http.StatusBadRequest, err)
			return
		}
	}

	var r2 []AccountingData
	if err != sql.ErrNoRows {
		var err error
		r2, err = FindAccountingDataByAccountNumber(
			db,
			accountData.AccountNumber,
			accountData.PeriodFrom,
			accountData.PeriodTo,
		)
		if err != nil {
			gumrest.ErrorResponse(ginCtx, http.StatusBadRequest, err)
			return
		}
	}

	r := mergeAccountingData(r1, r2)

	ginCtx.JSON(http.StatusOK, r)
}

func mergeAccountingData(a1, a2 []AccountingData) []AccountingData {
	ids := map[int64]bool{}
	r := a1

	for _, e := range a1 {
		ids[e.ID] = true
	}

	for _, e := range a2 {
		if _, ok := ids[e.ID]; !ok {
			r = append(r, e)
		}
	}

	return r
}

func ReadDocID(c *gin.Context) (int64, error) {
	tmp := c.Params.ByName("id")
	labelID, err := strconv.ParseInt(tmp, 10, 64)
	if err != nil {
		gumrest.ErrorResponse(
			c,
			http.StatusBadRequest,
			err,
		)
		return -1, err
	}

	return labelID, nil
}
