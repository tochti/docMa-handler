package bebber

import (
  "net/http"
  "github.com/gin-gonic/gin"
  "gopkg.in/mgo.v2"
  "gopkg.in/mgo.v2/bson"
)

const (
  sessionsCollection = "sessions"
  tokenHeaderField = "X-XSRF-TOKEN"
)

func Auth() gin.HandlerFunc {
  return func(c *gin.Context) {

    sessionKey := c.Request.Header.Get(tokenHeaderField)
    if sessionKey == "" {
      c.JSON(http.StatusUnauthorized, ErrorResponse{"fail", "Header not found"})
      c.Abort()
      return
    }

    session, err := mgo.Dial(GetSettings("BEBBER_DB_SERVER"))
    if err != nil {
      c.JSON(http.StatusUnauthorized, ErrorResponse{"fail", err.Error()})
      c.Abort()
      return
    }
    defer session.Close()

    sessionsC := session.DB(GetSettings("BEBBER_DB_NAME")).C(sessionsCollection)
    n, err := sessionsC.Find(bson.M{"key": sessionKey}).Count()
    if (err != nil) {
      c.JSON(http.StatusUnauthorized, ErrorResponse{"fail", err.Error()})
      c.Abort()
      return
    }

    if n != 1 {
      c.JSON(http.StatusUnauthorized, ErrorResponse{"fail", "Session not found"})
      c.Abort()
      return
    } else {
      c.Next()
    }
  }
}
