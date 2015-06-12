package bebber

import (
  "os"
  "time"
  "bytes"
  "testing"
  "net/http"
  "net/http/httptest"
  "path/filepath"

  "gopkg.in/mgo.v2"
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

type TestRequest struct {
  Body string
  Handler http.Handler
  Header http.Header
}

func (t *TestRequest) DialToken(method, path, token string) *httptest.ResponseRecorder {
  reqData := *t
  body := bytes.NewBufferString(reqData.Body)
  reqData.Header.Add("X-XSRF-TOKEN", token)

	req, _ := http.NewRequest(method, path, body)
  req.Header = reqData.Header
	w := httptest.NewRecorder()
	reqData.Handler.ServeHTTP(w, req)
  *t = reqData
	return w
}

func CreateTestUserSession(user, token string, db *mgo.Database, t *testing.T) {
  sessionsC := db.C(SessionsCollection)
  expires := time.Now().AddDate(0,0,1)
  userSession := UserSession{Token: token, User: user, Expires: expires}
  err = sessionsC.Insert(userSession)
  if err != nil {
    t.Fatal(err.Error())
  }
}

func SetupEnvs(t *testing.T) {
  os.Setenv("BEBBER_DB_SERVER", "127.0.0.1")
  os.Setenv("BEBBER_DB_NAME", "bebber_test")
}
