package bebber

import (
	"os"
	"strings"
	"testing"
	"time"
)

func Test_GetSettings(t *testing.T) {
	os.Setenv("TEST_ENV", "TEST_VALUE")
	if GetSettings("TEST_ENV") != "TEST_VALUE" {
		t.Error("TEST_ENV is missing!")
	}
}

func Test_SubListOK(t *testing.T) {
	a := []string{"1", "2", "3"}
	b := []string{"2", "3"}

	diff := SubList(a, b)
	if len(diff) != 1 {
		t.Fatal("Diff should be [1] but is ", diff)
	}
	if diff[0] != "1" {
		t.Error("Diff should be [1] but is ", diff)
	}

}

func Test_SubListEmpty(t *testing.T) {
	a := []string{}
	b := []string{}

	diff := SubList(a, b)
	if len(diff) != 0 {
		t.Fatal("Diff should be [] but is ", diff)
	}
}

func Test_ParseGermanDate(t *testing.T) {
	d := time.Date(1999, time.January, 1, 0, 0, 0, 0, time.UTC)
	result, err := ParseGermanDate("01.01.1999", ".")
	if err != nil {
		t.Fatal(err)
	}
	if d != result {
		t.Error("Expect ", d, " was ", result)
	}
}

func Test_Month(t *testing.T) {
	m, _ := Month(1)
	if m != time.January {
		t.Error("Expect ", time.January, " was ", m)
	}

	m, err := Month(13)
	if err == nil {
		t.Error("Expect to throw an error due to month is out of range")
	}
}

func Test_ZeroDate(t *testing.T) {
	z := GetZeroDate()
	if z.IsZero() != true {
		t.Error("Expect zero date was ", z)
	}
}

func Test_ParseDate_OK(t *testing.T) {
	d, err := ParseDate("01012015")
	if err != nil {
		t.Fatal("Expect nil was", err)
	}
	if d.Day() != 1 {
		t.Fatal("Expect 1 was", d.Day())
	}
	if d.Month() != time.January {
		t.Fatal("Exepct", time.January, "was", d.Month())
	}
	if d.Year() != 2015 {
		t.Fatal("Expect 2015 was", d.Year())
	}
}

func Test_ParseDate_Fail(t *testing.T) {
	_, err := ParseDate("112015")
	errMsg := "Cannot parse date"
	if strings.Contains(err.Error(), errMsg) == false {
		t.Fatal("Expect", errMsg, "was", err)
	}
}
