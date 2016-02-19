package docs

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	"github.com/tochti/docMa-handler/labels"

	"gopkg.in/gorp.v1"
)

type (
	Interval struct {
		From time.Time
		To   time.Time
	}

	SearchForm struct {
		Labels     string
		DocNumbers string
		DateOfScan Interval
	}
)

func SearchDocs(db *gorp.DbMap, searchForm SearchForm) ([]Doc, error) {

	// If no search param set return
	if len(searchForm.Labels) == 0 &&
		len(searchForm.DocNumbers) == 0 &&
		searchForm.DateOfScan.From.IsZero() &&
		searchForm.DateOfScan.To.IsZero() {
		return []Doc{}, nil
	}

	selParam := []interface{}{}
	filters := []*bytes.Buffer{}
	froms := []string{}

	// Create labels filter
	if len(searchForm.Labels) > 0 {
		froms = append(froms, fmt.Sprintf("labels as %v", labels.LabelsTable))
		froms = append(froms, fmt.Sprintf("docs_labels as %v", DocsLabelsTable))
		l := parseQueryString(searchForm.Labels)

		filter := bytes.NewBufferString("labels.name = ?")
		selParam = append(selParam, l[0])
		for _, v := range l[1:] {
			selParam = append(selParam, v)
			filter.WriteString(" OR labels.name = ?")
		}

		filters = append(filters,
			bytes.NewBufferString(
				fmt.Sprintf(`(
					(%v)
					AND docs_labels.label_id = labels.id
					AND docs.id = docs_labels.doc_id
				)`, filter.String()),
			),
		)
	}

	// Create doc number filter
	if len(searchForm.DocNumbers) > 0 {
		froms = append(froms, fmt.Sprintf("doc_numbers as %v", DocNumbersTable))
		l := parseQueryString(searchForm.DocNumbers)

		filter := bytes.NewBufferString("doc_numbers.number = ?")
		selParam = append(selParam, l[0])
		for _, v := range l[1:] {
			selParam = append(selParam, v)
			filter.WriteString(" OR doc_numbers.number = ?")
		}

		filters = append(filters,
			bytes.NewBufferString(
				fmt.Sprintf(`(
					(%v)
					AND docs.id = doc_numbers.doc_id
				)`, filter.String()),
			),
		)

	}

	// Create date of scan filter
	if !searchForm.DateOfScan.From.IsZero() || !searchForm.DateOfScan.To.IsZero() {
		if !searchForm.DateOfScan.From.IsZero() && !searchForm.DateOfScan.To.IsZero() {
			// Filter from x to y
			filters = append(filters,
				bytes.NewBufferString("(docs.date_of_scan BETWEEN ? AND ?)"),
			)
			selParam = append(selParam,
				searchForm.DateOfScan.From,
				searchForm.DateOfScan.To,
			)
		} else if !searchForm.DateOfScan.From.IsZero() {
			// Filter from X to infinity
			filters = append(filters,
				bytes.NewBufferString("(docs.date_of_scan >= ?)"),
			)
			selParam = append(selParam, searchForm.DateOfScan.From)
		} else if !searchForm.DateOfScan.To.IsZero() {
			// Filter -infinity to X
			filters = append(filters,
				bytes.NewBufferString("(docs.date_of_scan <= ?)"),
			)
			selParam = append(selParam, searchForm.DateOfScan.To)
		}
	}

	sel := bytes.NewBufferString(`
	SELECT
		docs.name,
		docs.barcode,
		docs.date_of_scan,
		docs.date_of_receipt`)

	sel.WriteString(fmt.Sprintf(`
	FROM
		docs as %v
	`, DocsTable))
	if len(froms) == 1 {
		sel.WriteString(fmt.Sprintf(",%v\n", froms[0]))
	} else if len(froms) > 1 {
		for _, v := range froms[:len(froms)-1] {
			sel.WriteString(fmt.Sprintf(",\n%v\n", v))
		}
		sel.WriteString(fmt.Sprintf(",%v\n", froms[len(froms)-1]))
	}

	sel.WriteString("WHERE\n")
	sel.WriteString(filters[0].String())
	for _, v := range filters[1:] {
		sel.WriteString(fmt.Sprintf("\t\nAND %v", v.String()))
	}

	sel.WriteString("\nGROUP BY docs.name")

	//fmt.Println(sel)
	//fmt.Println(selParam)

	r := []Doc{}
	_, err := db.Select(&r, sel.String(), selParam...)
	if err != nil {
		return []Doc{}, err
	}

	return r, nil
}

func dateSQLFormat(t time.Time) string {
	return fmt.Sprintf("%v-%v-%v 00:00:00",
		t.Year(), t.Day(), t.Day(),
	)
}

func parseQueryString(str string) []string {
	l := strings.Split(str, ",")

	for i, v := range l {
		l[i] = strings.TrimSpace(v)
		if l[i] == "" {
			l = append(l[0:i], l[i+1:]...)
		}

	}

	return l
}
