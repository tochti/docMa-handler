package docs

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/tochti/docMa-handler/accountingData"
	"github.com/tochti/docMa-handler/labels"
	"github.com/tochti/docMa-handler/valid"
	"github.com/tochti/gin-gum/gumrest"
	"gopkg.in/gorp.v1"
)

func CreateDocHandler(ginCtx *gin.Context, db *gorp.DbMap) {
	doc := Doc{}
	if err := ginCtx.BindJSON(&doc); err != nil {
		gumrest.ErrorResponse(ginCtx, http.StatusBadRequest, err)
		return
	}

	err := valid.Struct(doc)
	if err != nil {
		gumrest.ErrorResponse(ginCtx, http.StatusBadRequest, err)
		return
	}

	err = db.Insert(&doc)
	if err != nil {
		gumrest.ErrorResponse(ginCtx, http.StatusBadRequest, err)
		return
	}

	ginCtx.JSON(http.StatusCreated, doc)
}

func ReadOneDocHandler(ginCtx *gin.Context, db *gorp.DbMap) {
	id, err := ReadDocID(ginCtx)
	if err != nil {
		return
	}

	doc := Doc{}
	err = db.SelectOne(&doc, "SELECT * FROM docs WHERE id=?", id)
	if err != nil {
		gumrest.ErrorResponse(ginCtx, http.StatusBadRequest, err)
		return
	}

	ginCtx.JSON(http.StatusOK, doc)

}

// Attention: its only possible to update the complet doc
func UpdateDocHandler(ginCtx *gin.Context, db *gorp.DbMap) {
	id, err := ReadDocID(ginCtx)
	if err != nil {
		return
	}

	doc := Doc{}
	if err := ginCtx.BindJSON(&doc); err != nil {
		gumrest.ErrorResponse(ginCtx, http.StatusBadRequest, err)
		return
	}

	doc.ID = id
	_, err = db.Update(&doc)
	if err != nil {
		gumrest.ErrorResponse(ginCtx, http.StatusBadRequest, err)
		return
	}

	ginCtx.JSON(http.StatusOK, doc)
}

func UpdateDocNameHandler(ginCtx *gin.Context, db *gorp.DbMap) {
	id, err := ReadDocID(ginCtx)
	if err != nil {
		return
	}

	doc := Doc{}
	if err := ginCtx.BindJSON(&doc); err != nil {
		gumrest.ErrorResponse(ginCtx, http.StatusBadRequest, err)
		return
	}

	if err := valid.Struct(doc); err != nil {
		gumrest.ErrorResponse(ginCtx, http.StatusBadRequest, err)
		return
	}

	q := Q("UPDATE %v SET name=? WHERE id=?", DocsTable)
	_, err = db.Exec(q, doc.Name, id)
	if err != nil {
		gumrest.ErrorResponse(ginCtx, http.StatusBadRequest, err)
		return
	}

	ginCtx.JSON(http.StatusOK, nil)
}

func CreateDocNumberHandler(ginCtx *gin.Context, db *gorp.DbMap) {
	docNumber := DocNumber{}
	if err := ginCtx.BindJSON(&docNumber); err != nil {
		gumrest.ErrorResponse(ginCtx, http.StatusBadRequest, err)
		return
	}

	if err := valid.Struct(docNumber); err != nil {
		gumrest.ErrorResponse(ginCtx, http.StatusBadRequest, err)
		return
	}

	if err := db.Insert(&docNumber); err != nil {
		gumrest.ErrorResponse(ginCtx, http.StatusBadRequest, err)
		return
	}

	ginCtx.JSON(http.StatusCreated, docNumber)
}

func ReadAllDocNumbersHandler(ginCtx *gin.Context, db *gorp.DbMap) {
	id, err := ReadDocID(ginCtx)
	if err != nil {
		return
	}

	docNumbers := []DocNumber{}
	q := Q("SELECT * FROM %v WHERE doc_id=?", DocNumbersTable)
	if _, err := db.Select(&docNumbers, q, id); err != nil {
		gumrest.ErrorResponse(ginCtx, http.StatusBadRequest, err)
		return
	}

	ginCtx.JSON(http.StatusOK, docNumbers)
}

func DeleteDocNumberHandler(ginCtx *gin.Context, db *gorp.DbMap) {
	id, err := ReadDocID(ginCtx)
	if err != nil {
		return
	}

	number, ok := ginCtx.Params.Get("docNumber")
	if !ok {
		err := errors.New("Missing number parameter")
		gumrest.ErrorResponse(ginCtx, http.StatusBadRequest, err)
		return
	}

	docNumber := DocNumber{
		DocID:  id,
		Number: number,
	}

	if _, err = db.Delete(&docNumber); err != nil {
		gumrest.ErrorResponse(ginCtx, http.StatusBadRequest, err)
		return
	}

	ginCtx.JSON(http.StatusOK, nil)

}

func CreateDocAccountDataHandler(ginCtx *gin.Context, db *gorp.DbMap) {
	docAccountData := DocAccountData{}
	if err := ginCtx.BindJSON(&docAccountData); err != nil {
		gumrest.ErrorResponse(ginCtx, http.StatusBadRequest, err)
		return
	}

	if err := valid.Struct(docAccountData); err != nil {
		gumrest.ErrorResponse(ginCtx, http.StatusBadRequest, err)
		return
	}

	if err := db.Insert(&docAccountData); err != nil {
		gumrest.ErrorResponse(ginCtx, http.StatusBadRequest, err)
		return
	}

	ginCtx.JSON(http.StatusCreated, docAccountData)
}

