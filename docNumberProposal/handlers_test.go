package docNumberProposal

import (
	"encoding/json"
	"net/http"
	"testing"

	"gopkg.in/gorp.v1"

	"github.com/gin-gonic/gin"
	"github.com/tochti/docMa-handler/common"
	"github.com/tochti/docMa-handler/dbVars"
	"github.com/tochti/gin-gum/gumtest"
	"github.com/tochti/gin-gum/gumwrap"
)

func Test_ReadDocNumberProposalHandler(t *testing.T) {
	db := initDB(t)

	v := dbVars.DBVar{
		Name:  "docNumberProposal",
		Value: "value",
	}

	if err := db.Insert(&v); err != nil {
		t.Fatal(err)
	}

	r := gin.New()
	r.GET("/", gumwrap.Gorp(ReadDocNumberProposalHandler, db))
	resp := gumtest.NewRouter(r).ServeHTTP("GET", "/", "")

	expectResp := gumtest.JSONResponse{http.StatusOK, v}

	if err := gumtest.EqualJSONResponse(expectResp, resp); err != nil {
		t.Fatal(err)
	}

}

func Test_NextDocNumberProposalHandler(t *testing.T) {
	db := initDB(t)

	v := dbVars.DBVar{
		Name:  "docNumberProposal",
		Value: "1",
	}

	if err := db.Insert(&v); err != nil {
		t.Fatal(err)
	}

	r := gin.New()
	r.GET("/", gumwrap.Gorp(NextDocNumberProposalHandler, db))
	resp := gumtest.NewRouter(r).ServeHTTP("GET", "/", "")

	v.Value = "2"
	expectResp := gumtest.JSONResponse{http.StatusOK, v}

	if err := gumtest.EqualJSONResponse(expectResp, resp); err != nil {
		t.Fatal(err)
	}

}

func Test_UpdateDocNumberProposalHandler(t *testing.T) {
	db := initDB(t)

	v := dbVars.DBVar{
		Name:  "docNumberProposal",
		Value: "2",
	}

	if err := db.Insert(&v); err != nil {
		t.Fatal(err)
	}

	body, err := json.Marshal(v)
	if err != nil {
		t.Fatal(err)
	}

	r := gin.New()
	r.PUT("/", gumwrap.Gorp(UpdateDocNumberProposalHandler, db))
	resp := gumtest.NewRouter(r).ServeHTTP("PUT", "/", string(body))

	expectResp := gumtest.JSONResponse{http.StatusOK, v}
	if err := gumtest.EqualJSONResponse(expectResp, resp); err != nil {
		t.Fatal(err)
	}

}

func initDB(t *testing.T) *gorp.DbMap {
	return common.InitTestDB(t, dbVars.AddTables)
}
