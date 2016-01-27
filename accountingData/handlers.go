package accountingData

import (
	"net/http"

	"github.com/gin-gonic/gin"
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
