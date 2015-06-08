package bebber

import (
  _"fmt"
  "os"
  "time"
  "bytes"
  "testing"
  "net/http"
  "encoding/json"
  "gopkg.in/mgo.v2"
  "github.com/gin-gonic/gin"
)

func TestVerifyAuthOK(t *testing.T) {
  os.Setenv("BEBBER_DB_SERVER", "127.0.0.1")
  os.Setenv("BEBBER_DB_NAME", "bebber_test")
  session, err := mgo.Dial("127.0.0.1")
  if err != nil {
    t.Fatal(err)
  }
  defer session.Close()
  sessionsC := session.DB("bebber_test").C("sessions")
  expires := time.Now().AddDate(0,0,1)
  userSession := UserSession{Token: "123", User: "loveMaster_999", Expires: expires}
  sessionsC.Insert(userSession)
  defer session.DB("bebber_test").DropDatabase()

  h := gin.New()
  afterAuth := func (c *gin.Context) {
    se, _ := c.Get("session")
    c.JSON(http.StatusOK, se)
  }
  h.GET("/", Auth(), afterAuth)
  header := http.Header{}
  header.Add("X-XSRF-TOKEN", "123")
  body := bytes.NewBufferString("")
  resp := PerformRequestHeader(h, "GET", "/", body, &header)

  if resp.Code != 200 {
    t.Fatal("Response code should be 200 was", resp.Code)
  }

  respSession := UserSession{}
  err = json.Unmarshal([]byte(resp.Body.String()), &respSession)
  if err != nil {
    t.Fatal(err.Error())
  }

  if (respSession.User != userSession.User) ||
   (respSession.Token != userSession.Token) {
    t.Fatal("Expect", userSession, "was", respSession)
  }

}

func TestVerifyAuthFail(t *testing.T) {
  os.Setenv("BEBBER_DB_SERVER", "127.0.0.1")
  os.Setenv("BEBBER_DB_NAME", "bebber_test")
  session, err := mgo.Dial("127.0.0.1")
  if err != nil {
    t.Fatal(err)
  }
  defer session.Close()
  sessionsC := session.DB("bebber_test").C("sessions")
  expires := time.Now()
  sessionsC.Insert(UserSession{Token: "12", User: "loveMaster_999", Expires: expires})
  defer session.DB("bebber_test").DropDatabase()

  h := gin.New()
  h.GET("/", Auth(), func(c *gin.Context){c.JSON(http.StatusOK, gin.H{"some":"thing"})})
  header := http.Header{}
  header.Add("X-XSRF-TOKEN", "123")
  body := bytes.NewBufferString("")
  resp := PerformRequestHeader(h, "GET", "/", body, &header)

  if resp.Code != 401 {
    t.Fatal("Response code should be 401 was", resp.Code)
  }

}

func TestVerifyAuthExpiresFail(t *testing.T) {
  os.Setenv("BEBBER_DB_SERVER", "127.0.0.1")
  os.Setenv("BEBBER_DB_NAME", "bebber_test")
  session, err := mgo.Dial("127.0.0.1")
  if err != nil {
    t.Fatal(err)
  }
  defer session.Close()
  sessionsC := session.DB("bebber_test").C("sessions")
  expires := time.Now().AddDate(0,0,-1)
  sessionsC.Insert(UserSession{Token: "123", User: "loveMaster_999", Expires: expires})
  defer session.DB("bebber_test").DropDatabase()

  h := gin.New()
  h.GET("/", Auth(), func(c *gin.Context){c.JSON(http.StatusOK, gin.H{"some":"thing"})})
  header := http.Header{}
  header.Add("X-XSRF-TOKEN", "123")
  body := bytes.NewBufferString("")
  resp := PerformRequestHeader(h, "GET", "/", body, &header)

  if resp.Code != 401 {
    t.Fatal("Response code should be 401 was", resp.Code)
  }

}

