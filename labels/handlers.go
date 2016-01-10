package labels

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/tochti/gin-gum/gumrest"
	"gopkg.in/gorp.v1"
)

var (
	ErrNameMissing = errors.New("label name is missing")
)

func CreateLabel(c *gin.Context, db *gorp.DbMap) {
	label := &Label{}

	err := c.BindJSON(label)
	if err != nil {
		gumrest.ErrorResponse(
			c,
			http.StatusBadRequest,
			err,
		)
		return
	}

	if label.Name == "" {
		gumrest.ErrorResponse(
			c,
			http.StatusBadRequest,
			ErrNameMissing,
		)
		return
	}

	err = db.Insert(label)
	if err != nil {
		gumrest.ErrorResponse(
			c,
			http.StatusBadRequest,
			err,
		)
		return
	}

	c.JSON(http.StatusCreated, label)
}

func ReadAllLabels(c *gin.Context, db *gorp.DbMap) {
	labels := []Label{}

	q := Q("SELECT id, name FROM %v", LabelsTable)

	name := c.Query("name")
	var err error
	if name != "" {
		q = fmt.Sprintf("%v WHERE name=?", q)
		_, err = db.Select(&labels, q, name)
	} else {
		_, err = db.Select(&labels, q)
	}

	if err != nil {
		gumrest.ErrorResponse(
			c,
			http.StatusNotFound,
			err,
		)
		return
	}

	c.JSON(http.StatusOK, labels)
}

func ReadOneLabel(c *gin.Context, db *gorp.DbMap) {
	labelID, err := ReadLabelID(c)
	if err != nil {
		return
	}

	label := Label{}
	q := Q("SELECT id, name FROM %v WHERE id=?", LabelsTable)
	err = db.SelectOne(&label, q, labelID)
	if err != nil {
		gumrest.ErrorResponse(
			c,
			http.StatusNotFound,
			err,
		)
		return
	}

	c.JSON(http.StatusOK, label)
}

func DeleteLabel(c *gin.Context, db *gorp.DbMap) {
	labelID, err := ReadLabelID(c)
	if err != nil {
		return
	}

	_, err = db.Delete(&Label{ID: labelID})
	if err != nil {
		gumrest.ErrorResponse(
			c,
			http.StatusBadRequest,
			err,
		)
		return
	}

	c.JSON(http.StatusOK, nil)
}

func Q(q string, p ...interface{}) string {
	return fmt.Sprintf(q, p...)
}

func ReadLabelID(c *gin.Context) (int64, error) {
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
