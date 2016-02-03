package common

import (
	"log"

	"github.com/kelseyhightower/envconfig"
)

var (
	AppName = "docma"
)

type (
	Specs struct {
		Public          string `envconfig:"PUBLIC"`
		PDFViewerPublic string `envconfig:"PDFVIEWER_PUBLIC"`
		Files           string `envconfig:"FILES"`
	}
)

func LoadSpecs() Specs {
	specs := Specs{}
	err := envconfig.Process(AppName, &specs)
	if err != nil {
		log.Fatal(err)
	}

	return specs
}
