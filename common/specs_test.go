package common

import (
	"os"
	"testing"
)

func Test_LoadSpecs(t *testing.T) {
	os.Setenv("DOCMA_PUBLIC", "1")
	os.Setenv("DOCMA_PDFVIEWER_PUBLIC", "2")
	os.Setenv("DOCMA_FILES", "3")

	specs := LoadSpecs()

	if specs.Public != "1" {
		t.Fatalf("Expect %v was %v", "1", specs.Public)
	}

	if specs.PDFViewerPublic != "2" {
		t.Fatalf("Expect %v was %v", "2", specs.PDFViewerPublic)
	}

	if specs.Files != "3" {
		t.Fatalf("Expect %v was %v", "3", specs.Files)
	}
}
