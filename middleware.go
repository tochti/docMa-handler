package bebber

import (
	_ "fmt"
	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"time"
)

func Auth(c *gin.Context, g Globals) {

	token := c.Request.Header.Get(TokenHeaderField)
	if token == "" {
		cookie, err := c.Request.Cookie(XSRFCookieName)
		if err != nil {
			c.JSON(http.StatusUnauthorized, FailResponse{"fail", "Cookie not found"})
			c.Abort()
			return
		}
		token = cookie.Value
		if token == "" {
			c.JSON(http.StatusUnauthorized, FailResponse{"fail", "Header not found"})
			c.Abort()
			return
		}
	}

	session := g.MongoDB.Session.Copy()
	defer session.Close()

	sessionsColl := session.DB(g.MongoDB.DBName).C(SessionsColl)
	query := sessionsColl.Find(bson.M{"token": token})
	n, err := query.Count()
	if err != nil {
		c.JSON(http.StatusUnauthorized, FailResponse{"fail", err.Error()})
		c.Abort()
		return
	}
	if n != 1 {
		c.JSON(http.StatusUnauthorized, FailResponse{"fail", "Session not found"})
		c.Abort()
		return

	}

	userSession := UserSession{}
	err = query.One(&userSession)
	if err != nil {
		c.JSON(http.StatusUnauthorized, FailResponse{"fail", err.Error()})
		c.Abort()
		return
	}
	if userSession.Expires.Before(time.Now()) {
		c.JSON(http.StatusUnauthorized, FailResponse{"fail", "Session expired"})
		c.Abort()
		return
	} else {
		c.Set("session", userSession)
		c.Next()
	}

}
