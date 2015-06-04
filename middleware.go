package bebber

import (
  _"fmt"
  "time"
  "net/http"
  "github.com/gin-gonic/gin"
  "gopkg.in/mgo.v2"
  "gopkg.in/mgo.v2/bson"
)

const (
  SessionsCollection = "sessions"
  tokenHeaderField = "X-XSRF-TOKEN"
)

func Auth() gin.HandlerFunc {
  return func(c *gin.Context) {

    token := c.Request.Header.Get(tokenHeaderField)
    if token == "" {
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

    sessionsC := session.DB(GetSettings("BEBBER_DB_NAME")).C(SessionsCollection)
    query := sessionsC.Find(bson.M{"token": token})
    n, err := query.Count()
    if err != nil {
      c.JSON(http.StatusUnauthorized, ErrorResponse{"fail", err.Error()})
      c.Abort()
      return
    }
    if n != 1 {
      c.JSON(http.StatusUnauthorized, ErrorResponse{"fail", "Session not found"})
      c.Abort()
      return

    }

    userSession := UserSession{}
    err = query.One(&userSession)
    if err != nil {
      c.JSON(http.StatusUnauthorized, ErrorResponse{"fail", err.Error()})
      c.Abort()
      return
    }
    if userSession.Expires.Before(time.Now()) {
      c.JSON(http.StatusUnauthorized, ErrorResponse{"fail", "Session expired"})
      c.Abort()
      return
    } else {
      c.Next()
    }

  }

}
