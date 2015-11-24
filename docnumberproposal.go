package bebber

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func (_ DocNumberProposal) Curr(db *mgo.Database) (DocNumberProposal, error) {
	coll := db.C(VarsColl)
	query := coll.Find(
		bson.M{"name": DocNumberProposalName},
	)
	result := VarDocNumberProposal{}
	err := query.One(&result)
	if err != nil {
		return -1, err
	} else {
		return result.Proposal, nil
	}
}

func (prop DocNumberProposal) Save(db *mgo.Database) error {
	coll := db.C(VarsColl)
	_, err := coll.Upsert(
		bson.M{"name": DocNumberProposalName},
		VarDocNumberProposal{
			Name:     DocNumberProposalName,
			Proposal: prop,
		},
	)
	if err != nil {
		return err
	}

	return nil
}

func (prop DocNumberProposal) Next(db *mgo.Database) (DocNumberProposal, error) {
	no, err := prop.Curr(db)
	if err != nil {
		return -1, err
	} else {
		return no + 1, nil
	}
}
