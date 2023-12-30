package acid_test

import (
	"reflect"
	"testing"

	"github.com/jeffrydegrande/acid"
)

func TestPinning(t *testing.T) {
	t.Parallel()
	acid.UseCDN(acid.JsDelivr)
	acid.Pin("package1", "1.0.0")
	acid.Pin("package2", "2.0.0")

	expected := []acid.Package{
		{"package1", "https://cdn.jsdelivr.net/npm/package1@1.0.0/+esm"},
		{"package2", "https://cdn.jsdelivr.net/npm/package2@2.0.0/+esm"},
	}

	if p := acid.Packages(); !reflect.DeepEqual(p, expected) {
		t.Errorf("Expected %v, got %v", expected, p)
	}
}
