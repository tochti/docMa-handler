package docs

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
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

func DeleteDocNumberHandler(ginCtx *gin.Context, db *gorp.DbMap) {
	id, err := ReadDocID(ginCtx)
	if err != nil {
		return
	}

	number, ok := ginCtx.Params.Get("number")
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

func Q(q string, p ...interface{}) string {
	return fmt.Sprintf(q, p...)
}