func ReadOneDocAccountDataHandler(ginCtx *gin.Context, db *gorp.DbMap) {
	id, err := ReadDocID(ginCtx)
	if err != nil {
		return
	}

	docAccountData := DocAccountData{}
	q := Q("SELECT * FROM %v WHERE doc_id=?", DocAccountDataTable)
	if err := db.SelectOne(&docAccountData, q, id); err != nil {
		gumrest.ErrorResponse(ginCtx, http.StatusBadRequest, err)
		return
	}

	ginCtx.JSON(http.StatusOK, docAccountData)
}

func UpdateDocAccountDataHandler(ginCtx *gin.Context, db *gorp.DbMap) {
	id, err := ReadDocID(ginCtx)
	if err != nil {
		return
	}

	docAccountData := DocAccountData{}
	if err := ginCtx.BindJSON(&docAccountData); err != nil {
		gumrest.ErrorResponse(ginCtx, http.StatusBadRequest, err)
		return
	}

	docAccountData.DocID = id
	if err := valid.Struct(docAccountData); err != nil {
		gumrest.ErrorResponse(ginCtx, http.StatusBadRequest, err)
		return
	}

	if _, err := db.Update(&docAccountData); err != nil {
		gumrest.ErrorResponse(ginCtx, http.StatusBadRequest, err)
		return
	}

	ginCtx.JSON(http.StatusOK, docAccountData)

}

func FindAllLabelsOfDocHandler(ginCtx *gin.Context, db *gorp.DbMap) {
	id, err := ReadDocID(ginCtx)
	if err != nil {
		return
	}

	labelList, err := FindLabelsOfDoc(db, id)
	if err != nil {
		gumrest.ErrorResponse(ginCtx, http.StatusBadRequest, err)
		return
	}

	ginCtx.JSON(http.StatusOK, labelList)
}

func FindAllAccountingDataOfDocHandler(ginCtx *gin.Context, db *gorp.DbMap) {
	id, err := ReadDocID(ginCtx)
	if err != nil {
		return
	}

	tmp, err := ReadDocNumbers(db, id)
	if err != nil {
		gumrest.ErrorResponse(ginCtx, http.StatusBadRequest, err)
		return
	}
	docNumbers := []string{}
	for _, x := range tmp {
		docNumbers = append(docNumbers, x.Number)
	}

	var r1 []accountingData.AccountingData
	if len(docNumbers) > 0 {
		var err error
		r1, err = accountingData.FindAccountingDataByDocNumbers(db, docNumbers)
		if err != nil {
			gumrest.ErrorResponse(ginCtx, http.StatusBadRequest, err)
			return
		}
	}

	accountData, err := ReadAccountData(db, id)
	if err != nil {
		if err != sql.ErrNoRows {
			gumrest.ErrorResponse(ginCtx, http.StatusBadRequest, err)
			return
		}
	}

	var r2 []accountingData.AccountingData
	if err != sql.ErrNoRows {
		var err error
		r2, err = accountingData.FindAccountingDataByAccountNumber(
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

func JoinLabelHandler(ginCtx *gin.Context, db *gorp.DbMap) {
	docsLabels := DocsLabels{}
	if err := ginCtx.BindJSON(&docsLabels); err != nil {
		gumrest.ErrorResponse(ginCtx, http.StatusBadRequest, err)
		return
	}

	if err := valid.Struct(docsLabels); err != nil {
		gumrest.ErrorResponse(ginCtx, http.StatusBadRequest, err)
		return
	}

	if err := db.Insert(&docsLabels); err != nil {
		gumrest.ErrorResponse(ginCtx, http.StatusBadRequest, err)
		return
	}

	ginCtx.JSON(http.StatusCreated, docsLabels)
}

func DetachLabelHandler(ginCtx *gin.Context, db *gorp.DbMap) {
	docID, err := ReadIntParam(ginCtx, "docID")
	if err != nil {
		return
	}
	labelID, err := ReadIntParam(ginCtx, "labelID")
	if err != nil {
		return
	}

	docsLabels := DocsLabels{
		DocID:   int64(docID),
		LabelID: int64(labelID),
	}

	if _, err := db.Delete(&docsLabels); err != nil {
		gumrest.ErrorResponse(ginCtx, http.StatusBadRequest, err)
		return
	}

	ginCtx.JSON(http.StatusOK, nil)
}

func FindDocsWithLabelHandler(ginCtx *gin.Context, db *gorp.DbMap) {
	labelID, err := labels.ReadLabelID(ginCtx)
	if err != nil {
		return
	}

	docs, err := FindDocsWithLabel(db, labelID)
	if err != nil {
		gumrest.ErrorResponse(ginCtx, http.StatusBadRequest, err)
		return
	}

	ginCtx.JSON(http.StatusOK, docs)
}

func mergeAccountingData(a1, a2 []accountingData.AccountingData) []accountingData.AccountingData {
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
	i, err := ReadIntParam(c, "docID")
	return int64(i), err
}

func ReadIntParam(c *gin.Context, name string) (int, error) {
	tmp := c.Params.ByName(name)
	i, err := strconv.Atoi(tmp)
	if err != nil {
		gumrest.ErrorResponse(
			c,
			http.StatusBadRequest,
			err,
		)
		return -1, err
	}

	return i, nil
}

func Q(q string, p ...interface{}) string {
	return fmt.Sprintf(q, p...)
}
