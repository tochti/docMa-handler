package bebber

import (
	"crypto/sha1"
	"errors"
	"fmt"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func (user *User) Read(username string, db *mgo.Database) error {
	u := *user
	usersColl := db.C(UsersColl)
	query := usersColl.Find(bson.M{"username": username})

	n, err := query.Count()
	if err != nil {
		return err
	}
	if n != 1 {
		errMsg := fmt.Sprintf("Cannot find user %v", username)
		return errors.New(errMsg)
	}

	err = query.One(&u)
	if err != nil {
		return err
	}
	*user = u
	return nil
}

func (user *User) Save(db *mgo.Database) error {
	u := *user
	u.Password = fmt.Sprintf("%x", sha1.Sum([]byte("tt")))
	usersColl := db.C(UsersColl)
	err := usersColl.Insert(u)
	if err != nil {
		return err
	}

	return nil
}
