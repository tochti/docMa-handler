package bebber

import (
  "gopkg.in/mgo.v2"
  "github.com/gin-gonic/gin"
)

// Struct to store obj which are interesting for many Handler Functions
type Globals struct {
  Config Config
  MongoDB MongoDBConn
}

type Config map[string]string

type MongoDBConn struct {
  Session *mgo.Session
  DialInfo *mgo.DialInfo
  DBName string
}

func MakeGlobalsHandler(fn func(*gin.Context, Globals), globals Globals) func(*gin.Context) {
  return func (c *gin.Context) {
           fn(c, globals)
         }
}
