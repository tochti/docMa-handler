package bebber

import (
  _"fmt"
  "os"
  "time"
  "bytes"
  "testing"
  "net/http"
  "crypto/sha1"
  "gopkg.in/mgo.v2"
  "gopkg.in/mgo.v2/bson"
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
  userId := bson.NewObjectId()
  sha1Pass := sha1.Sum([]byte("test"))
  user := bson.M{"_id": userId, "username": "loveMaster_999", "password": sha1Pass}
  usersC := session.DB("bebber_test").C("users")
  usersC.Insert(user)
  sessionsC := session.DB("bebber_test").C("sessions")
  createDate := time.Now()
  sessionsC.Insert(bson.M{"key": "123", "user": userId, "createDate": createDate})
  defer session.DB("bebber_test").DropDatabase()

  h := gin.New()
  h.GET("/", func(c *gin.Context){c.JSON(http.StatusOK, gin.H{"some":"thing"})}, Auth())
  header := http.Header{}
  header.Add("X-AUTH-TOKEN", "123")
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
  sha1Pass := sha1.Sum([]byte("test"))
  user := bson.M{"username": "loveMaster_999", "password": sha1Pass}
  usersC := session.DB("bebber_test").C("users")
  usersC.Insert(user)
  sessionsC := session.DB("bebber_test").C("sessions")
  createDate := time.Now()
  sessionsC.Insert(bson.M{"key": "12", "user": "loveMaster_999", "createDate": createDate})
  defer session.DB("bebber_test").DropDatabase()

  h := gin.New()
  h.GET("/", Auth(), func(c *gin.Context){c.JSON(http.StatusOK, gin.H{"some":"thing"})})
  header := http.Header{}
  header.Add("X-AUTH-TOKEN", "123")
  body := bytes.NewBufferString("")
  resp := PerformRequestHeader(h, "GET", "/", body, &header)

  if resp.Code != 401 {
    t.Fatal("Response code should be 401 was", resp.Code)
  }

}

