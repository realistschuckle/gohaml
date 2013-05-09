package gohaml

import (
	"encoding/json"
	"github.com/realistschuckle/gohaml"
	"os"
	"testing"
)

const spec_fn = "test/tests.json"

var scope map[string]interface{}

type haml_test struct {
	Haml   string
	Html   string
  Optional bool
	Config map[string]string
  Locals map[string]interface{}
}

type test_set map[string]haml_test
type test_all map[string]test_set

func TestSpecs(t *testing.T) {
	num_tests := 0
	num_passed := 0
	num_failed := 0

	var file *os.File
	var err error

	if file, err = os.Open(spec_fn); err != nil {
		t.Fatalf(err.Error())
	}
	var tests test_all
	decoder := json.NewDecoder(file)
	if err = decoder.Decode(&tests); err != nil {
		t.Fatalf(err.Error())
	}
	for category_name := range tests {
		t.Log("Category: " + category_name)
		category := tests[category_name]
		for test_name := range category {
			t.Log("Test :" + test_name)
			num_tests += 1
			test := category[test_name]
			engine, herr := gohaml.NewEngine(test.Haml)
			if herr != nil {
				t.Error(herr.Error())
				num_failed += 1
			} else {
        var html string
        if test.Locals != nil {
				  html = engine.Render(test.Locals)
        } else {
          html = engine.Render(scope)
        }
				if html != test.Html {
					t.Errorf("expected: ->%s\n", test.Html)
					t.Errorf("got     : ->%s\n", html)
					num_failed += 1
				} else {
					num_passed += 1
				}
			}
		} // test in category
	} // category in sll
  t.Logf("passed: %d failed %d of %d\n", num_passed, num_failed, num_tests);
}
