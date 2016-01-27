package dbVars

var (
	DBVarsTable = "db_vars"
)

type (
	DBVar struct {
		Name  string `db:"name" json:"name" valid:"required"`
		Value string `db:"value" json:"value" valid:"required"`
	}
)
