package bebber

import (
  _"fmt"
  "os"
  "time"
  "bytes"
  "testing"
  "net/http"
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
  sessionsC.Insert(UserSession{Token: "123", User: "loveMaster_999", Expires: expires})
  defer session.DB("bebber_test").DropDatabase()

  h := gin.New()
  h.GET("/", Auth(), func(c *gin.Context){c.JSON(http.StatusOK, gin.H{"some":"thing"})})
  header := http.Header{}
  header.Add("X-XSRF-TOKEN", "123")
  body := bytes.NewBufferString("")
  resp := PerformRequestHeader(h, "GET", "/", body, &header)

  if resp.Code != 200 {
    t.Fatal("Response code should be 200 was", resp.Code)
  }

  if resp.Body.String() != "{\"some\":\"thing\"}\n" {
    t.Fatal("Response Body should be {\"some\":\"thing\"} was", resp.Body.String())
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

