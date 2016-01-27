package docNumberProposal

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/tochti/docMa-handler/dbVars"
	"github.com/tochti/docMa-handler/valid"
	"github.com/tochti/gin-gum/gumrest"
	"gopkg.in/gorp.v1"
)

var (
	varName = "docNumberProposal"
)

func ReadDocNumberProposalHandler(ginCtx *gin.Context, db *gorp.DbMap) {
	v, err := db.Get(dbVars.DBVar{}, varName)
	if err != nil {
		gumrest.ErrorResponse(ginCtx, http.StatusInternalServerError, err)
	}

	ginCtx.JSON(http.StatusOK, v)
}

func NextDocNumberProposalHandler(ginCtx *gin.Context, db *gorp.DbMap) {
	v := dbVars.DBVar{}
	q := fmt.Sprintf("SELECT * FROM %v WHERE name=?", dbVars.DBVarsTable)
	if err := db.SelectOne(&v, q, varName); err != nil {
		gumrest.ErrorResponse(ginCtx, http.StatusInternalServerError, err)
		return
	}

	i, err := strconv.Atoi(v.Value)
	if err != nil {
		gumrest.ErrorResponse(ginCtx, http.StatusInternalServerError, err)
		return
	}

	v.Value = strconv.Itoa(i + 1)
	ginCtx.JSON(http.StatusOK, v)
}

func UpdateDocNumberProposalHandler(ginCtx *gin.Context, db *gorp.DbMap) {
	v := dbVars.DBVar{}
	if err := ginCtx.BindJSON(&v); err != nil {
		gumrest.ErrorResponse(ginCtx, http.StatusBadRequest, err)
		return
	}

	if err := valid.Struct(v); err != nil {
		gumrest.ErrorResponse(ginCtx, http.StatusBadRequest, err)
		return
	}

	if _, err := db.Update(&v); err != nil {
		gumrest.ErrorResponse(ginCtx, http.StatusBadRequest, err)
		return
	}

	ginCtx.JSON(http.StatusOK, v)
}
