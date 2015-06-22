package bebber

import (
  "time"
  "testing"
  "net/http"
  "encoding/json"
  "github.com/gin-gonic/gin"
)

func Test_VerifyAuth_OK(t *testing.T) {
  globals := MakeTestGlobals(t)
  db := globals.MongoDB.Session.DB(TestDBName)
  defer globals.MongoDB.Session.Close()
  defer db.DropDatabase()

  sessionsColl := db.C(SessionsColl)
  expires := time.Now().AddDate(0,0,1)
  expectSession := UserSession{Token: "123", User: "loveMaster_999", Expires: expires}
  sessionsColl.Insert(expectSession)

  h := gin.New()
  afterAuth := func (c *gin.Context) {
    se, _ := c.Get("session")
    c.JSON(http.StatusOK, se)
  }
  h.GET("/", MakeGlobalsHandler(Auth, globals), afterAuth)

  request := TestRequest{
                Body: "",
                Header: http.Header{},
                Handler: h,
             }
  response := request.SendWithToken("GET", "/", "123")

  if response.Code != 200 {
    t.Fatal("Response code should be 200 was", response.Code)
  }

  responseSession := UserSession{}
  err = json.Unmarshal([]byte(response.Body.String()), &responseSession)
  if err != nil {
    t.Fatal(err.Error())
  }

  if (responseSession.User != expectSession.User) ||
   (responseSession.Token != expectSession.Token) {
    t.Fatal("Expect", expectSession, "was", responseSession)
  }

}

func Test_VerifyAuth_Fail(t *testing.T) {
  globals := MakeTestGlobals(t)
  defer globals.MongoDB.Session.Close()
  db := globals.MongoDB.Session.DB(TestDBName)
  defer db.DropDatabase()

  sessionsColl := db.C(SessionsColl)
  expires := time.Now()
  sessionsColl.Insert(UserSession{Token: "12", User: "loveMaster_999", Expires: expires})

  h := gin.New()
  auth := MakeGlobalsHandler(Auth, globals)
  h.GET("/", auth, func(c *gin.Context){c.JSON(http.StatusOK, gin.H{"some":"thing"})})
  request := TestRequest{
    Body: "",
    Header: http.Header{},
    Handler: h,
  }
  response := request.SendWithToken("GET", "/", "123")

  if response.Code != 401 {
    t.Fatal("Response code should be 401 was", response.Code)
  }

}

func Test_VerifyAuth_ExpiresFail(t *testing.T) {
  globals := MakeTestGlobals(t)
  defer globals.MongoDB.Session.Close()
  db := globals.MongoDB.Session.DB(TestDBName)
  defer db.DropDatabase()

  sessionsColl := db.C(SessionsColl)
  expires := time.Now().AddDate(0,0,-1)
  sessionsColl.Insert(UserSession{Token: "123", User: "loveMaster_999", Expires: expires})

  h := gin.New()
  auth := MakeGlobalsHandler(Auth, globals)
  h.GET("/", auth, func(c *gin.Context){c.JSON(http.StatusOK, gin.H{"some":"thing"})})
  request := TestRequest{
              Body: "",
              Header: http.Header{},
              Handler: h,
            }
  response := request.SendWithToken("GET", "/", "123")

  if response.Code != 401 {
    t.Fatal("Response code should be 401 was", response.Code)
  }

}

