package bebber

import (
	"github.com/gin-gonic/gin"
	"testing"
)

func Test_MakeGlobalsHandler(t *testing.T) {

	globals := Globals{Config: map[string]string{"Burn": "Motherfucker"}}

	var globalsReturn Globals
	handler := func(c *gin.Context, globals Globals) {
		globalsReturn = globals
	}

	fn := MakeGlobalsHandler(handler, globals)
	tmpContext := gin.Context{}
	fn(&tmpContext)

	v, ok := globalsReturn.Config["Burn"]
	if (ok != true) || (v != "Motherfucker") {
		t.Fatal("Globals missing")
	}

}
