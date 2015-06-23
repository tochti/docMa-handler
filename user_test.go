package bebber

import (
  "fmt"
  "testing"
  "crypto/sha1"

  "gopkg.in/mgo.v2"
  "gopkg.in/mgo.v2/bson"
)

func Test_ReadUser_OK(t *testing.T) {
  session, err := mgo.Dial("127.0.0.1")
  if err != nil {
    t.Fatal(err.Error())
  }
  defer session.Close()

  userTmp := User{Username: "XXX", Password: ""}
  userExpect := User{Username: "Haschel", Password: ""}
  db := session.DB(TestDBName)
  col := db.C(UsersColl)
  defer db.DropDatabase()

  err = col.Insert(userExpect, userTmp)
  if err != nil {
    t.Fatal(err.Error())
  }

  user := User{}
  err = user.Read("Haschel", db)

  if (userExpect.Username != user.Username) && (err == nil) {
    t.Fatal("Expect,", userExpect, "was,", user)
  }

}

func Test_ReadUser_Fail(t *testing.T) {
  session, err := mgo.Dial("127.0.0.1")
  if err != nil {
    t.Fatal(err.Error())
  }
  defer session.Close()

  db := session.DB(TestDBName)
  defer db.DropDatabase()

  user := User{}
  err = user.Read("Haschel", db)

  if err.Error() != "Cannot find user Haschel" {
    t.Fatal("Expect 'Cannot found user Haschel' error was", err.Error())
  }

}

func Test_SaveUser_OK(t *testing.T) {
  session, err := mgo.Dial("127.0.0.1")
  if err != nil {
    t.Fatal(err.Error())
  }
  defer session.Close()

  db := session.DB(TestDBName)
  defer db.DropDatabase()

  sha1Pass := fmt.Sprintf("%x", sha1.Sum([]byte("tt")))
  userExpect := User{Username: "test", Password: "tt"}
  err = userExpect.Save(db)

  user := User{}
  usersColl := db.C(UsersColl)
  err = usersColl.Find(bson.M{"username": "test"}).One(&user)
  if err != nil {
    t.Fatal(err.Error())
  }

  if (userExpect.Username != user.Username) && (err == nil) {
    t.Fatal("Expect", userExpect, "was", user)
  }

  if (sha1Pass != user.Password) {
    t.Fatal("Expect", sha1Pass, "was", user.Password)
  }

}
