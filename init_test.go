package bebber

import (
  "bytes"
  "net/http"
  "net/http/httptest"
  "path/filepath"
)

var testDir, err = filepath.Abs("./testdata")

func PerformRequest(r http.Handler, method, path string, body *bytes.Buffer) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, body)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func PerformRequestHeader(r http.Handler, method, path string, body *bytes.Buffer, header *http.Header) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, body)
  req.Header = *header
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

