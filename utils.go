package bebber

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func Month(m int64) (time.Month, error) {
	if (m < 1) || (m > 12) {
		return time.April, errors.New("Month out of range")
	}
	months := []time.Month{
		time.January,
		time.February,
		time.March,
		time.April,
		time.May,
		time.June,
		time.July,
		time.August,
		time.September,
		time.October,
		time.November,
		time.December,
	}

	return months[m-1], nil
}

func MakeFailResponse(c *gin.Context, msg string) {
	c.JSON(http.StatusOK, FailResponse{Status: "fail", Msg: msg})
}

func GetSettings(k string) string {
	v := os.Getenv(k)
	if v == "" {
		fmt.Println(k + " env is missing")
		os.Exit(2)
	}

	return v
}

func SubList(a, b []string) []string {
	m := make(map[string]bool)
	for _, v := range b {
		m[v] = false
	}

	var res []string
	for _, v := range a {
		_, ok := m[v]
		if ok == false {
			res = append(res, v)
		}
	}

	return res
}

// d have to be in following format ddmmyyyy
func ParseDate(d string) (time.Time, error) {
	if len(d) != 8 {
		return GetZeroDate(), errors.New("Cannot parse date " + d)
	}

	year, err := strconv.Atoi(d[4:8])
	month, err := strconv.Atoi(d[2:4])
	day, err := strconv.Atoi(d[0:2])
	if err != nil {
		return GetZeroDate(), err
	}
	m, err := Month(int64(month))
	if err != nil {
		return GetZeroDate(), err
	}

	return time.Date(year, m, day, 0, 0, 0, 0, time.UTC), nil
}

func ParseJSONRequest(c *gin.Context, s interface{}) error {
	buf := new(bytes.Buffer)
	buf.ReadFrom(c.Request.Body)

	err := json.Unmarshal(buf.Bytes(), s)
	return err
}

func ParseGermanDate(d string, sep string) (time.Time, error) {
	tmp := strings.Split(d, sep)

	dtmp, err := strconv.ParseInt(tmp[0], 10, 0)
	if err != nil {
		return GetZeroDate(), err
	}
	mtmp, err := strconv.ParseInt(tmp[1], 10, 0)
	if err != nil {
		return GetZeroDate(), err
	}
	ytmp, err := strconv.ParseInt(tmp[2], 10, 0)
	if err != nil {
		return GetZeroDate(), err
	}
	m, err := Month(mtmp)
	if err != nil {
		return GetZeroDate(), err
	}
	return time.Date(int(ytmp), m, int(dtmp), 0, 0, 0, 0, time.UTC), nil
}

func GetZeroDate() time.Time {
	return time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC)
}

func ParseFloatComma(s string) (float64, error) {
	fStr := strings.Replace(s, ".", "", -1)
	fStr = strings.Replace(fStr, ",", ".", -1)

	f, err := strconv.ParseFloat(fStr, 64)
	if err != nil {
		return 0, err
	} else {
		return f, nil
	}

}

func DateToString(t time.Time) string {
	y, m, d := t.Date()
	date := fmt.Sprintf("%02d.%02d.%d", d, int(m), y)
	return date
}
