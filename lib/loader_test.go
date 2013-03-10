package gohaml

import (
	"testing"
	"fmt"
)

func TestLoadFile (t *testing.T) {
	var fsl Loader
	var err error

	if fsl, err = NewFileSystemLoader("."); err != nil {
		t.Errorf("couldn't create fileSystemLoader: %s", err)
	}

	if _, err = fsl.Load(1); err == nil {
		t.Errorf("rats! expected error")
	}

	if _, err = fsl.Load("test.haml"); err != nil {
		t.Errorf("couldn't load: test.haml: %s", err);
	} 


	if fsl, err = NewFileSystemLoader("blsadfasdf"); err == nil {
		t.Errorf("rats! expected error for non existing dir ...  ")
	}
}
