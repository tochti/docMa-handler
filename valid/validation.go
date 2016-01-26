package valid

import "gopkg.in/go-playground/validator.v8"

var (
	validate *validator.Validate
	TagName  = "valid"
	config   = &validator.Config{TagName: TagName}
	Validate = validator.New(config)
	Struct   = Validate.Struct
)
