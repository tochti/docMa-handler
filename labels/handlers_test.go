package labels

import (
	"net/http"
	"testing"

	"gopkg.in/gorp.v1"

	"github.com/gin-gonic/gin"
	"github.com/tochti/docMa-handler/common"
	"github.com/tochti/gin-gum/gumtest"
	"github.com/tochti/gin-gum/gumwrap"
)

func Test_CreateLabel(t *testing.T) {
	db := initDB(t)

	r := gin.New()
	r.POST("/", gumwrap.Gorp(CreateLabel, db))

	body := `{"name": "daemon"}`

	resp := gumtest.NewRouter(r).ServeHTTP("POST", "/", body)

	label := Label{ID: 1, Name: "daemon"}
	expectResp := gumtest.JSONResponse{http.StatusCreated, label}

	if err := gumtest.EqualJSONResponse(expectResp, resp); err != nil {
		t.Fatal(err)
	}
}

func Test_CreateLabelFail(t *testing.T) {
	db := initDB(t)

	r := gin.New()
	r.POST("/", gumwrap.Gorp(CreateLabel, db))

	body := `{"ddd": "daemon"}`

	resp := gumtest.NewRouter(r).ServeHTTP("POST", "/", body)

	if resp.Code != http.StatusBadRequest {
		t.Fatalf("Expect %v was %v", http.StatusBadRequest, resp.Code)
	}
}

func Test_ReadAllLabels(t *testing.T) {
	db := initDB(t)
	labels := fillTestDB(t, db)

	r := gin.New()
	r.GET("/", gumwrap.Gorp(ReadAllLabels, db))

	resp := gumtest.NewRouter(r).ServeHTTP("GET", "/", "")

	expectResp := gumtest.JSONResponse{http.StatusOK, labels}

	if err := gumtest.EqualJSONResponse(expectResp, resp); err != nil {
		t.Fatal(err)
	}
}

func Test_ReadAllLabelsWithFilter(t *testing.T) {
	db := initDB(t)
	labels := fillTestDB(t, db)

	r := gin.New()
	r.GET("/", gumwrap.Gorp(ReadAllLabels, db))

	resp := gumtest.NewRouter(r).ServeHTTP("GET", "/?name=bad", "")

	expectResp := gumtest.JSONResponse{http.StatusOK, []*Label{labels[0]}}

	if err := gumtest.EqualJSONResponse(expectResp, resp); err != nil {
		t.Fatal(err)
	}
}

func Test_ReadOneLabel(t *testing.T) {
	db := initDB(t)
	labels := fillTestDB(t, db)

	r := gin.New()
	r.GET("/:labelID", gumwrap.Gorp(ReadOneLabel, db))

	resp := gumtest.NewRouter(r).ServeHTTP("GET", "/1", "")

	expectResp := gumtest.JSONResponse{http.StatusOK, labels[0]}

	if err := gumtest.EqualJSONResponse(expectResp, resp); err != nil {
		t.Fatal(err)
	}

}

func Test_ReadOneLabelFail(t *testing.T) {
	db := initDB(t)
	fillTestDB(t, db)

	r := gin.New()
	r.GET("/:labelID", gumwrap.Gorp(ReadOneLabel, db))

	resp := gumtest.NewRouter(r).ServeHTTP("GET", "/5", "")

	if http.StatusNotFound != resp.Code {
		t.Fatalf("Expect %v was %v", http.StatusNotFound, resp.Code)
	}
}

func Test_DeleteLabel(t *testing.T) {
	db := initDB(t)
	fillTestDB(t, db)

	r := gin.New()
	r.DELETE("/:labelID", gumwrap.Gorp(DeleteLabel, db))

	resp := gumtest.NewRouter(r).ServeHTTP("DELETE", "/1", "")

	expectResp := gumtest.JSONResponse{http.StatusOK, nil}

	if err := gumtest.EqualJSONResponse(expectResp, resp); err != nil {
		t.Fatal(err)
	}
}

func fillTestDB(t *testing.T, db *gorp.DbMap) []*Label {
	labels := []*Label{
		{1, "bad"},
		{2, "good"},
	}

	db.Insert(gumtest.IfaceSlice(labels)...)

	return labels
}

func initDB(t *testing.T) *gorp.DbMap {
	return common.InitTestDB(t, AddTables)
}
