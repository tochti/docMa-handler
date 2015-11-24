package bebber

import (
	"strings"
	"testing"
)

func Test_ReadCurrDocNumberProposal_OK(t *testing.T) {
	g := MakeTestGlobals(t)
	session := g.MongoDB.Session.Copy()
	defer session.Close()
	db := session.DB(TestDBName)
	defer db.DropDatabase()

	varsColl := db.C(VarsColl)

	expectProposal := DocNumberProposal(1234)
	err := varsColl.Insert(
		VarDocNumberProposal{
			Name:     DocNumberProposalName,
			Proposal: expectProposal,
		},
	)
	if err != nil {
		t.Fatal(err.Error())
	}

	var docNumberProposal DocNumberProposal
	no, err := docNumberProposal.Curr(db)
	if err != nil {
		t.Fatal(err.Error())
	}

	if expectProposal != no {
		t.Fatal("Execpt", expectProposal, "was", no)
	}
}

func Test_ReadCurrDocNumberProposal_NoneFail(t *testing.T) {
	g := MakeTestGlobals(t)
	session := g.MongoDB.Session.Copy()
	defer session.Close()
	db := session.DB(TestDBName)
	defer db.DropDatabase()

	var docNumberProposal DocNumberProposal
	_, err := docNumberProposal.Curr(db)
	errMsg := "not found"
	if strings.Contains(errMsg, err.Error()) == false {
		t.Fatal("Exepct", errMsg, "was", err.Error())
	}
}

func Test_MakeDocNumberProposal_OK(t *testing.T) {
	g := MakeTestGlobals(t)
	session := g.MongoDB.Session.Copy()
	defer session.Close()
	db := session.DB(TestDBName)
	defer db.DropDatabase()

	expectDocNumberProposal := DocNumberProposal(1234)
	err := expectDocNumberProposal.Save(db)
	if err != nil {
		t.Fatal(err.Error())
	}

	no, err := expectDocNumberProposal.Curr(db)
	if err != nil {
		t.Fatal(err.Error())
	}

	if expectDocNumberProposal != no {
		t.Fatal("Expect", expectDocNumberProposal, "was", no)
	}
}

func Test_ChangeDocNumberProposal_OK(t *testing.T) {
	g := MakeTestGlobals(t)
	session := g.MongoDB.Session.Copy()
	defer session.Close()
	db := session.DB(TestDBName)
	defer db.DropDatabase()

	docNumberProposal := DocNumberProposal(12)
	err := docNumberProposal.Save(db)
	if err != nil {
		t.Fatal(err.Error())
	}

	expectDocNumberProposal := DocNumberProposal(1234)
	err = expectDocNumberProposal.Save(db)
	if err != nil {
		t.Fatal(err.Error())
	}

	no, err := expectDocNumberProposal.Curr(db)
	if err != nil {
		t.Fatal(err.Error())
	}

	if expectDocNumberProposal != no {
		t.Fatal("Expect", expectDocNumberProposal, "was", no)
	}
}

func Test_NextDocNumberProposal_OK(t *testing.T) {
	g := MakeTestGlobals(t)
	session := g.MongoDB.Session.Copy()
	defer session.Close()
	db := session.DB(TestDBName)
	defer db.DropDatabase()

	docNumberProposal := DocNumberProposal(12)
	err := docNumberProposal.Save(db)
	if err != nil {
		t.Fatal(err.Error())
	}

	no, err := docNumberProposal.Next(db)
	if err != nil {
		t.Fatal(err.Error())
	}

	if 13 != no {
		t.Fatal("Expect 13 was", no)
	}
}
